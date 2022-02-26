package lxlib

import (
	"github.com/lxbot/lxlib/v2/common"
	"github.com/lxbot/lxlib/v2/lxtypes"
)

type Adapter struct {
	common    *common.LxCommon
	eventCh   *chan *lxtypes.Event
	messageCh *chan *lxtypes.Message
}

func NewAdapter() (*Adapter, *chan *lxtypes.Message) {
	messageCh := make(chan *lxtypes.Message)
	eventCh := make(chan *lxtypes.Event)

	c := common.NewLxCommon()
	adapter := &Adapter{
		common:    c,
		eventCh:   &eventCh,
		messageCh: &messageCh,
	}

	go c.Listen(&eventCh)
	go adapter.listen()
	adapter.Raw(lxtypes.NewEvent(lxtypes.ReadyEvent, nil))

	return adapter, &messageCh
}

func (this *Adapter) listen() {
	for {
		eventPtr := <-*this.eventCh
		switch eventPtr.Event {
		case lxtypes.OutgoingMessageEvent:
			*this.messageCh <- eventPtr.Payload.(*lxtypes.Message)
		}
	}
}

func (this *Adapter) Raw(event *lxtypes.Event) {
	go this.common.Send(event)
}

func (this *Adapter) Send(message *lxtypes.Message) {
	go this.common.Send(lxtypes.NewEvent(lxtypes.IncomingMessageEvent, message))
}
