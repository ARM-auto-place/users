package middlewares

import (
	"context"

	externalapi "github.com/ARM-auto-place/users/pkg/interfaces/externalapi"
	customErrors "github.com/ARM-auto-place/users/pkg/models/errors"
	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/google/uuid"
)

type accessRoleClient interface {
	GetRole(context.Context, uuid.UUID) (string, error)
}

type accessMiddleware struct {
	next   externalapi.ExternalAPI
	access accessRoleClient
}

func NewAccessMiddleware(next externalapi.ExternalAPI, access accessRoleClient) externalapi.ExternalAPI {
	return &accessMiddleware{next: next, access: access}
}

func (m *accessMiddleware) GetProfile(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	if err := m.checkRole(ctx, id); err != nil {
		return models.Profile{}, err
	}
	return m.next.GetProfile(ctx, id)
}

func (m *accessMiddleware) GetMyCars(ctx context.Context, id uuid.UUID) ([]models.MyCar, error) {
	if err := m.checkRole(ctx, id); err != nil {
		return nil, err
	}
	return m.next.GetMyCars(ctx, id)
}

func (m *accessMiddleware) GetFavoriteCars(ctx context.Context, id uuid.UUID) ([]models.FavoriteCar, error) {
	if err := m.checkRole(ctx, id); err != nil {
		return nil, err
	}
	return m.next.GetFavoriteCars(ctx, id)
}

func (m *accessMiddleware) UpdateProfile(ctx context.Context, id uuid.UUID, patch models.ProfilePatch) (models.Profile, error) {
	if err := m.checkRole(ctx, id); err != nil {
		return models.Profile{}, err
	}
	return m.next.UpdateProfile(ctx, id, patch)
}

func (m *accessMiddleware) UpdateSettings(ctx context.Context, id uuid.UUID, patch models.SettingsPatch) (models.Settings, error) {
	if err := m.checkRole(ctx, id); err != nil {
		return models.Settings{}, err
	}
	return m.next.UpdateSettings(ctx, id, patch)
}

func (m *accessMiddleware) checkRole(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return customErrors.Unauthorized()
	}
	role, err := m.access.GetRole(ctx, id)
	if err != nil || (role != "admin" && role != "owner") {
		return customErrors.Forbidden()
	}
	return nil
}
