package net

import (
	"chain/net/ws"
)

// ParserListenser interface
type ParserListenser interface {
	OnData(err error)
	OnError(err error)
}

// Neighbor record neighbor node
type Neighbor struct {
	id       int
	hostname string
	pubkey   string
	status   string
	reader   *ws.WsConnection // client to server
	writer   *ws.WsConnection // server to client
}

// NeighborPool record all neighbor node
// type NeighborPool struct {
// 	pool map[int]Neighbor
// }

func NodeParseFromBytes(data []byte) {}

func BlockParseFromBytes(data []byte) {}

func VoteParseFromBytes(data []byte) {}
