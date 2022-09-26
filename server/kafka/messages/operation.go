package messages

import (
	"time"

	"github.com/google/uuid"
)

type OperationCreatedMsgData struct {
	ID          uuid.UUID `json:"id"`
	DocumentID  uuid.UUID `json:"document_id"`
	IsProcessed bool      `json:"is_processed"`
	CreatedAt   time.Time `json:"created_at"`
}

type OperationCreatedMsg struct {
	ID   uuid.UUID               `json:"id"`
	Type string                  `json:"type"`
	Data OperationCreatedMsgData `json:"data"`
}

func NewOperationCreatedMsg(data OperationCreatedMsgData) *OperationCreatedMsg {
	msg := new(OperationCreatedMsg)
	msg.ID = uuid.New()
	msg.Type = "operation.created"
	msg.Data = data

	return msg
}
