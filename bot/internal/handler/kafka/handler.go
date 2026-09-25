package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/broker"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain/notification"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/handler/kafka/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/logger"
)

type Service interface {
	Notify(ctx context.Context, notice *notification.Notification) error
}

type handler struct {
	log     *logger.Logger
	service Service
}

func NewHandler(log *logger.Logger, service Service) *handler {
	return &handler{
		log:     log,
		service: service,
	}
}

func (h *handler) Handle(ctx context.Context, message broker.Message) error {
	var notice dto.Notification

	if err := json.Unmarshal(message.Value, &notice); err != nil {
		return fmt.Errorf("decode notification: %w", err)
	}

	h.log.InfoContext(
		ctx,
		"processing notification",
		"type", notice.Type,
		"max_chat_id", notice.MaxChatID,
		"tracked_location_id", notice.TrackedLocation.ID,
		"subject_external_id", notice.Subject.ExternalID,
	)

	if err := h.service.Notify(ctx, dto.ToDomain(&notice)); err != nil {
		if errors.Is(err, errs.ErrChatNotFound) {
			h.log.WarnContext(
				ctx,
				"notification skipped because MAX chat was not found",
				"max_chat_id", notice.MaxChatID,
			)
			return nil
		}
		return fmt.Errorf("process notification: %w", err)
	}

	return nil
}
