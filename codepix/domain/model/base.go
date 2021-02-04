package model

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

// Base represents base information for various models
type Base struct {
	ID string `json:"id" gorm:"column:id;type:uuid;primary_key" valid:"uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp" valid:"-"`
}