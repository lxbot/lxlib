package lxlib

import (
	"encoding/json"

	"github.com/lxbot/lxlib/v2/common"
	"github.com/lxbot/lxlib/v2/lxtypes"
	"github.com/mitchellh/mapstructure"
)

type Store struct {
	common  *common.LxCommon
	eventCh *chan *lxtypes.Event
	getCh   *chan *lxtypes.StoreEvent
	setCh   *chan *lxtypes.StoreEvent
}

func NewStore() (store *Store, getCh *chan *lxtypes.StoreEvent, setCh *chan *lxtypes.StoreEvent) {
	common.TraceLog("(store)", "lxlib.NewStore()", "start")
	defer common.TraceLog("(store)", "lxlib.NewStore()", "end")

	gCh := make(chan *lxtypes.StoreEvent)
	sCh := make(chan *lxtypes.StoreEvent)
	eventCh := make(chan *lxtypes.Event)

	c := common.NewLxCommon()
	store = &Store{
		common:  c,
		eventCh: &eventCh,
		getCh:   &gCh,
		setCh:   &sCh,
	}

	go c.Listen(&eventCh)
	go store.listen()
	store.Raw(lxtypes.NewEvent(lxtypes.ReadyEvent, lxtypes.ReadyEventPayload{
		Mode:     lxtypes.StdIOMode,
		Endpoint: "",
	}))

	return store, &gCh, &sCh
}

func (this *Store) listen() {
	common.TraceLog("(store)", "lxlib.listen()", "start")
	defer common.TraceLog("(store)", "lxlib.listen()", "end")

	for {
		common.TraceLog("(store)", "lxlib.listen()", "waiting event...")

		eventPtr := <-*this.eventCh

		common.TraceLog("(store)", "lxlib.listen()", "event received")

		switch eventPtr.Event {
		case lxtypes.GetStorageEvent:
			common.TraceLog("(store)", "lxlib.listen()", "event received", "type:", eventPtr.Event)
			*this.getCh <- this.eventToStoreEvent(eventPtr)
		case lxtypes.SetStorageEvent:
			common.TraceLog("(store)", "lxlib.listen()", "event received", "type:", eventPtr.Event)
			*this.setCh <- this.eventToStoreEvent(eventPtr)
		default:
			common.TraceLog("(store)", "lxlib.listen()", "unknown event received", "type:", eventPtr.Event)
		}
	}
}

func (this *Store) eventToStoreEvent(eventPtr *lxtypes.Event) *lxtypes.StoreEvent {
	json := eventPtr.Payload.(json.RawMessage)
	payload, err := common.FromJSON(json)
	if err != nil {
		common.ErrorLog(err)
	}
	kv := new(lxtypes.KV)
	if err := mapstructure.WeakDecode(payload, kv); err != nil {
		common.ErrorLog(err)
	}
	return lxtypes.NewStoreEvent(eventPtr, kv)
}

func (this *Store) Raw(event *lxtypes.Event) {
	common.TraceLog("(store)", "lxlib.Raw()", "start")
	defer common.TraceLog("(store)", "lxlib.Raw()", "end")

	go this.common.Send(event)
}

func (this *Store) SendGetResult(event *lxtypes.StoreEvent) {
	common.TraceLog("(store)", "lxlib.SendGetResult()", "start")
	defer common.TraceLog("(store)", "lxlib.SendGetResult()", "end")

	common.TraceLog("(store)", "lxlib.SendGetResult()", "payload:", event.KV)
	ev := lxtypes.NewEvent(lxtypes.GetStorageEvent, event.KV)
	ev.SetID(event.GetID())
	go this.common.Send(ev)
}
