package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID     uuid.UUID `json:"userId"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	MiddleName *string   `json:"middleName"`
	Surname    *string   `json:"surname"`
	Phone      string    `json:"phone"`
	PhotoPath  *string   `json:"photoPath"`
	Age        *int32    `json:"age"`
	City       *string   `json:"city"`
	Address    *string   `json:"address"`
	Settings   Settings  `json:"settings"`
}

type Settings struct {
	DarkMode    bool `json:"darkMode"`
	PrivateData bool `json:"privateData"`
}

type ProfilePatch struct {
	Name       *string `json:"name,omitempty"`
	MiddleName *string `json:"middleName,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	PhotoPath  *string `json:"photoPath,omitempty"`
	Age        *int32  `json:"age,omitempty"`
	City       *string `json:"city,omitempty"`
	Address    *string `json:"address,omitempty"`
	present    map[string]json.RawMessage
}

func (p *ProfilePatch) UnmarshalJSON(data []byte) error {
	type fields ProfilePatch
	*p = ProfilePatch{}
	return decodePatch(data, (*fields)(p), &p.present,
		"name", "middleName", "surname", "phone", "photoPath", "age", "city", "address")
}

// Has reports whether the JSON included a field, including explicit null.
func (p ProfilePatch) Has(field string) bool {
	_, ok := p.present[field]
	return ok
}

type SettingsPatch struct {
	DarkMode    *bool `json:"darkMode,omitempty"`
	PrivateData *bool `json:"privateData,omitempty"`
}

func (p *SettingsPatch) UnmarshalJSON(data []byte) error {
	type fields SettingsPatch
	*p = SettingsPatch{}
	var present map[string]json.RawMessage
	if err := decodePatch(data, (*fields)(p), &present, "darkMode", "privateData"); err != nil {
		return err
	}
	for field, value := range present {
		if bytes.Equal(value, []byte("null")) {
			return fmt.Errorf("%s must be a boolean", field)
		}
	}
	return nil
}

func decodePatch(data []byte, target any, present *map[string]json.RawMessage, allowed ...string) error {
	if err := json.Unmarshal(data, present); err != nil {
		return err
	}
	if *present == nil {
		return fmt.Errorf("expected a JSON object")
	}
	for field := range *present {
		if !slices.Contains(allowed, field) {
			return fmt.Errorf("unknown field %q", field)
		}
	}
	return json.Unmarshal(data, target)
}

type CarCard struct {
	CarID      uuid.UUID `json:"carId"`
	Mark       string    `json:"mark"`
	Model      string    `json:"model"`
	BuildYear  int32     `json:"buildYear"`
	Price      *int32    `json:"price"`
	Kilometre  *int32    `json:"kilometre"`
	Condition  *string   `json:"condition"`
	Color      *string   `json:"color"`
	CoverImage *string   `json:"coverImage"`
}

type MyCar struct {
	CarCard   `json:",inline"`
	CreatedAt time.Time `json:"createdAt"`
}

type FavoriteCar struct {
	CarCard     `json:",inline"`
	FavoritedAt time.Time `json:"favoritedAt"`
}

type ResponseMeta struct {
	Error            bool           `json:"error"`
	ErrorText        string         `json:"errorText"`
	AdditionalErrors map[string]any `json:"additionalErrors"`
}

type Response struct {
	ResponseMeta `json:",inline"`
	Data         any `json:"data"`
}
