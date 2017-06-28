package app

type ClientEvent struct {
	ID    uint32 `json:"id"`
	IsOut bool   `json:"is_out"`
}

type CuratorEvent struct {
	LoserIds []uint32 `json:"loser_ids"`
}

func NewClientEvent() *ClientEvent {
	return &ClientEvent{}
}

func NewCuratorEvent() *CuratorEvent {
	return &CuratorEvent{}
}
