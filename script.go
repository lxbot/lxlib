package lxlib

import (
	"encoding/json"

	"github.com/lxbot/lxlib/v2/common"
	"github.com/lxbot/lxlib/v2/lxtypes"
	"github.com/mitchellh/mapstructure"
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
	common.TraceLog("(script)", "lxlib.NewScript()", "start")
	defer common.TraceLog("(script)", "lxlib.NewScript()", "end")

	messageCh := make(chan *lxtypes.Message)
	eventCh := make(chan *lxtypes.Event)

	c := common.NewLxCommon()
	script := &Script{
		common:    c,
		events:    make(map[string]Event),
		eventCh:   &eventCh,
		messageCh: &messageCh,
	}

	go c.Listen(&eventCh)
	go script.listen()
	script.Raw(lxtypes.NewEvent(lxtypes.ReadyEvent, lxtypes.ReadyEventPayload{
		Mode:     lxtypes.StdIOMode,
		Endpoint: "",
	}))

	return script, &messageCh
}

func (this *Script) listen() {
	common.TraceLog("(script)", "lxlib.listen()", "start")
	defer common.TraceLog("(script)", "lxlib.listen()", "end")

	for {
		common.TraceLog("(script)", "lxlib.listen()", "waiting event...")

		eventPtr := <-*this.eventCh

		switch eventPtr.Event {
		case lxtypes.IncomingMessageEvent:
			json := eventPtr.Payload.(json.RawMessage)
			common.TraceLog("(script)", "lxlib.listen()", "event received", "type:", eventPtr.Event)
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
			common.TraceLog("(script)", "lxlib.listen()", "event received", "type:", eventPtr.Event)
			if event, ok := this.events[eventPtr.ID]; ok {
				common.TraceLog("(script)", "lxlib.listen()", "found registered event", "id:", eventPtr.ID)
				*event.eventCh <- eventPtr
			}
		default:
			common.TraceLog("(script)", "lxlib.listen()", "unknown event received", "type:", eventPtr.Event)
		}
	}
}

func (this *Script) Raw(event *lxtypes.Event) {
	common.TraceLog("(script)", "lxlib.Raw()", "start")
	defer common.TraceLog("(script)", "lxlib.Raw()", "end")

	go this.common.Send(event)
}

func (this *Script) SendMessage(message *lxtypes.Message) error {
	common.TraceLog("(script)", "lxlib.SendMessage()", "start")
	defer common.TraceLog("(script)", "lxlib.SendMessage()", "end")

	m, err := message.ToMap()
	if err != nil {
		return err
	}
	common.TraceLog("(script)", "lxlib.SendMessage()", "payload:", m)

	this.Raw(lxtypes.NewEvent(lxtypes.OutgoingMessageEvent, m))
	return nil
}

func (this *Script) GetStorage(key string) interface{} {
	common.TraceLog("(script)", "lxlib.GetStorage()", "start")
	defer common.TraceLog("(script)", "lxlib.GetStorage()", "end")

	event := lxtypes.NewEvent(lxtypes.GetStorageEvent, lxtypes.KV{
		Key: key,
	})
	common.TraceLog("(script)", "lxlib.GetStorage()", "payload:", event.Payload)

	eventCh := make(chan *lxtypes.Event)
	this.events[event.ID] = Event{
		eventCh: &eventCh,
	}
	defer delete(this.events, event.ID)
	this.Raw(event)

	common.TraceLog("(script)", "lxlib.GetStorage()", "waiting response...")

	result := <-eventCh
	json := result.Payload.(json.RawMessage)
	payload, err := common.FromJSON(json)
	if err != nil {
		common.WarnLog(err)
	}
	resultKV := new(lxtypes.KV)
	if err := mapstructure.WeakDecode(payload, resultKV); err != nil {
		common.WarnLog(err)
		resultKV = &lxtypes.KV{
			Key:   key,
			Value: nil,
		}
	}

	common.TraceLog("(script)", "lxlib.GetStorage()", "response received", "result:", resultKV)

	return resultKV.Value
}

func (this *Script) SetStorage(key string, value interface{}) {
	common.TraceLog("(script)", "lxlib.SetStorage()", "start")
	defer common.TraceLog("(script)", "lxlib.SetStorage()", "end")

	this.Raw(lxtypes.NewEvent(lxtypes.SetStorageEvent, lxtypes.KV{
		Key:   key,
		Value: value,
	}))
}

func (this *Script) Close() {
	common.TraceLog("(script)", "lxlib.Close()", "start")
	defer common.TraceLog("(script)", "lxlib.Close()", "end")

	go this.common.Close()
}
