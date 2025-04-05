package ddd

type IAggregateRoot interface {
	DomainEvents() []IDomainEvent
	ClearDomainEvents()
}
