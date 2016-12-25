package ws

import (
	"chain/bus"
	"chain/common"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Msg struct {
	T    string
	Data string
}

type ServerHub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *ServerHub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewClient(h *ServerHub, c *websocket.Conn) Client {
	return Client{
		hub:  h,
		conn: c,
		send: make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg Msg
		json.Unmarshal(message, &msg)
		data, _ := common.HexStringToBytes(msg.Data)
		switch msg.T {
		case "TX":
			tx := bus.DealAndVerifyTx(data)
			if tx != nil {
			}
			break
		case "SA":
			sa := bus.DealSyncAsk(data)
			if sa != nil {
			}
			break
		case "SP":
			sp := bus.DealSyncResp(data)
			// fmt.Printf("server received: %+v\n", sp)
			if sp != nil && bus.STATE.SyncNeed {
				bus.PIPE.SP <- *sp
			}
			break
		case "VR":
			vr := bus.DealViewResp(data)
			if vr != nil {
				bus.PIPE.VR <- *vr
			}
			break
		case "CR":
			cr := bus.DealCommitResp(data)
			if cr != nil {
				// fmt.Printf("server received: %+v\n", cr)
				bus.PIPE.CR <- *cr
			}
			break
		}
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewServerHub() *ServerHub {
	shub := ServerHub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	go shub.run()
	return &shub
}

func (h *ServerHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *ServerHub) RegisterClient(c *Client) {
	h.register <- c
}

func (h *ServerHub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

func (h *ServerHub) GetClientsSize() int {
	return len(h.clients)
}
