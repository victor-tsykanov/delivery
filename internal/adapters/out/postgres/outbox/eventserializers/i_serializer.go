package eventserializers

import "github.com/victor-tsykanov/delivery/internal/common/ddd"

type ISerializer interface {
	Serialize(event ddd.IDomainEvent) ([]byte, error)
	Deserialize(data []byte) (ddd.IDomainEvent, error)
}
