package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Operation struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	DocumentID  uuid.UUID `json:"document_id" gorm:"type:char(36);primary_key"`
	IsProcessed bool      `json:"is_processed"`
	CreatedAt   time.Time `json:"created_at"`
}

func (operation *Operation) BeforeCreate(tx *gorm.DB) (err error) {
	operation.ID = uuid.New()

	return
}
