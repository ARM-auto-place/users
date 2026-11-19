package custom_handlers

import (
	"encoding/json"
	"mime"

	externalapi "github.com/ARM-auto-place/users/pkg/interfaces/externalapi"
	customErrors "github.com/ARM-auto-place/users/pkg/models/errors"
	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetProfile(ctx *fiber.Ctx, svc externalapi.ExternalAPI) error {
	return handle(ctx, func(id uuid.UUID) (any, error) {
		return svc.GetProfile(ctx.UserContext(), id)
	})
}

func GetMyCars(ctx *fiber.Ctx, svc externalapi.ExternalAPI) error {
	return handle(ctx, func(id uuid.UUID) (any, error) {
		return svc.GetMyCars(ctx.UserContext(), id)
	})
}

func GetFavoriteCars(ctx *fiber.Ctx, svc externalapi.ExternalAPI) error {
	return handle(ctx, func(id uuid.UUID) (any, error) {
		return svc.GetFavoriteCars(ctx.UserContext(), id)
	})
}

func UpdateProfile(ctx *fiber.Ctx, svc externalapi.ExternalAPI) error {
	return handle(ctx, func(id uuid.UUID) (any, error) {
		var patch models.ProfilePatch
		if err := parseBody(ctx, &patch); err != nil {
			return nil, err
		}
		return svc.UpdateProfile(ctx.UserContext(), id, patch)
	})
}

func UpdateSettings(ctx *fiber.Ctx, svc externalapi.ExternalAPI) error {
	return handle(ctx, func(id uuid.UUID) (any, error) {
		var patch models.SettingsPatch
		if err := parseBody(ctx, &patch); err != nil {
			return nil, err
		}
		return svc.UpdateSettings(ctx.UserContext(), id, patch)
	})
}

func handle(ctx *fiber.Ctx, call func(uuid.UUID) (any, error)) error {
	id, ok := ctx.Locals(sessionUserKey{}).(uuid.UUID)
	if !ok {
		return SendResponse(ctx, nil, customErrors.Unauthorized())
	}
	data, err := call(id)
	if err != nil {
		return SendResponse(ctx, nil, err)
	}
	return SendResponse(ctx, data, nil)
}

func parseBody(ctx *fiber.Ctx, patch any) error {
	mediaType, _, err := mime.ParseMediaType(ctx.Get(fiber.HeaderContentType))
	if err != nil || mediaType != fiber.MIMEApplicationJSON {
		return customErrors.BadRequest("Content-Type must be application/json")
	}
	if err := json.Unmarshal(ctx.Body(), patch); err != nil {
		return customErrors.BadRequest("invalid JSON update")
	}
	return nil
}
