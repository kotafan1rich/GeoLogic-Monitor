package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs/app"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/logger"
)

type UserRepository interface {
	Upsert(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByMaxUserID(ctx context.Context, maxUserID int64) (*domain.User, error)
}
type service struct {
	repo UserRepository
	log  *logger.Logger
}

func NewUserService(log *logger.Logger, repo UserRepository) *service {
	return &service{log: log, repo: repo}
}

func (s *service) Upsert(
	ctx context.Context,
	maxUserID, maxChatID int64,
) (*domain.User, error) {
	user, err := domain.NewUser(maxUserID, maxChatID)
	if err != nil {
		return nil, app.ValidationError(err)
	}

	user, err = s.repo.Upsert(ctx, user)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to upsert user",
			slog.Int64("max_user_id", maxUserID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return user, nil
}

func (s *service) GetByMaxUserID(ctx context.Context, maxUserID int64) (*domain.User, error) {
	user, err := s.repo.GetByMaxUserID(ctx, maxUserID)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, app.Wrap(err, app.ErrNotFound)
		}

		s.log.ErrorContext(ctx,
			"failed to get user by MAX user ID",
			slog.Int64("max_user_id", maxUserID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return user, nil
}
