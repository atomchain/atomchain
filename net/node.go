package net

import (
	"chain/account"
	"chain/common/config"
	"chain/net/ws"

	"chain/crypto"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// INIT Constants
const (
	INIT       = 0
	HAND       = 1
	HANDSHAKE  = 2
	HANDSHAKED = 3
	ESTABLISH  = 4
	INACTIVITY = 5
)

var cli *account.ClientImpl
var pool *Pool
var shub *ws.ServerHub

type Msg struct {
	T    string
	Data string
}

// Node to export
type Node struct {
	state     uint32 // node state
	chF       chan func() error
	id        uint64   // The nodes's id
	cap       [32]byte // The node capability set
	version   uint32   // The network protocol the node used
	services  uint64   // The services the node supplied
	relay     bool     // The relay capability of the node (merge into capbility flag)
	height    uint64   // The node latest block height
	publicKey *crypto.PubKey
	Hub       Hub
	Shub      *ws.ServerHub
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func (node *Node) backend() {
	for f := range node.chF {
		f()
	}
}

func server(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := ws.NewClient(shub, conn)
	shub.RegisterClient(&client)
	go client.WritePump()
	go client.ReadPump()
}

// func connect_neighbors(neighbors []string) {
// 	// done := make(chan struct{})
// 	// var wg sync.WaitGroup
// 	// wg.Add(len(neighbors))
// 	for _, d := range neighbors {
// 		u := url.URL{Scheme: "ws", Host: d, Path: "/sync"}
// 		log.Info(fmt.Sprintf("connecting to %s", u.String()))
// 		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 		if err != nil {
// 			log.Warn(fmt.Sprintf("Error When Connecting to %s %s", u.String(), err))
// 			return
// 		}

// 		go func() {
// 			for {
// 				_, message, err := c.ReadMessage()
// 				if err != nil {
// 					log.Println("read:", err)
// 					return
// 				}
// 				var msg Msg
// 				json.Unmarshal(message, &msg)
// 				data, _ := common.HexStringToBytes(msg.Data)
// 				switch msg.T {
// 				case "TX":
// 					tx := transaction.Unserialize(data)
// 					bpub, _ := common.HexStringToBytes(tx.Public)
// 					from, _ := common.HexStringToBytes(tx.From)
// 					pub := account.DePubkey(bpub)
// 					addr := account.PubToAddress(*pub)
// 					if !bytes.Equal(from, addr) {
// 						break
// 					}
// 					sign, _ := common.HexStringToBytes(tx.Sign)
// 					hash := transaction.GetHash(tx)
// 					if !transaction.Verify(bpub, sign, hash) {
// 						break
// 					}
// 					fmt.Printf("received tx: %+v\n", tx)
// 					PushTx(pool, message)
// 					break
// 				}
// 			}
// 		}()
// 	}
// 	// wg.Wait()
// }

func startServer(addr string) *http.Server {
	srv := &http.Server{Addr: addr}
	log.Infof("start server: %s", addr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "The AquaChain PRC server is running\n")
	})
	http.HandleFunc("/sync", server)
	// Quit := make(chan bool)
	go func() {
		log.Infof("start ws server listen on: %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Httpserver: ListenAndServe() error: %s", err)
		}
		// ok := <-Quit
	}()
	return srv
}

//RetryConnAddrs struct
type RetryConnAddrs struct {
	sync.RWMutex
	RetryAddrs map[string]int
}

//ConnectingNodes struct
type ConnectingNodes struct {
	sync.RWMutex
	ConnectingAddrs []string
}

func rmNode(node *Node) {
	log.Debug(fmt.Sprintf("Remove unused/deuplicate node: 0x%0x", node.id))
}

func GetHashHexString(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func (node *Node) bootstrap(seed string, bnodes []string) {
	params := config.Parameters
	log.Infof("Query the seed string%s (HashHex:%s)", params.SeedString, GetHashHexString(params.SeedString))
	// var NgList []string
	// NgList = make([]string, 100)
	// seedHash := GetHashHexString(params.SeedString)
	// for _, bnode := range bnodes {
	// 	queryURL := bnode + "/" + seedHash
	// 	resp, err := http.Get(queryURL)
	// 	if err != nil {
	// 		log.Warnf("dht query error, queryURL:%s", queryURL)
	// 	} else {
	// 		defer resp.Body.Close()
	// 		body, _ := ioutil.ReadAll(resp.Body)
	// 		err = json.Unmarshal(body, &NgList)
	// 	}
	// }
	// go connect_neighbors(NgList)
	go func() {
		seedHash := GetHashHexString(params.SeedString)
		node.Hub = CreateConnections(seedHash)
	}()
}

// NewNode create node
func NewNode(bnodes []string) *Node {
	params := config.Parameters
	n := Node{
		state: INIT,
		chF:   make(chan func() error),
	}
	// runtime.SetFinalizer(&n, rmNode)
	log.Infof("params.ServerAddr:%s", params.ServerAddr)
	shub = ws.NewServerHub()
	n.Shub = shub
	startServer(params.ServerAddr)
	n.bootstrap(params.SeedString, bnodes)
	return &n
}

func SetPool(p *Pool) {
	pool = p
}

func SetClient(client *account.ClientImpl) {
	cli = client
}
