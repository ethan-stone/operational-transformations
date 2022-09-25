package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;"`
	CreatedAt time.Time `json:"created_at"`
}

func (document *Document) BeforeCreate(tx *gorm.DB) (err error) {
	document.ID = uuid.New()

	return
}
