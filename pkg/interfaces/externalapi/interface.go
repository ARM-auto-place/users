package externalapi

import (
	"context"

	"github.com/ARM-auto-place/users/swaggers/externalapi/models"
	"github.com/google/uuid"
)

// @tg version=1.0.0
// @tg backend=users-api
// @tg title=`Users API`
// @tg servers=
//
//go:generate tg transport --services . --out ../../../internal/transport/jsonRPC/externalapi --outSwagger ../../../swaggers/externalapi/swagger.yaml
//go:generate sed -i.bak -e "/name: user_id/,/format: uuid/s/format: uuid/format: jwt/" ../../../swaggers/externalapi/swagger.yaml
//go:generate rm -f ../../../internal/transport/jsonRPC/externalapi/metrics.go ../../../swaggers/externalapi/swagger.yaml.bak

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
	// @tg http-cookies=userID|user_id
	// @tg log-skip=userID,response
	// @tg handler=github.com/ARM-auto-place/users/internal/transport/custom-handels:GetProfile
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:ProfileResponse
	// @tg summary=`Профиль текущего пользователя`
	GetProfile(ctx context.Context, userID uuid.UUID) (profile models.Profile, err error)

	// @tg http-method=GET
	// @tg http-path=/v1/car/my
	// @tg http-cookies=userID|user_id
	// @tg log-skip=userID,response
	// @tg handler=github.com/ARM-auto-place/users/internal/transport/custom-handels:GetMyCars
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:MyCarsResponse
	// @tg summary=`Опубликованные автомобили пользователя`
	GetMyCars(ctx context.Context, userID uuid.UUID) (cars []models.MyCar, err error)

	// @tg http-method=GET
	// @tg http-path=/v1/car/favorite
	// @tg http-cookies=userID|user_id
	// @tg log-skip=userID,response
	// @tg handler=github.com/ARM-auto-place/users/internal/transport/custom-handels:GetFavoriteCars
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:FavoriteCarsResponse
	// @tg summary=`Избранные автомобили пользователя`
	GetFavoriteCars(ctx context.Context, userID uuid.UUID) (cars []models.FavoriteCar, err error)

	// @tg http-method=POST
	// @tg http-path=/v1/profile
	// @tg http-cookies=userID|user_id
	// @tg log-skip=userID,patch,ProfilePatch,response
	// @tg handler=github.com/ARM-auto-place/users/internal/transport/custom-handels:UpdateProfile
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:ProfileResponse
	// @tg patch.tags=json:inline
	// @tg summary=`Частичное обновление профиля`
	UpdateProfile(ctx context.Context, userID uuid.UUID, patch models.ProfilePatch) (profile models.Profile, err error)

	// @tg http-method=POST
	// @tg http-path=/v1/settings
	// @tg patch.tags=json:inline
	// @tg http-cookies=userID|user_id
	// @tg log-skip=userID,patch,SettingsPatch,response
	// @tg handler=github.com/ARM-auto-place/users/internal/transport/custom-handels:UpdateSettings
	// @tg 200=github.com/ARM-auto-place/users/swaggers/externalapi/models:SettingsResponse
	// @tg summary=`Частичное обновление настроек`
	UpdateSettings(ctx context.Context, userID uuid.UUID, patch models.SettingsPatch) (settings models.Settings, err error)
}
