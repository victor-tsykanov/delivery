package main

import (
	"context"
	"log"

	"github.com/victor-tsykanov/delivery/internal/core/domain/services"
)

type CompositionRoot struct {
	DomainServices DomainServices
}

type DomainServices struct {
	DispatchService services.IDispatchService
}

func NewCompositionRoot(_ context.Context) CompositionRoot {
	dispatchService, err := services.NewDispatchService()
	if err != nil {
		log.Fatalf("faied to create DispatchService: %v", err)
	}

	compositionRoot := CompositionRoot{
		DomainServices: DomainServices{
			DispatchService: dispatchService,
		},
	}

	return compositionRoot
}
