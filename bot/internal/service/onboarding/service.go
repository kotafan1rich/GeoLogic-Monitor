package onboarding

import (
	"context"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type Messager interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type ApiRepository interface {
	PutUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type service struct {
	apiRepo  ApiRepository
	messager Messager
}

func NewService(apiRepo ApiRepository, messager Messager) *service {
	return &service{
		apiRepo:  apiRepo,
		messager: messager,
	}
}

func (s *service) RegisterUser(ctx context.Context, maxUserID, maxChatID int64) error {
	user := domain.NewUser(maxUserID, maxChatID)
	user, err := s.apiRepo.PutUser(ctx, user)
	if err != nil {
		return err
	}

	err = s.messager.SendMessage(ctx, maxChatID, OnboardingMessage)
	if err != nil {
		return err
	}
	return nil
}
