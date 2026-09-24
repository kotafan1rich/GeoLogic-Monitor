package user

import (
	"context"
	"errors"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/database"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/domain"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/errs"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/user/dto"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/user/model"
	"github.com/kotafan1rich/GeoLogic-Monitor/api/internal/repository/database/user/query"
)

type repository struct {
	db database.DBTX
}

func NewRepository(db database.DBTX) *repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, user *domain.User) (*domain.User, error) {
	userModel := dto.ToModel(*user)
	err := scanUser(
		r.db.QueryRow(ctx, query.Upsert, userModel.MaxUserID, userModel.MaxChatID),
		userModel,
	)
	if err != nil {
		return nil, err
	}

	return dto.ToDomain(*userModel), nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return r.get(ctx, query.GetByID, id)
}

func (r *repository) GetByMaxUserID(ctx context.Context, maxUserID int64) (*domain.User, error) {
	return r.get(ctx, query.GetByMaxUserID, maxUserID)
}

func (r *repository) get(ctx context.Context, queryString string, arg any) (*domain.User, error) {
	userModel := model.User{}
	err := scanUser(r.db.QueryRow(ctx, queryString, arg), &userModel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return dto.ToDomain(userModel), nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(row scanner, user *model.User) error {
	return row.Scan(
		&user.ID,
		&user.MaxUserID,
		&user.MaxChatID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
