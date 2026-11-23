package models

import (
	"encoding/json"
	"fmt"
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
	if err := json.Unmarshal(data, &p.present); err != nil {
		return err
	}
	if p.present == nil {
		return fmt.Errorf("expected a JSON object")
	}
	for field := range p.present {
		switch field {
		case "name", "middleName", "surname", "phone", "photoPath", "age", "city", "address":
		default:
			return fmt.Errorf("unknown field %q", field)
		}
	}
	return json.Unmarshal(data, (*fields)(p))
}

func (p ProfilePatch) MarshalJSON() ([]byte, error) {
	type fields ProfilePatch
	data, err := json.Marshal(fields(p))
	if err != nil || len(p.present) == 0 {
		return data, err
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	for field := range p.present {
		if _, ok := values[field]; !ok {
			values[field] = json.RawMessage("null")
		}
	}
	return json.Marshal(values)
}

func (p ProfilePatch) Has(field string) bool {
	_, ok := p.present[field]
	return ok
}

func (p ProfilePatch) Empty() bool {
	return len(p.present) == 0 && p.Name == nil && p.MiddleName == nil && p.Surname == nil &&
		p.Phone == nil && p.PhotoPath == nil && p.Age == nil && p.City == nil && p.Address == nil
}

type SettingsPatch struct {
	DarkMode    *bool `json:"darkMode,omitempty"`
	PrivateData *bool `json:"privateData,omitempty"`
}

func (p *SettingsPatch) UnmarshalJSON(data []byte) error {
	var fields map[string]*bool
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if fields == nil {
		return fmt.Errorf("expected a JSON object")
	}
	for field, value := range fields {
		if field != "darkMode" && field != "privateData" {
			return fmt.Errorf("unknown field %q", field)
		}
		if value == nil {
			return fmt.Errorf("%s must be a boolean", field)
		}
	}
	*p = SettingsPatch{DarkMode: fields["darkMode"], PrivateData: fields["privateData"]}
	return nil
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
