package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Value превращает структуру в JSON для базы данных
func (u UserRequest) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan превращает JSON из базы данных обратно в структуру Go
func (u *UserRequest) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, u)
}
