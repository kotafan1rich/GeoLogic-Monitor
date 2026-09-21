package api

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type ApiClient interface {
	PutUser(ctx context.Context, maxUserID, maxChatID int64) (*domain.User, error)
}

type respository struct {
	apiClient ApiClient
}

func NewRepository(apiClient ApiClient) *respository {
	return &respository{
		apiClient: apiClient,
	}
}

func (r *respository) PutUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	user, err := r.apiClient.PutUser(ctx, user.MaxUserId, user.MaxChatId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
