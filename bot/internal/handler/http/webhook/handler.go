package webhook

import (
	"context"
	"net/http"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, update domain.Update) error
}

type handler struct {
	dispatcher Dispatcher
}

func NewHandler(dispatcher Dispatcher) *handler {
	return &handler{
		dispatcher: dispatcher,
	}
}

func (h *handler) Handle(w http.Response, r *http.Request) {

}
