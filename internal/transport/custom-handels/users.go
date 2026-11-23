package custom_handlers

import (
	externalapi "github.com/ARM-auto-place/users/pkg/interfaces/externalapi"
	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetProfile(ctx *fiber.Ctx, svc externalapi.ExternalAPI, id uuid.UUID) error {
	return handle(ctx, func() (any, error) {
		return svc.GetProfile(ctx.UserContext(), id)
	})
}

func GetMyCars(ctx *fiber.Ctx, svc externalapi.ExternalAPI, id uuid.UUID) error {
	return handle(ctx, func() (any, error) {
		return svc.GetMyCars(ctx.UserContext(), id)
	})
}

func GetFavoriteCars(ctx *fiber.Ctx, svc externalapi.ExternalAPI, id uuid.UUID) error {
	return handle(ctx, func() (any, error) {
		return svc.GetFavoriteCars(ctx.UserContext(), id)
	})
}

func UpdateProfile(ctx *fiber.Ctx, svc externalapi.ExternalAPI, id uuid.UUID, patch models.ProfilePatch) error {
	return handle(ctx, func() (any, error) {
		return svc.UpdateProfile(ctx.UserContext(), id, patch)
	})
}

func UpdateSettings(ctx *fiber.Ctx, svc externalapi.ExternalAPI, id uuid.UUID, patch models.SettingsPatch) error {
	return handle(ctx, func() (any, error) {
		return svc.UpdateSettings(ctx.UserContext(), id, patch)
	})
}

func handle(ctx *fiber.Ctx, call func() (any, error)) error {
	data, err := call()
	if err != nil {
		return SendResponse(ctx, nil, err)
	}
	return SendResponse(ctx, data, nil)
}
