package webhook

import (
	"context"

	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, update model.Update)
}

type handler struct {
	dispatcher Dispatcher
}

func NewHandler(dispatcher Dispatcher) *handler {
	return &handler{
		dispatcher: dispatcher,
	}
}

func (h *handler) Handle(ctx context.Context, update model.Update) {
	h.dispatcher.Dispatch(ctx, update)
}
