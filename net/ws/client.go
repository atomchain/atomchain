package ws

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// WsConnection is websocket connection
// methods:
//		NewWsConnection(addr),
//  	OnData(callback: func([]byte))
// 		SendMessage([]byte)
//		Close()
type WsConnection struct {
	ready bool
	// 连接地址
	Addr string
	// websocket 连接器
	ws *RecConn
	// 发送信息的缓冲 channel
	send chan []byte
	// 接受信息的缓冲 channel
	recv chan []byte
	// 失败重连的触发 channel
	reconnect chan struct{}
}

// NewWsConnection create websocket connection
func NewWsConnection(addr string) (*WsConnection, error) {
	// ctx, _ := context.WithCancel(context.Background())
	msgSend := make(chan []byte, 20000)
	msgRecv := make(chan []byte, 20000)

	connection := &WsConnection{
		ready:     false,
		Addr:      addr,
		send:      msgSend,
		recv:      msgRecv,
		reconnect: make(chan struct{}),
		ws:        nil,
	}
	ws := RecConn{}
	connection.ws = &ws
	u := url.URL{Scheme: "ws", Host: connection.Addr, Path: "/sync"}
	ws.Dial(u.String(), nil)
	// Read Goro
	go func() {
		for {
			if !ws.IsConnected() {
				connection.ready = false
				log.Printf("Websocket disconnected %s", ws.GetURL())
				time.Sleep(2 * time.Second)
				continue
			}
			connection.ready = true
			_, message, err := ws.Conn.ReadMessage()
			if err != nil {
				log.Printf("Error: ReadMessage %s", ws.GetURL())
				ws.CloseAndReconnect()
			}
			msgRecv <- message
		}
	}()

	//Write Goro
	go func() {
		for {

			if !ws.IsConnected() {
				log.Printf("Websocket disconnected %s", ws.GetURL())
				time.Sleep(2 * time.Second)
				continue
			}

			select {
			case message := <-msgSend:
				if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Printf("Error: WriteMessage %s", ws.GetURL())
				}
			}
		}

	}()

	return connection, nil
}

// SendMessage send data to web socket
func (connection *WsConnection) SendMessage(msg []byte) {
	if connection.ready {
		connection.send <- msg
		return
	}
	log.Printf("connection is not ready")
}

// OnData process data by callback
func (connection *WsConnection) OnData(dealData func([]byte) []byte) {
	go func() {
		for {
			select {
			case data := <-connection.recv:
				// fmt.Print(data)
				toSend := dealData(data)
				if toSend != nil {
					connection.SendMessage(toSend)
				}
			}
		}
	}()
}

func (connection *WsConnection) IsReady() bool {
	return connection.ready
}

func (connection *WsConnection) Colse() {
	connection.ws.Close()
}
