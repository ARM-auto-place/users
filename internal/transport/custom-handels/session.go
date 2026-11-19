package custom_handlers

import (
	"fmt"
	"net/url"
	"slices"

	customErrors "github.com/ARM-auto-place/users/pkg/models/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type sessionUserKey struct{}

func Authenticate(signingKey string, allowedOrigins []string) fiber.Handler {
	key := []byte(signingKey)
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired(),
		jwt.WithIssuer("auth"), jwt.WithAudience("users"))
	return func(ctx *fiber.Ctx) error {
		id, err := sessionUser(parser, key, ctx.Cookies("user_id"))
		if err != nil {
			return SendResponse(ctx, nil, customErrors.Unauthorized())
		}
		if ctx.Method() == fiber.MethodPost && !allowedOrigin(ctx, allowedOrigins) {
			return SendResponse(ctx, nil, customErrors.Forbidden())
		}
		ctx.Locals(sessionUserKey{}, id)
		return ctx.Next()
	}
}

func sessionUser(parser *jwt.Parser, key []byte, cookie string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	if _, err := parser.ParseWithClaims(cookie, claims, func(*jwt.Token) (any, error) {
		return key, nil
	}); err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("invalid session user")
	}
	return id, nil
}

func allowedOrigin(ctx *fiber.Ctx, allowedOrigins []string) bool {
	origin := ctx.Get(fiber.HeaderOrigin)
	if origin == "" || slices.Contains(allowedOrigins, origin) {
		return true
	}
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Scheme == ctx.Protocol() && parsed.Host == ctx.Get(fiber.HeaderHost) &&
		parsed.Host != "" && parsed.User == nil && parsed.Path == "" && parsed.RawQuery == "" && parsed.Fragment == ""
}
