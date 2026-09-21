package user

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type ApiRepository interface {
	PutUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type service struct {
	apiRepo ApiRepository
}

func NewService(apiRepo ApiRepository) *service {
	return &service{
		apiRepo: apiRepo,
	}
}

func (s *service) RegisterUser(ctx context.Context, maxUserID, maxChatID int64) error {
	user := domain.NewUser(maxUserID, maxChatID)
	user, err := s.apiRepo.PutUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
