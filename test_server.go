package main

import (
	"chain/account"
	"chain/consensus/transaction"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type PingMsg struct {
	T    string
	Data []byte
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

// Handler is http upgrade to Websocket Handler
func Handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	client, _ := account.Open("wallet.dat")
	account := client.GetAccount()
	wallet := client.GetWallet()
	cpub := client.Getcpub()
	addr, _ := wallet.AddressBytes(*account)
	tx := transaction.Transfer{
		From:   addr,
		To:     addr,
		Amount: int64(120000),
	}
	tx.Sign = transaction.Sign(&tx, account, wallet)
	tx.Public = cpub
	msg := PingMsg{
		T:    "TX",
		Data: transaction.Serialize(tx),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("%s", err)
	}
	time.Sleep(1 * time.Second)
	err = ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
}

func testWsServer() {
	addr := "127.0.0.1:1222"
	srv := &http.Server{Addr: addr}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "The AquaChain PRC server is running\n")
	})
	http.HandleFunc("/sync", Handler)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	log.Printf("ws server listen on: %s", addr)
	select {}
}

func main() {
	testWsServer()
}
