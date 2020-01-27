package lxlib

import (
	"github.com/go-test/deep"
	"testing"
)

func dummyMessage() M {
	return M{
		"user": M{
			"id":   "user.id",
			"name": "user.name",
		},
		"room": M{
			"id":          "room.id",
			"name":        "room.name",
			"description": "room.description",
		},
		"message": M{
			"id":   "message.id",
			"text": "message.text",
			"attachments": []M{
				{
					"url":         "message.attachments[0].url",
					"description": "message.attachments[0].description",
				},
				{
					"url":         "message.attachments[1].url",
					"description": "message.attachments[1].description",
				},
			},
		},
		"mode": "",
	}
}

func TestNewLXMessage(t *testing.T) {
	msg := dummyMessage()
	lxm, err := NewLXMessage(msg)
	if err != nil || lxm == nil {
		t.Error("LXMessage create error:", err)
	}
}

func TestLXMessage_ToMap(t *testing.T) {
	msg := dummyMessage()
	lxm, _ := NewLXMessage(msg)

	m, err := lxm.ToMap()
	if err != nil || m == nil {
		t.Error("LXMessage ToMap error:", err)
	}

	t.Logf("%v", msg)
	t.Logf("%v", m)

	if diff := deep.Equal(msg, m); diff != nil {
		t.Error("invalid ToMap:", diff)
	}
}
