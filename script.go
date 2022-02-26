package lxlib

import (
	"github.com/lxbot/lxlib/v2/common"
	"github.com/lxbot/lxlib/v2/lxtypes"
)

type (
	Script struct {
		common    *common.LxCommon
		events    map[string]Event
		eventCh   *chan *lxtypes.Event
		messageCh *chan *lxtypes.Message
	}
	Event struct {
		eventCh *chan *lxtypes.Event
	}
)

func NewScript() (*Script, *chan *lxtypes.Message) {
	messageCh := make(chan *lxtypes.Message)
	eventCh := make(chan *lxtypes.Event)

	common := common.NewLxCommon()
	script := &Script{
		common:    common,
		events:    make(map[string]Event),
		eventCh:   &eventCh,
		messageCh: &messageCh,
	}

	go common.Listen(&eventCh)
	go script.listen()
	script.Raw(lxtypes.NewEvent(lxtypes.ReadyEvent, nil))

	return script, &messageCh
}

func (this *Script) listen() {
	for {
		eventPtr := <-*this.eventCh
		switch eventPtr.Event {
		case lxtypes.IncomingMessageEvent:
			*this.messageCh <- eventPtr.Payload.(*lxtypes.Message)
		case lxtypes.GetStorageEvent:
			if event, ok := this.events[eventPtr.ID]; ok {
				*event.eventCh <- eventPtr
			}
		}
	}
}

func (this *Script) Raw(event *lxtypes.Event) {
	go this.common.Send(event)
}

func (this *Script) SendMessage(message *lxtypes.Message) error {
	m, err := message.ToMap()
	if err != nil {
		return err
	}
	this.Raw(lxtypes.NewEvent(lxtypes.OutgoingMessageEvent, m))
	return nil
}

func (this *Script) GetStorage(key string) interface{} {
	event := lxtypes.NewEvent(lxtypes.GetStorageEvent, lxtypes.KV{
		Key: key,
	})
	eventCh := make(chan *lxtypes.Event)
	this.events[event.ID] = Event{
		eventCh: &eventCh,
	}
	this.Raw(event)

	result := <-eventCh
	resultKV := result.Payload.(lxtypes.KV)
	return resultKV.Value
}

func (this *Script) SetStorage(key string, value interface{}) {
	this.Raw(lxtypes.NewEvent(lxtypes.SetStorageEvent, lxtypes.KV{
		Key:   key,
		Value: value,
	}))
}

func (this *Script) Close() {
	go this.common.Close()
}
