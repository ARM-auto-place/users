package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	commonErrors "github.com/ARM-auto-place/users/internal/models/errors"
	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const defaultTimeout = 5 * time.Second

type storage struct {
	users ConnectManager
	cars  ConnectManager
}

func NewUsersStorage(users, cars ConnectManager) UsersStorage {
	return &storage{users: users, cars: cars}
}

func (s *storage) GetProfile(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	row := s.users.GetMasterConn().QueryRow(ctx, sqlSelectProfile, id)
	return scanProfile(row, "get profile")
}

func (s *storage) UpdateProfile(ctx context.Context, id uuid.UUID, patch models.ProfilePatch) (models.Profile, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	row := s.users.GetMasterConn().QueryRow(ctx, sqlUpdateProfile,
		id, patch.Name,
		patch.Has("middleName") || patch.MiddleName != nil, patch.MiddleName,
		patch.Has("surname") || patch.Surname != nil, patch.Surname,
		patch.Phone,
		patch.Has("photoPath") || patch.PhotoPath != nil, patch.PhotoPath,
		patch.Has("age") || patch.Age != nil, patch.Age,
		patch.Has("city") || patch.City != nil, patch.City,
		patch.Has("address") || patch.Address != nil, patch.Address,
	)
	return scanProfile(row, "update profile")
}

func scanProfile(row pgx.Row, operation string) (profile models.Profile, err error) {
	err = row.Scan(
		&profile.UserID, &profile.Email, &profile.Name, &profile.MiddleName, &profile.Surname, &profile.Phone,
		&profile.PhotoPath, &profile.Age, &profile.City, &profile.Address,
		&profile.Settings.DarkMode, &profile.Settings.PrivateData,
	)
	if err != nil {
		return models.Profile{}, profileError(operation, err)
	}
	return profile, nil
}

func (s *storage) UpdateSettings(ctx context.Context, id uuid.UUID, patch models.SettingsPatch) (settings models.Settings, err error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err = pgx.BeginTxFunc(ctx, s.users.GetMasterConn(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(tx pgx.Tx) error {
		// Lock the profile so concurrent settings writes for this user run in order.
		var profileID uuid.UUID
		if err := tx.QueryRow(ctx, sqlLockProfile, id).Scan(&profileID); err != nil {
			return err
		}
		return tx.QueryRow(ctx, sqlUpsertSettings, profileID, patch.DarkMode, patch.PrivateData).
			Scan(&settings.DarkMode, &settings.PrivateData)
	})
	if err != nil {
		return models.Settings{}, profileError("update settings", err)
	}
	return settings, nil
}

func profileError(operation string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return commonErrors.ErrProfileNotFound
	}
	return fmt.Errorf("%s: postgres: %w", operation, err)
}

func (s *storage) GetMyCars(ctx context.Context, id uuid.UUID) ([]models.MyCar, error) {
	return listCars(ctx, s, id, sqlSelectMyCars, "get my cars", func(card models.CarCard, createdAt time.Time) models.MyCar {
		return models.MyCar{CarCard: card, CreatedAt: createdAt}
	})
}

func (s *storage) GetFavoriteCars(ctx context.Context, id uuid.UUID) ([]models.FavoriteCar, error) {
	return listCars(ctx, s, id, sqlSelectFavoriteCars, "get favorite cars", func(card models.CarCard, createdAt time.Time) models.FavoriteCar {
		return models.FavoriteCar{CarCard: card, FavoritedAt: createdAt}
	})
}

func listCars[T any](ctx context.Context, s *storage, id uuid.UUID, query, operation string, item func(models.CarCard, time.Time) T) ([]T, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	rows, err := s.cars.GetReplicaConn().Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: postgres: %w", operation, err)
	}
	defer rows.Close()

	cars := make([]T, 0)
	for rows.Next() {
		var card models.CarCard
		var image json.RawMessage
		var createdAt time.Time
		err = rows.Scan(
			&card.CarID, &card.Mark, &card.Model, &card.BuildYear, &card.Price, &card.Kilometre,
			&card.Condition, &card.Color, &image, &createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan car: %w", operation, err)
		}
		card.CoverImage = coverImage(image)
		cars = append(cars, item(card, createdAt.UTC()))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: read cars: %w", operation, err)
	}
	return cars, nil
}

func coverImage(raw json.RawMessage) *string {
	var path string
	if json.Unmarshal(raw, &path) == nil && path != "" {
		return &path
	}
	var image struct {
		URL  string `json:"url"`
		Path string `json:"path"`
	}
	if json.Unmarshal(raw, &image) == nil {
		if image.URL != "" {
			return &image.URL
		}
		if image.Path != "" {
			return &image.Path
		}
	}
	return nil
}
