package producer

import (
	"errors"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	ErrInvalidLogger    = errors.New("logger is nil")
	ErrProduceMsg       = errors.New("failed to produce message")
	ErrDeliveryMsg      = errors.New("failed to delivery message")
	ErrUnknownEventType = errors.New("invalid event type")
)

func isRetriable(err error) bool {
	if errors.Is(err, kgo.ErrRecordTimeout) || errors.Is(err, kgo.ErrRecordRetries) {
		return true
	}
	return kerr.IsRetriable(err)
}
