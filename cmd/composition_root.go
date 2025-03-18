package main

import (
	"context"
	"log"

	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/order"
	"github.com/victor-tsykanov/delivery/internal/core/domain/services"
	ports "github.com/victor-tsykanov/delivery/internal/core/ports/out"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	DomainServices DomainServices
	Repositories   Repositories
}

type DomainServices struct {
	DispatchService services.IDispatchService
}

type Repositories struct {
	TransactionManager ports.ITransactionManager
	OrderRepository    ports.IOrderRepository
}

func NewCompositionRoot(_ context.Context, gormDb *gorm.DB) CompositionRoot {
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

	compositionRoot := CompositionRoot{
		DomainServices: DomainServices{
			DispatchService: dispatchService,
		},
		Repositories: Repositories{
			TransactionManager: transactionManager,
			OrderRepository:    orderRepository,
		},
	}

	return compositionRoot
}
