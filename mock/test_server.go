package main

import (
	"chain/account"
	. "chain/common"
	"chain/common/config"
	"chain/consensus/transaction"
	"encoding/hex"
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
	Data string
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
		From:   BytesToHexString(addr),
		To:     "173E9CAdA8427ECd33ad728A91E6004Fa3eACf63",
		Amount: int64(120000),
	}
	tx.Sign = BytesToHexString(transaction.Sign(&tx, account, wallet))
	tx.Public = BytesToHexString(cpub)
	msg := PingMsg{
		T:    "TX",
		Data: BytesToHexString(transaction.Serialize(tx)),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("%s", err)
	}
	time.Sleep(1 * time.Second)
	log.Printf("[%s]recv data:%v", time.Now().String(), msg)
	log.Printf("[%s]The AquaChain PRC server is running", time.Now().String())
	ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
}

func testWsServer() {
	addr := ":1666"
	srv := &http.Server{Addr: addr}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "The AquaChain PRC server is running\n")
	})
	http.HandleFunc("/sync", Handler)
	log.Printf("ws server listen on: %s", addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	select {}
}

func testBytesToHexString() {
	client, _ := account.Open(config.Parameters.WalletFile)
	account := client.GetAccount()
	wallet := client.GetWallet()
	// cpub := client.Getcpub()
	addr, _ := wallet.AddressBytes(*account)
	fmt.Printf("\n")
	fmt.Printf("addr = %s\n", hex.EncodeToString(addr))
	fmt.Printf("addr = %q\n", addr)
}

func main() {
	// testBytesToHexString()
	testWsServer()
}
