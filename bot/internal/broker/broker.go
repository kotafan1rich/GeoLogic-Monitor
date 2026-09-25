package broker

import "context"

type Message struct {
	Key   []byte
	Value []byte
}

type Handler func(ctx context.Context, message Message) error

type Consumer interface {
	Consume(ctx context.Context, handler Handler) error
	Close() error
}
