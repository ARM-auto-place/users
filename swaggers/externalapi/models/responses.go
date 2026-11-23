package models

type ProfileResponse struct {
	ResponseMeta `json:",inline"`
	Data         Profile `json:"data"`
}

type SettingsResponse struct {
	ResponseMeta `json:",inline"`
	Data         Settings `json:"data"`
}

type MyCarsResponse struct {
	ResponseMeta `json:",inline"`
	Data         []MyCar `json:"data"`
}

type FavoriteCarsResponse struct {
	ResponseMeta `json:",inline"`
	Data         []FavoriteCar `json:"data"`
}
