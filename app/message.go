package app

type Message struct {
	Type uint32
	Msg  []byte
}

type ClientEvent struct {
	ID    uint32 `json:"id"`
	IsOut bool   `json:"is_out"`
}

type CuratorEvent struct {
	LoserIds []uint32 `json:"loser_ids"`
}

func NewMessage(t uint32, msg []byte) *Message {
	return &Message{
		Type: t,
		Msg:  msg,
	}
}

func NewClientEvent() *ClientEvent {
	return &ClientEvent{}
}

func NewCuratorEvent() *CuratorEvent {
	return &CuratorEvent{}
}
