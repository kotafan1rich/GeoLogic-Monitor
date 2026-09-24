package service

import "context"

type OnboardingService interface {
	RegisterUser(ctx context.Context, maxUserID, maxChatID int64) error
}
