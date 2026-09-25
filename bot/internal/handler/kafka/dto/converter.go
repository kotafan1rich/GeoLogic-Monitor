package dto

import (
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain/notification"
)

func ToDomain(notice *Notification) *notification.Notification {
	return &notification.Notification{
		Type:      notice.Type,
		MaxChatID: notice.MaxChatID,
		TrackedLocation: domain.TrackedLocation{
			ID:            notice.TrackedLocation.ID,
			Name:          notice.TrackedLocation.Name,
			Address:       notice.TrackedLocation.Address,
			BusnessTypeID: notice.TrackedLocation.BusinessTypeID,
		},
		Subject:    notification.Subject(notice.Subject),
		Route:      notification.Route(notice.Route),
		Assessment: notification.Assessment(notice.Assessment),
	}
}
