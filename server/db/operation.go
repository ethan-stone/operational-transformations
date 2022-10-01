package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Operation struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	DocumentID    uuid.UUID `json:"document_id" gorm:"type:char(36);primary_key"`
	IsProcessed   bool      `json:"is_processed"`
	Action        string    `json:"action"`         // Should be INSERT or DELETE
	StartPosition int32     `json:"start_position"` // Where to insert into or where to start deleting from
	EndPosition   int32     `json:"end_position"`   // Where to stop deleting from. For INSERT operations this will always be the same as the start position and ignored
	Text          string    `json:"text"`           // The text to insert for an INSERT operation. This will always be an empty string and ignored for DELETE operations
	CreatedAt     time.Time `json:"created_at"`
}

func (operation *Operation) BeforeCreate(tx *gorm.DB) (err error) {
	operation.ID = uuid.New()

	return
}
