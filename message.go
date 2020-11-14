package lxlib

import (
	"encoding/json"
	"errors"
)

type (
	M         = map[string]interface{}
	LXMessage struct {
		User    User    `json:"user"`
		Room    Room    `json:"room"`
		Message Message `json:"message"`
		Mode    string  `json:"mode"`
		Raw     interface{} `json:"raw"`
	}
	pack struct {
		User    User    `json:"user"`
		Room    Room    `json:"room"`
		Message Message `json:"message"`
		Raw     interface{} `json:"raw"`
	}
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	Room struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	Message struct {
		ID          string       `json:"id"`
		Text        string       `json:"text"`
		Attachments []Attachment `json:"attachments"`
	}
	Attachment struct {
		Url         string `json:"url"`
		Description string `json:"description"`
	}
)

func NewLXMessage(msg M) (*LXMessage, error) {
	t, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	pack := new(pack)
	if err := json.Unmarshal(t, pack); err != nil {
		return nil, err
	}
	return &LXMessage{
		User: pack.User,
		Room: pack.Room,
		Message: pack.Message,
		Mode: "",
		Raw: pack.Raw,
	}, nil
}

func (this *LXMessage) SetText(text string) *LXMessage {
	this.Message.Text = text
	return this
}

func (this *LXMessage) Send() *LXMessage {
	this.Mode = "send"
	return this
}

func (this *LXMessage) Reply() *LXMessage {
	this.Mode = "reply"
	return this
}

func (this *LXMessage) ToMap() (M, error) {
	r, err := toMap(this)
	if err != nil {
		return nil, err
	}

	r["message"].(M)["attachments"], err = toArrayMap(this.Message.Attachments)
	return r, err
}

func toMap(i interface{}) (M, error) {
	t, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var r M
	if err := json.Unmarshal(t, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func toArrayMap(t interface{}) ([]M, error) {
	switch t.(type) {
	case []Attachment:
		s := t.([]Attachment)
		r := make([]map[string]interface{}, len(s))
		var err error
		for i, v := range s {
			r[i], err = toMap(v)
			if err != nil {
				return nil, err
			}
		}
		return r, nil
	default:
		return nil, errors.New("invalid type")
	}
}
