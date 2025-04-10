package eventserializers

import "github.com/victor-tsykanov/delivery/internal/common/ddd"

type IRegistry interface {
	Serialize(event ddd.IDomainEvent) ([]byte, error)
	Deserialize(eventType string, data []byte) (ddd.IDomainEvent, error)
}
