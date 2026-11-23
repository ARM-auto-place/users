// @tg version=1.0.0
// @tg backend=users-api
// @tg title=`Users API`
// @tg servers=/
package externalapi

import (
	"context"

	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/google/uuid"
)

//go:generate tg transport --services . --out ../../../internal/transport/jsonRPC/externalapi --outSwagger ../../../swaggers/externalapi/swagger.yaml
//go:generate rm -f ../../../internal/transport/jsonRPC/externalapi/metrics.go

// ExternalAPI serves the current user's profile, settings, and car lists.
//
// @tg http-server log
// @tg http-prefix=/api/users
// @tg 400=github.com/ARM-auto-place/users/swaggers/externalapi/models:Response
// @tg 401=github.com/ARM-auto-place/users/swaggers/externalapi/models:Response
// @tg 403=github.com/ARM-auto-place/users/swaggers/externalapi/models:Response
// @tg 404=github.com/ARM-auto-place/users/swaggers/externalapi/models:Response
// @tg 500=github.com/ARM-auto-place/users/swaggers/externalapi/models:Response
type ExternalAPI interface {
	// @tg http-method=GET
	// @tg http-path=/v1/profile
	// @tg http-headers=userID|X-User-Id
	// @tg userID.tags=json:-
	// @tg log-skip=userID,response
	// @tg http-response=github.com/ARM-auto-place/users/internal/transport/custom-handels:GetProfile
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:ProfileResponse
	// @tg summary=`Профиль текущего пользователя`
	GetProfile(ctx context.Context, userID uuid.UUID) (profile models.Profile, err error)

	// @tg http-method=GET
	// @tg http-path=/v1/car/my
	// @tg http-headers=userID|X-User-Id
	// @tg userID.tags=json:-
	// @tg log-skip=userID,response
	// @tg http-response=github.com/ARM-auto-place/users/internal/transport/custom-handels:GetMyCars
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:MyCarsResponse
	// @tg summary=`Опубликованные автомобили пользователя`
	GetMyCars(ctx context.Context, userID uuid.UUID) (cars []models.MyCar, err error)

	// @tg http-method=GET
	// @tg http-path=/v1/car/favorite
	// @tg http-headers=userID|X-User-Id
	// @tg userID.tags=json:-
	// @tg log-skip=userID,response
	// @tg http-response=github.com/ARM-auto-place/users/internal/transport/custom-handels:GetFavoriteCars
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:FavoriteCarsResponse
	// @tg summary=`Избранные автомобили пользователя`
	GetFavoriteCars(ctx context.Context, userID uuid.UUID) (cars []models.FavoriteCar, err error)

	// @tg http-method=POST
	// @tg http-path=/v1/profile
	// @tg http-headers=userID|X-User-Id
	// @tg userID.tags=json:-
	// @tg log-skip=userID,profilePatch,ProfilePatch,response
	// @tg http-response=github.com/ARM-auto-place/users/internal/transport/custom-handels:UpdateProfile
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:ProfileResponse
	// @tg profilePatch.tags=json:inline
	// @tg summary=`Частичное обновление профиля`
	UpdateProfile(ctx context.Context, userID uuid.UUID, profilePatch models.ProfilePatch) (profile models.Profile, err error)

	// @tg http-method=POST
	// @tg http-path=/v1/settings
	// @tg settingsPatch.tags=json:inline
	// @tg http-headers=userID|X-User-Id
	// @tg userID.tags=json:-
	// @tg log-skip=userID,settingsPatch,SettingsPatch,response
	// @tg http-response=github.com/ARM-auto-place/users/internal/transport/custom-handels:UpdateSettings
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:SettingsResponse
	// @tg summary=`Частичное обновление настроек`
	UpdateSettings(ctx context.Context, userID uuid.UUID, settingsPatch models.SettingsPatch) (settings models.Settings, err error)
}
