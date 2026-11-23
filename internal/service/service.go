package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	commonErrors "github.com/ARM-auto-place/users/internal/models/errors"
	"github.com/ARM-auto-place/users/internal/storage/postgres"
	externalapi "github.com/ARM-auto-place/users/pkg/interfaces/externalapi"
	customErrors "github.com/ARM-auto-place/users/pkg/models/errors"
	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type usersService struct {
	logger  zerolog.Logger
	storage postgres.UsersStorage
}

func NewUsersService(logger zerolog.Logger, storage postgres.UsersStorage) externalapi.ExternalAPI {
	return &usersService{logger: logger, storage: storage}
}

func (s *usersService) GetProfile(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	profile, err := s.storage.GetProfile(ctx, id)
	return profile, s.storageError(ctx, "get profile", err)
}

func (s *usersService) GetMyCars(ctx context.Context, id uuid.UUID) ([]models.MyCar, error) {
	cars, err := s.storage.GetMyCars(ctx, id)
	return cars, s.storageError(ctx, "get my cars", err)
}

func (s *usersService) GetFavoriteCars(ctx context.Context, id uuid.UUID) ([]models.FavoriteCar, error) {
	cars, err := s.storage.GetFavoriteCars(ctx, id)
	return cars, s.storageError(ctx, "get favorite cars", err)
}

func (s *usersService) UpdateProfile(ctx context.Context, id uuid.UUID, patch models.ProfilePatch) (models.Profile, error) {
	if err := validateProfilePatch(patch); err != nil {
		return models.Profile{}, err
	}
	profile, err := s.storage.UpdateProfile(ctx, id, patch)
	return profile, s.storageError(ctx, "update profile", err)
}

func (s *usersService) UpdateSettings(ctx context.Context, id uuid.UUID, patch models.SettingsPatch) (models.Settings, error) {
	if patch.DarkMode == nil && patch.PrivateData == nil {
		return models.Settings{}, customErrors.BadRequest("at least one setting is required")
	}
	settings, err := s.storage.UpdateSettings(ctx, id, patch)
	return settings, s.storageError(ctx, "update settings", err)
}

func (s *usersService) storageError(ctx context.Context, operation string, err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, commonErrors.ErrProfileNotFound):
		return customErrors.ProfileNotFound()
	default:
		logger := zerolog.Ctx(ctx)
		if logger.GetLevel() == zerolog.Disabled {
			logger = &s.logger
		}
		logger.Error().Str("operation", operation).Type("errorType", err).Msg("storage operation failed")
		return customErrors.InternalServerError(err)
	}
}

func validateProfilePatch(patch models.ProfilePatch) error {
	for _, field := range []struct {
		name  string
		value *string
	}{
		{"name", patch.Name},
		{"phone", patch.Phone},
	} {
		if field.value == nil && patch.Has(field.name) ||
			field.value != nil && strings.TrimSpace(*field.value) == "" {
			return customErrors.BadRequest(field.name + " must be a non-empty string")
		}
	}
	if patch.Age != nil && *patch.Age < 0 {
		return customErrors.BadRequest("age cannot be negative")
	}
	if patch.PhotoPath != nil && utf8.RuneCountInString(*patch.PhotoPath) > 255 {
		return customErrors.BadRequest("photoPath cannot exceed 255 characters")
	}
	if patch.Empty() {
		return customErrors.BadRequest("at least one profile field is required")
	}
	return nil
}
