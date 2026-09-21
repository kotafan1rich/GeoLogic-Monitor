package bot

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type dispatcher struct {
}

func NewDispatcher() *dispatcher {
	return &dispatcher{}
}

func (d *dispatcher) Dispatch(ctx context.Context, update domain.Update) error
