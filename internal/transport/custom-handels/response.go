package custom_handlers

import (
	"errors"
	"time"

	customErrors "github.com/ARM-auto-place/users/pkg/models/errors"
	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func SendResponse(ctx *fiber.Ctx, data any, err error) error {
	status := fiber.StatusOK
	response := models.Response{Data: data}
	if err != nil {
		var public *customErrors.Error
		if !errors.As(err, &public) {
			public = customErrors.InternalServerError(err)
		}
		status = public.GetStatusCode()
		response.Error = true
		response.ErrorText = public.GetPublicText()
		response.AdditionalErrors = map[string]any{"reason": public.Error()}
	}
	ctx.Set(fiber.HeaderCacheControl, "no-store")
	return ctx.Status(status).JSON(response)
}

func LogRequests(ctx *fiber.Ctx) error {
	start := time.Now()
	err := ctx.Next()
	status := ctx.Response().StatusCode()
	if err != nil {
		status = fiber.StatusInternalServerError
		var httpError *fiber.Error
		if errors.As(err, &httpError) {
			status = httpError.Code
		}
	}
	level := zerolog.InfoLevel
	if status >= 500 {
		level = zerolog.ErrorLevel
	}
	zerolog.Ctx(ctx.UserContext()).WithLevel(level).
		Str("method", ctx.Method()).
		Str("path", ctx.Route().Path).
		Int("status", status).
		Dur("duration", time.Since(start)).
		Msg("http request")
	return err
}
