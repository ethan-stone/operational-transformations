package db

import (
	"time"

	"database/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Operation struct {
	ID            uuid.UUID      `json:"id" gorm:"type:char(36);primary_key"`
	DocumentID    uuid.UUID      `json:"document_id" gorm:"type:char(36);primary_key"`
	IsProcessed   bool           `json:"is_processed"`
	Action        string         `json:"action"`         // Should be INSERT or DELETE
	StartPosition int32          `json:"start_position"` // Where to insert into or where to start deleting from
	EndPosition   sql.NullInt32  `json:"end_position"`   // Where to stop deleting from. For INSERT operations this will always be null
	Text          sql.NullString `json:"text"`           // The text to insert for an INSERT operation. This will always be null for a DELETE operation
	CreatedAt     time.Time      `json:"created_at"`
}

func (operation *Operation) BeforeCreate(tx *gorm.DB) (err error) {
	operation.ID = uuid.New()

	return
}
