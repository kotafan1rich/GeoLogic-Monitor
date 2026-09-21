package bot

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type StartHandler interface {
	Start()
}

type dispatcher struct {
	startHandler StartHandler
}

func NewDispatcher(startHandler StartHandler) *dispatcher {
	return &dispatcher{
		startHandler: startHandler,
	}
}

func (d *dispatcher) Dispatch(ctx context.Context, update domain.Update) error
