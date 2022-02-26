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
	gCh := make(chan *lxtypes.KV)
	sCh := make(chan *lxtypes.KV)
	eventCh := make(chan *lxtypes.Event)

	common := common.NewLxCommon()
	store = &Store{
		common:  common,
		eventCh: &eventCh,
		getCh:   &gCh,
		setCh:   &sCh,
	}

	go common.Listen(&eventCh)
	go store.listen()
	store.Raw(lxtypes.NewEvent(lxtypes.ReadyEvent, nil))

	return store, &gCh, &sCh
}

func (this *Store) listen() {
	for {
		eventPtr := <-*this.eventCh
		switch eventPtr.Event {
		case lxtypes.GetStorageEvent:
			*this.getCh <- eventPtr.Payload.(*lxtypes.KV)
		case lxtypes.SetStorageEvent:
			*this.setCh <- eventPtr.Payload.(*lxtypes.KV)
		}
	}
}

func (this *Store) Raw(event *lxtypes.Event) {
	go this.common.Send(event)
}

func (this *Store) SendGetResult(kv *lxtypes.KV) {
	go this.common.Send(lxtypes.NewEvent(lxtypes.GetStorageEvent, kv))
}
