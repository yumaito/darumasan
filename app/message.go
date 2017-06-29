package app

type GameMessage struct {
	ID          string   `json:"id"`
	ClientType  uint32   `json:"client_type"`
	Clients     []string `json:"clients"`
	DeadClients []string `json:"dead_clients"`
	CuratorID   string   `json:"curator_id"`
	IsWatched   bool     `json:"is_watched"`
}

type ClientMessage struct {
	ID         string `json:"id"`
	ClientType uint32 `json:"client_type"`
	Status     bool   `json:"status"`
}
