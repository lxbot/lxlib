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
		"raw": "raw_message",
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
	msg["message"].(M)["attachments"] = []M{}

	t.Logf("%v", msg)
	t.Logf("%v", m)

	if diff := deep.Equal(msg, m); diff != nil {
		t.Error("invalid ToMap:", diff)
	}
}

func TestLXMessage_ToMapWithAttachment(t *testing.T) {
	msg := dummyMessage()
	lxm, _ := NewLXMessage(msg)

	m, err := lxm.SetAttachments(lxm.Message.Attachments).ToMap()
	if err != nil || m == nil {
		t.Error("LXMessage ToMap error:", err)
	}

	t.Logf("%v", msg)
	t.Logf("%v", m)

	if diff := deep.Equal(msg, m); diff != nil {
		t.Error("invalid ToMap:", diff)
	}
}

func TestLXMessage_ToSliceMap(t *testing.T) {
	msg := dummyMessage()
	lxm, _ := NewLXMessage(msg)

	m, err := lxm.ToSliceMap()
	if err != nil || m == nil {
		t.Error("LXMessage ToSliceMap error:", err)
	}
	msg["message"].(M)["attachments"] = []M{}

	t.Logf("%v", msg)
	t.Logf("%v", m)

	if len(m) != 1 {
		t.Error("invalid ToSliceMap length:", len(m))
	} else if diff := deep.Equal(msg, m[0]); diff != nil {
		t.Error("invalid ToSliceMap:", diff)
	}
}

func TestLXMessage_ToSliceMapWithAttachment(t *testing.T) {
	msg := dummyMessage()
	lxm, _ := NewLXMessage(msg)

	m, err := lxm.SetAttachments(lxm.Message.Attachments).ToSliceMap()
	if err != nil || m == nil {
		t.Error("LXMessage ToSliceMap error:", err)
	}

	t.Logf("%v", msg)
	t.Logf("%v", m)

	if len(m) != 1 {
		t.Error("invalid ToSliceMap length:", len(m))
	} else if diff := deep.Equal(msg, m[0]); diff != nil {
		t.Error("invalid ToSliceMap:", diff)
	}
}
