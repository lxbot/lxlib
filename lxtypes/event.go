package lxtypes

import (
	"encoding/json"

	"github.com/rs/xid"
)

type (
	Event struct {
		ID      string      `json:"id"`
		Event   EventType   `json:"event"`
		Payload interface{} `json:"payload"`
	}
	ReadyEventPayload struct {
		Mode     IOMode `json:"mode"`
		Endpoint string `json:"endpoint"`
	}
	IOMode        string
	StdInOutEvent struct {
		ID      string          `json:"id"`
		Event   EventType       `json:"event"`
		Payload json.RawMessage `json:"payload"`
	}
	EventType string
)

const (
	StdIOMode IOMode = "stdio"
)

const (
	ReadyEvent           EventType = "ready"
	OutgoingEvent                  = "outgoning"
	CloseEvent                     = "close"
	IncomingMessageEvent           = "incoming_message"
	OutgoingMessageEvent           = "outgoing_message"
	GetStorageEvent                = "get_storage"
	SetStorageEvent                = "set_storage"
)

func NewEvent(event EventType, payload interface{}) *Event {
	return &Event{
		ID:      xid.New().String(),
		Event:   event,
		Payload: payload,
	}
}

func (this *Event) CopyID(copyFrom *Event) {
	this.ID = copyFrom.ID
}

func (this *Event) PayloadAsMap() M {
	return this.Payload.(M)
}
