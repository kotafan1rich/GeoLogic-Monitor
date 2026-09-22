package start

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
)

type Service interface {
	RegisterUser(ctx context.Context, maxUserID, maxChatID int64) error
}

type handler struct {
	log         *logger.Logger
	userService Service
}

func NewStartHandler(log *logger.Logger, userService Service) *handler {
	return &handler{
		log:         log,
		userService: userService,
	}
}

func (h *handler) Start(ctx context.Context, chatID, userID int64) {
	err := h.userService.RegisterUser(ctx, userID, chatID)
	if err != nil {
		h.log.Error("failed to process command", "err", err)
	}
}
