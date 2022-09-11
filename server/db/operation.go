package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Operation struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	DocumentID uuid.UUID `json:"document_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (operation *Operation) BeforeCreate(tx *gorm.DB) (err error) {
	operation.ID = uuid.New()

	return
}
