package lxlib

import (
	"encoding/json"

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
	common.TraceLog("lxlib.NewScript()", "start")
	defer common.TraceLog("lxlib.NewScript()", "end")

	messageCh := make(chan *lxtypes.Message)
	eventCh := make(chan *lxtypes.Event)

	c := common.NewLxCommon()
	script := &Script{
		common:    c,
		events:    make(map[string]Event),
		eventCh:   &eventCh,
		messageCh: &messageCh,
	}

	go script.listen()
	script.Raw(lxtypes.NewEvent(lxtypes.ReadyEvent, nil))

	return script, &messageCh
}

func (this *Script) listen() {
	common.TraceLog("lxlib.listen()", "start")
	defer common.TraceLog("lxlib.listen()", "end")

	go this.common.Listen(this.eventCh)

	for {
		common.TraceLog("lxlib.listen()", "waiting event...")

		eventPtr := <-*this.eventCh

		common.TraceLog("lxlib.listen()", "event received")
		switch eventPtr.Event {
		case lxtypes.IncomingMessageEvent:
			json := eventPtr.Payload.(json.RawMessage)
			common.TraceLog("lxlib.listen()", "event received", "type:", eventPtr.Event, "json:", json)
			payload, err := common.FromJSON(json)
			if err != nil {
				common.ErrorLog(err)
				continue
			}
			message, err := lxtypes.NewLXMessage(payload)
			if err != nil {
				common.ErrorLog(err)
				continue
			}
			*this.messageCh <- message
		case lxtypes.GetStorageEvent:
			common.TraceLog("lxlib.listen()", "event received", "type:", eventPtr.Event)
			if event, ok := this.events[eventPtr.ID]; ok {
				common.TraceLog("lxlib.listen()", "found registered event", "id:", eventPtr.ID)
				*event.eventCh <- eventPtr
			}
		default:
			common.TraceLog("lxlib.listen()", "unknown event received", "type:", eventPtr.Event)
		}
	}
}

func (this *Script) Raw(event *lxtypes.Event) {
	common.TraceLog("lxlib.Raw()", "start")
	defer common.TraceLog("lxlib.Raw()", "end")

	go this.common.Send(event)
}

func (this *Script) SendMessage(message *lxtypes.Message) error {
	common.TraceLog("lxlib.SendMessage()", "start")
	defer common.TraceLog("lxlib.SendMessage()", "end")

	m, err := message.ToMap()
	if err != nil {
		return err
	}
	common.TraceLog("lxlib.SendMessage()", "payload:", m)

	this.Raw(lxtypes.NewEvent(lxtypes.OutgoingMessageEvent, m))
	return nil
}

func (this *Script) GetStorage(key string) interface{} {
	common.TraceLog("lxlib.GetStorage()", "start")
	defer common.TraceLog("lxlib.GetStorage()", "end")

	event := lxtypes.NewEvent(lxtypes.GetStorageEvent, lxtypes.KV{
		Key: key,
	})
	common.TraceLog("lxlib.GetStorage()", "payload:", event.Payload)

	eventCh := make(chan *lxtypes.Event)
	this.events[event.ID] = Event{
		eventCh: &eventCh,
	}
	this.Raw(event)

	common.TraceLog("lxlib.GetStorage()", "waiting response...")

	result := <-eventCh
	resultKV := result.Payload.(lxtypes.KV)

	common.TraceLog("lxlib.GetStorage()", "response received", "result:", resultKV)

	return resultKV.Value
}

func (this *Script) SetStorage(key string, value interface{}) {
	common.TraceLog("lxlib.SetStorage()", "start")
	defer common.TraceLog("lxlib.SetStorage()", "end")

	this.Raw(lxtypes.NewEvent(lxtypes.SetStorageEvent, lxtypes.KV{
		Key:   key,
		Value: value,
	}))
}

func (this *Script) Close() {
	common.TraceLog("lxlib.Close()", "start")
	defer common.TraceLog("lxlib.Close()", "end")

	go this.common.Close()
}
