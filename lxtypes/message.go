package lxtypes

import (
	"encoding/json"

	"github.com/rs/xid"
)

type (
	Message struct {
		User     User        `json:"user"`
		Room     Room        `json:"room"`
		Contents []Content   `json:"messages"`
		Mode     Mode        `json:"mode"`
		Raw      interface{} `json:"raw"`
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
	Content struct {
		ID          string       `json:"id"`
		Text        string       `json:"text"`
		Attachments []Attachment `json:"attachments"`
	}
	Attachment struct {
		Url         string `json:"url"`
		Description string `json:"description"`
	}
	Mode string
)

const (
	SendMode  Mode = "send"
	ReplyMode      = "reply"
)

func NewLXMessage(msg M) (*Message, error) {
	t, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	m := new(Message)
	if err := json.Unmarshal(t, m); err != nil {
		return nil, err
	}
	return &Message{
		User:     m.User,
		Room:     m.Room,
		Contents: m.Contents,
		Mode:     m.Mode,
		Raw:      m.Raw,
	}, nil
}

func (this *Message) Copy() (*Message, error) {
	m, err := this.ToMap()
	if err != nil {
		return nil, err
	}
	return NewLXMessage(m)
}

func (this *Message) ResetContents() *Message {
	this.Contents = make([]Content, 0)
	return this
}

func (this *Message) AddContent(text string, attachments ...Attachment) *Message {
	var a []Attachment
	if attachments == nil {
		a = make([]Attachment, 0)
	} else {
		a = attachments
	}

	this.Contents = append(this.Contents, Content{
		ID:          xid.New().String(),
		Text:        text,
		Attachments: a,
	})
	return this
}

func (this *Message) Send() *Message {
	this.Mode = SendMode
	return this
}

func (this *Message) Reply() *Message {
	this.Mode = ReplyMode
	return this
}

func (this *Message) ToMap() (M, error) {
	t, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	var r M
	if err := json.Unmarshal(t, &r); err != nil {
		return nil, err
	}
	return r, nil
}
