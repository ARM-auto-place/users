package postgres

import (
	"context"

	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/google/uuid"
)

type UsersStorage interface {
	GetProfile(context.Context, uuid.UUID) (models.Profile, error)
	GetMyCars(context.Context, uuid.UUID) ([]models.MyCar, error)
	GetFavoriteCars(context.Context, uuid.UUID) ([]models.FavoriteCar, error)
	UpdateProfile(context.Context, uuid.UUID, models.ProfilePatch) (models.Profile, error)
	UpdateSettings(context.Context, uuid.UUID, models.SettingsPatch) (models.Settings, error)
}
