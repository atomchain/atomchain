package main

import (
	"chain/net/ws"
	"log"
	"time"
)

func check(e error) {
	if e != nil {
		log.Fatalf("fatal: %s\n", e) // print and os.exit(1)
		// panic(e)
	}
}

func testWsClient() {
	addr := "localhost:8888"
	client, err := ws.NewWsConnection(addr)
	check(err)
	defer client.Colse()
	client.OnData(func(data []byte) {
		log.Printf("on data:[%s], byte length %d", time.Now().String(), len(data))
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				msg := []byte(t.String())
				client.SendMessage(msg)
			}
		}
	}()

	select {}
}

func main() {
	log.SetFlags(0)
	testWsClient()
}
