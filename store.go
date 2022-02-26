package lxlib

import (
	"github.com/lxbot/lxlib/v2/common"
	"github.com/lxbot/lxlib/v2/lxtypes"
)

type Store struct {
	common  *common.LxCommon
	eventCh *chan *lxtypes.Event
	getCh   *chan *lxtypes.KV
	setCh   *chan *lxtypes.KV
}

func NewStore() (store *Store, getCh *chan *lxtypes.KV, setCh *chan *lxtypes.KV) {
	common.TraceLog("(store)", "lxlib.NewStore()", "start")
	defer common.TraceLog("(store)", "lxlib.NewStore()", "end")

	gCh := make(chan *lxtypes.KV)
	sCh := make(chan *lxtypes.KV)
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
	store.Raw(lxtypes.NewEvent(lxtypes.ReadyEvent, nil))

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
			*this.getCh <- eventPtr.Payload.(*lxtypes.KV)
		case lxtypes.SetStorageEvent:
			common.TraceLog("(store)", "lxlib.listen()", "event received", "type:", eventPtr.Event)
			*this.setCh <- eventPtr.Payload.(*lxtypes.KV)
		default:
			common.TraceLog("(store)", "lxlib.listen()", "unknown event received", "type:", eventPtr.Event)
		}
	}
}

func (this *Store) Raw(event *lxtypes.Event) {
	common.TraceLog("(store)", "lxlib.Raw()", "start")
	defer common.TraceLog("(store)", "lxlib.Raw()", "end")

	go this.common.Send(event)
}

func (this *Store) SendGetResult(kv *lxtypes.KV) {
	common.TraceLog("(store)", "lxlib.SendGetResult()", "start")
	defer common.TraceLog("(store)", "lxlib.SendGetResult()", "end")

	common.TraceLog("(store)", "lxlib.SendGetResult()", "payload:", kv)
	go this.common.Send(lxtypes.NewEvent(lxtypes.GetStorageEvent, kv))
}
