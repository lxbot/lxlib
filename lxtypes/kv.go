package lxtypes

type (
	KV struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}
	StoreEvent struct {
		id string
		KV
	}
)

func NewStoreEvent(event *Event, kv *KV) *StoreEvent {
	return &StoreEvent{
		id: event.ID,
		KV: *kv,
	}
}

func (this *StoreEvent) GetID() string {
	return this.id
}
