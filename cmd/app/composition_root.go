package app

import (
	"context"
	"fmt"
	"log"

	kafkaInAdapters "github.com/victor-tsykanov/delivery/internal/adapters/in/kafka"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/grpc"
	kafkaOutAdapters "github.com/victor-tsykanov/delivery/internal/adapters/out/kafka"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/courier"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/order"
	"github.com/victor-tsykanov/delivery/internal/common/config"
	"github.com/victor-tsykanov/delivery/internal/common/eventdispatcher"
	"github.com/victor-tsykanov/delivery/internal/common/persistence"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/commands"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/events"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/queries"
	"github.com/victor-tsykanov/delivery/internal/core/domain/services"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	outPorts "github.com/victor-tsykanov/delivery/internal/core/ports/out"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	DomainServices    DomainServices
	Repositories      Repositories
	CommandHandlers   CommandHandlers
	QueryHandlers     QueryHandlers
	QueueConsumers    QueueConsumers
	QueueProducers    QueueProducers
	shutdownCallbacks []shutdownCallback
}

type shutdownCallback func(ctx context.Context) error

func (r *CompositionRoot) Shutdown(ctx context.Context) {
	for _, callback := range r.shutdownCallbacks {
		err := callback(ctx)
		if err != nil {
			fmt.Println("shutdown callback failed with error:", err)
		}
	}
}

type DomainServices struct {
	DispatchService services.IDispatchService
}

type Repositories struct {
	TransactionManager persistence.ITransactionManager
	CourierRepository  outPorts.ICourierRepository
	OrderRepository    outPorts.IOrderRepository
}

type CommandHandlers struct {
	CreateOrderCommandHandler  inPorts.ICreateOrderCommandHandler
	MoveCouriersCommandHandler inPorts.IMoveCouriersCommandHandler
	AssignOrdersCommandHandler inPorts.IAssignOrdersCommandHandler
}

type QueryHandlers struct {
	GetAllCouriersQueryHandler   inPorts.IGetAllCouriersQueryHandler
	GetPendingOrdersQueryHandler inPorts.IGetPendingOrdersQueryHandler
}

type QueueConsumers struct {
	BasketConfirmedConsumer *kafkaInAdapters.BasketConfirmedConsumer
}

type QueueProducers struct {
	OrderCompletedEventProducer *kafkaOutAdapters.OrderCompletedEventProducer
}

func NewCompositionRoot(
	_ context.Context,
	gormDb *gorm.DB,
	geoServiceAddress string,
	kafkaConfig *config.KafkaConfig,
) *CompositionRoot {
	dispatchService, err := services.NewDispatchService()
	if err != nil {
		log.Fatalf("faied to create DispatchService: %v", err)
	}

	transactionManager, err := postgres.NewGormTransactionManager(gormDb)
	if err != nil {
		log.Fatalf("faied to create GormTransactionManager: %v", err)
	}

	eventDispatcher := eventdispatcher.New()

	orderRepository, err := order.NewRepository(gormDb, eventDispatcher)
	if err != nil {
		log.Fatalf("faied to create order repository: %v", err)
	}

	courierRepository, err := courier.NewRepository(gormDb)
	if err != nil {
		log.Fatalf("faied to create courier repository: %v", err)
	}

	geoClient, err := grpc.NewGeoClient(geoServiceAddress)
	if err != nil {
		log.Fatalf("faied to create Geo client: %v", err)
	}

	createOrderCommandHandler, err := commands.NewCreateOrderCommandHandler(
		transactionManager,
		orderRepository,
		geoClient,
	)
	if err != nil {
		log.Fatalf("faied to create CreateOrderCommandHandler: %v", err)
	}

	moveCouriersCommandHandler, err := commands.NewMoveCouriersCommandHandler(
		transactionManager,
		courierRepository,
		orderRepository,
	)
	if err != nil {
		log.Fatalf("faied to create MoveCouriersCommandHandler: %v", err)
	}

	assignOrdersCommandHandler, err := commands.NewAssignOrdersCommandHandler(
		transactionManager,
		dispatchService,
		courierRepository,
		orderRepository,
	)
	if err != nil {
		log.Fatalf("faied to create AssignOrdersCommandHandler: %v", err)
	}

	getAllCouriersQueryHandler, err := queries.NewGetAllCouriersQueryHandler(gormDb)
	if err != nil {
		log.Fatalf("faied to create GetAllCouriersQueryHandler: %v", err)
	}

	getPendingOrdersQueryHandler, err := queries.NewGetPendingOrdersQueryHandler(gormDb)
	if err != nil {
		log.Fatalf("faied to create GetPendingOrdersQueryHandler: %v", err)
	}

	basketConfirmedConsumer, err := kafkaInAdapters.NewBasketConfirmedConsumer(createOrderCommandHandler, kafkaConfig)
	if err != nil {
		log.Fatalf("faied to create BasketConfirmedConsumer: %v", err)
	}

	orderCompletedProducer, err := kafkaOutAdapters.NewOrderCompletedEventProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("faied to create OrderCompletedEventProducer: %v", err)
	}

	orderCompletedEventHandler, err := events.NewOrderCompletedEventHandler(orderCompletedProducer)
	if err != nil {
		log.Fatalf("faied to create OrderCompletedEventHandler: %v", err)
	}

	eventdispatcher.Register(eventDispatcher, orderCompletedEventHandler)

	compositionRoot := CompositionRoot{
		DomainServices: DomainServices{
			DispatchService: dispatchService,
		},
		Repositories: Repositories{
			TransactionManager: transactionManager,
			CourierRepository:  courierRepository,
			OrderRepository:    orderRepository,
		},
		CommandHandlers: CommandHandlers{
			CreateOrderCommandHandler:  createOrderCommandHandler,
			MoveCouriersCommandHandler: moveCouriersCommandHandler,
			AssignOrdersCommandHandler: assignOrdersCommandHandler,
		},
		QueryHandlers: QueryHandlers{
			GetAllCouriersQueryHandler:   getAllCouriersQueryHandler,
			GetPendingOrdersQueryHandler: getPendingOrdersQueryHandler,
		},
		QueueConsumers: QueueConsumers{
			BasketConfirmedConsumer: basketConfirmedConsumer,
		},
		QueueProducers: QueueProducers{
			OrderCompletedEventProducer: orderCompletedProducer,
		},
		shutdownCallbacks: []shutdownCallback{
			func(_ context.Context) error {
				return geoClient.Close()
			},
			func(_ context.Context) error {
				return basketConfirmedConsumer.Close()
			},
			func(_ context.Context) error {
				return orderCompletedProducer.Close()
			},
		},
	}

	return &compositionRoot
}
