package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// --- Методы для UserRequest ---

// Value подготавливает структуру для сохранения в JSONB
func (u UserRequest) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan считывает данные из JSONB обратно в структуру
func (u *UserRequest) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed for UserRequest")
	}
	return json.Unmarshal(b, u)
}

// --- Методы для Chapters ---

// Value превращает весь список глав в JSON для БД
func (c Chapters) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

// Scan превращает JSON из БД обратно в слайс глав
func (c *Chapters) Scan(value interface{}) error {
	if value == nil {
		*c = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed for Chapters")
	}
	return json.Unmarshal(b, c)
}
