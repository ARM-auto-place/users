package postgres

import _ "embed"

//go:embed sql/select_profile.sql
var sqlSelectProfile string

//go:embed sql/update_profile.sql
var sqlUpdateProfile string

//go:embed sql/lock_profile.sql
var sqlLockProfile string

//go:embed sql/upsert_settings.sql
var sqlUpsertSettings string

//go:embed sql/select_my_cars.sql
var sqlSelectMyCars string

//go:embed sql/select_favorite_cars.sql
var sqlSelectFavoriteCars string
