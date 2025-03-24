package app

import (
	"context"
	"log"

	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/courier"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/order"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/commands"
	"github.com/victor-tsykanov/delivery/internal/core/application/usecases/queries"
	"github.com/victor-tsykanov/delivery/internal/core/domain/services"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	outPorts "github.com/victor-tsykanov/delivery/internal/core/ports/out"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	DomainServices  DomainServices
	Repositories    Repositories
	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers
}

type DomainServices struct {
	DispatchService services.IDispatchService
}

type Repositories struct {
	TransactionManager outPorts.ITransactionManager
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

func NewCompositionRoot(_ context.Context, gormDb *gorm.DB) *CompositionRoot {
	dispatchService, err := services.NewDispatchService()
	if err != nil {
		log.Fatalf("faied to create DispatchService: %v", err)
	}

	transactionManager, err := postgres.NewGormTransactionManager(gormDb)
	if err != nil {
		log.Fatalf("faied to create GormTransactionManager: %v", err)
	}

	orderRepository, err := order.NewRepository(gormDb)
	if err != nil {
		log.Fatalf("faied to create order repository: %v", err)
	}

	courierRepository, err := courier.NewRepository(gormDb)
	if err != nil {
		log.Fatalf("faied to create courier repository: %v", err)
	}

	createOrderCommandHandler, err := commands.NewCreateOrderCommandHandler(transactionManager, orderRepository)
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
	}

	return &compositionRoot
}
