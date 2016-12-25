package net

import (
	"chain/account"
	"chain/bus"
	"chain/common"
	. "chain/common"
	"chain/common/config"
	"chain/consensus/transaction"
	"chain/net/ws"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	dhtAPI = "http://40.85.147.33:5000/atomchain"
)

type PingMsg struct {
	T    string
	Data string
}

type Hub struct {
	cons []*ws.WsConnection
}

func (hub *Hub) Boardcast(data []byte) {
	for _, c := range hub.cons {
		log.Printf("send data to %s", c.Addr)
		c.SendMessage(data)
	}
}

func (hub *Hub) GetConnectedCons() int {
	count := 0
	for index := 0; index < len(hub.cons); index++ {
		if hub.cons[index].IsReady() {
			count += 1
		}
	}
	return count
}

func (hub *Hub) Broadcast(data []byte) {
	for _, c := range hub.cons {
		log.Printf("boardcast to %s", c.Addr)
		c.SendMessage(data)
	}
}

func makeURL(signature string) string {
	return fmt.Sprintf("%s/%s", dhtAPI, signature)
}

func isJSONByte(s []byte) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
}

func parseAddrList(body []byte) ([]string, error) {
	var addrList []string
	// if !isJSONByte(body) {
	// 	log.Printf("response is not an valid json, body:%s\n", string(body))
	// 	return addrList, nil
	// }
	err := json.Unmarshal(body, &addrList)
	if err != nil {
		return []string{}, nil
	}
	return addrList, nil
}

//过滤掉本机的ip地址
func filterLocalAddr(addrList []string) []string {
	localAddr := getOutboundIP() + config.Parameters.ServerAddr
	var nlist []string
	for _, addr := range addrList {
		if !strings.HasPrefix(addr, localAddr) {
			nlist = append(nlist, addr)
		}
	}
	return nlist
}

func getAddrList(url string) ([]string, error) {
	// return []string{"127.0.0.1:8080"}, nil
	return []string{"127.0.0.1:1666", "127.0.0.1:1668", "127.0.0.1:1670"}, nil
	resp, err := http.Get(url)
	if err != nil {
		// log.Printf("url=%s\nget http request failed\n", url)
		return nil, errors.New("get DHT request failed")
	}
	log.Infof("url=%s\nget DHT request success\n", url)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read http body error")
	}
	addrList, _ := parseAddrList(body)
	return filterLocalAddr(addrList), nil
}

func BuildMessage(toAddr []byte, amount int64) []byte {
	client, _ := account.Open(config.Parameters.WalletFile)
	account := client.GetAccount()
	wallet := client.GetWallet()
	cpub := client.Getcpub()
	fromAddr, _ := wallet.AddressBytes(*account)
	tx := transaction.Transfer{
		From:   BytesToHexString(fromAddr),
		To:     BytesToHexString(toAddr),
		Amount: amount,
	}
	tx.Sign = BytesToHexString(transaction.Sign(&tx, account, wallet))
	tx.Public = BytesToHexString(cpub)
	fmt.Printf("tx={\nFrom:%s\nTo:%s\nAmount:%d\nPublic:%s\nSign:%s\nTxid:%s\n}\n",
		tx.From, tx.To, tx.Amount, tx.Public, tx.Sign, tx.Txid)
	msg := PingMsg{
		T:    "TX",
		Data: BytesToHexString(transaction.Serialize(tx)),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("%s", err)
	}
	return data
}

func CreateConnections(signature string) Hub {
	addrList, err := getAddrList(makeURL(signature))
	if err != nil {
		log.Println(err)
	}
	hub := Hub{}
	log.Printf("createConnections %d neighbors [%v]", len(addrList), addrList)
	for _, addr := range addrList {
		log.Printf("connect to addr=%s\n", addr)

		client, err := ws.NewWsConnection(addr)
		if err != nil {
			log.Printf("connect to addr=%s failed\n", addr)
			continue
		}

		client.OnData(dataCallback)
		hub.cons = append(hub.cons, client)
	}
	return hub
}

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	log.Printf("self.OutboundIP=%s\n", localAddr.IP)
	return fmt.Sprintf("%s", localAddr.IP)
}

func dataCallback(message []byte) []byte {
	var msg Msg
	var ret []byte
	ret = nil
	json.Unmarshal(message, &msg)
	data, _ := common.HexStringToBytes(msg.Data)
	switch msg.T {
	case "TX":
		tx := bus.DealAndVerifyTx(data)
		if tx != nil {
			fmt.Printf("received tx: %+v\n", tx)
		}
		break
	case "SA":
		sa := bus.DealSyncAsk(data)
		if sa != nil {
			fmt.Printf("received Sync Ask: %+v\n", sa)
			ret = bus.PackSyncResp(bus.STATE.Root, bus.STATE.Index)
		}
		break
	case "SP":
		sp := bus.DealSyncResp(data)
		// if sp != nil {
		// 	fmt.Printf("received Sync Ask: %+v\n", sp)
		// }
		bus.PIPE.SP <- *sp
		break
	case "VS":
		vs := bus.DealViewSync(data)
		if vs != nil {
			ret = bus.PackViewResp(uint64(3))
		}
		// bus.PIPE.SP <- *sp
		break
	case "CS":
		cs := bus.DealCommitSync(data)
		if cs != nil && bus.STATE.InCommit == true && bus.STATE.CommitIndex == cs.Index {
			bus.PIPE.CS <- *cs
			ret = bus.PackCommitResp(cs.Index)
		}
		break
	case "CR":
		cr := bus.DealCommitResp(data)
		// fmt.Printf("client received CR %+v\n", cr)
		if cr != nil {
			bus.PIPE.CR <- *cr
		}
		break
	}
	return ret
	// log.Printf("on data:[%s], byte length %d", time.Now().String(), len(message))
	// var msg Msg
	// json.Unmarshal(message, &msg)
	// data, _ := common.HexStringToBytes(msg.Data)
	// switch msg.T {
	// case "TX":
	// 	tx := transaction.Unserialize(data)
	// 	bpub, _ := common.HexStringToBytes(tx.Public)
	// 	from, _ := common.HexStringToBytes(tx.From)
	// 	pub := account.DePubkey(bpub)
	// 	addr := account.PubToAddress(*pub)
	// 	if !bytes.Equal(from, addr) {
	// 		break
	// 	}
	// 	sign, _ := common.HexStringToBytes(tx.Sign)
	// 	hash := transaction.GetHash(tx)
	// 	if !transaction.Verify(bpub, sign, hash) {
	// 		break
	// 	}
	// 	fmt.Printf("received tx: %+v\n", tx)
	// 	PushTx(pool, message)
	// 	break

	// }
}
