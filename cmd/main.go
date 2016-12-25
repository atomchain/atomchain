package main

import (
	"chain/account"
	"chain/bus"
	"chain/common"
	"chain/common/config"
	"chain/common/password"
	"chain/consensus"
	"chain/consensus/block"
	"chain/crypto"
	"chain/helper"
	"chain/net"
	"chain/net/ws"
	"chain/params"
	"chain/storage"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	mlog "chain/common/log"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var iflag bool
var pool *net.Pool = net.NewPool()
var myselfAddress string

func parse(v string) string {
	return "transaction"
}

func InitData(filepath string) (error, *storage.DB) {
	rdb := storage.NewRocksDB(filepath)
	db := storage.New(rdb)
	return nil, db
}

func seedUInt64(d uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, d)
	return b
}

func waitBlockFromMaster(master string, height uint64) *block.Block {
	log.Info("#waitBlockFromMaster")
	message := net.GetTxs(pool)
	// var msg net.Msg
	// json.Unmarshal(message[0], &msg)
	b := block.NewBlock(height, message)
	return &b
}

func consensus_loop(node *net.Node, db *storage.DB) {
	fmt.Println("start consensus_loop")
	// log.Infof("The current BP is: %", consensus.GetFixedBP())
	pbftround := consensus.NewRound(&bus.STATE, node)
	// ticker := time.NewTicker(time.Duration(config.Parameters.Period) * time.Second)
	// pbftround.CheckViewNumber()
	for {
		log.Infof("The current height is %d", bus.STATE.Index)
		// pbftround.CheckViewNumber()
		master := consensus.GetNextBFTMaster(bus.STATE.Index)
		height := bus.STATE.Index + 1
		if myselfAddress == master[2:] {
			log.Infof("Run to be the master %s \n", myselfAddress)
			bus.STATE.Index = pbftround.Commit(height, nil, true)
		} else {
			bus.STATE.Index = pbftround.Commit(height, nil, false)
		}
	}
	// BPs := consensus.GetBP()
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			consensus.CheckView(pbft_round)
	// 			consensus.UpdateView(pbft_round)
	// 			log.Infof("Current Tx Pool Size is: %d \n", net.PoolSize(pool))
	// 			maddr := consensus.GetNextMaster(BPs, seedUInt64(lastBlock.Index()))
	// 			log.Infof("The master address is %s", maddr)
	// 			if myselfAddress == maddr {
	// 				log.Infof("I am the master %s", myselfAddress)
	// 				newBlock := block.NewBlock(lastBlock.Index()+1, net.GetTxs(pool)) // 打包交易
	// 				bdata, _ := newBlock.Marshal()
	// 				node.Hub.Broadcast(bdata) // 广播给其他的节点并受到回复
	// 				net.EmptyTxs(pool)
	// 				lastBlock = &newBlock
	// 			} else {
	// 				// Wait for OtherBlock
	// 				log.Infof("I am  not the master %s", myselfAddress)
	// 				lastBlock = waitBlockFromMaster(maddr, lastBlock.Index()+1)
	// 			}
	// 		}
	// 	}
	// }()
}

func GenerateTestBlock() *block.Block {
	data := helper.GetTrans()
	genesis := block.NewBlock(uint64(0), data)
	genesis.Hash()
	return &genesis
}

func waitEnoughtNodes(atomNode *net.Node) {
	ticker := time.NewTicker(time.Duration(config.Parameters.Period) * time.Second)
	for {
		tflag := false
		select {
		case <-ticker.C:
			if atomNode.Shub.GetClientsSize() >= consensus.PBFT_NUMBER_OF_NODE && atomNode.Hub.GetConnectedCons() >= consensus.PBFT_NUMBER_OF_NODE {
				tflag = true
			}
		}
		if tflag {
			break
		}
	}
}

func wait() {
	ticker := time.NewTicker(time.Duration(config.Parameters.Period) * time.Second)
	for {
		select {
		case <-ticker.C:
		}
	}
}

func initByArgs(arg string) {
	if strings.Contains(arg, "1") {
		iflag = true
		config.Init2("./mock/config-1.json")
	} else if strings.Contains(arg, "2") {
		config.Init2("./mock/config-2.json")
	} else if strings.Contains(arg, "3") {
		config.Init2("./mock/config-3.json")
	}
}

/*--------------------------*/

func main() {
	app := cli.NewApp()
	app.Name = "atomchain CLI"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "start a test network",
			Action: func(c *cli.Context) error {
				initByArgs(c.Args().Get(0))
				log.Infof("account.Open %s\n", config.Parameters.WalletFile)
				client, err := account.Open(config.Parameters.WalletFile)
				myselfAddress = client.GetAddress()
				log.Infof("Client Wallet Address is %s\n", myselfAddress)
				if err != nil {
					return err
				}
				net.SetClient(client)
				net.SetPool(pool)
				_, atomData := InitData(config.Parameters.atomData)
				atomNode := net.NewNode(params.TestnetBootnodes[:])
				if atomData.GetBlockIndex("BLK") == uint64(0) {
					genesis := GenerateTestBlock()
					atomData.SaveBlock(genesis, "BLK")
				}
				mtr := common.BytesToHexString(atomData.GetBlockMT("BLK"))
				bidx := atomData.GetBlockIndex("BLK")
				bus.SetupPipe(bidx, mtr, client)
				fmt.Printf("The current Merkle Root is %s\n", mtr)
				waitEnoughtNodes(atomNode)
				log.Info("Connected to Enought Nodes")
				// wait()
				if consensus.NeedSyncBlock(mtr, bidx, atomNode) {
					consensus.SyncBlock()
				}
				fmt.Print("Sync is compeleted")

				consensus_loop(atomNode, atomData)
				stop := make(chan os.Signal)
				select {
				case signal := <-stop:
					log.Info("Got signal:%v\n", signal)
				}
				return nil
			},
		},
		{
			Name:    "dev",
			Aliases: []string{"t"},
			Usage:   "",
			Action: func(c *cli.Context) error {
				ws.NewWsConnection("127.0.0.1:8080")
				// u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/sync"}
				// co, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				// if err != nil {
				// 	log.Fatal("dial:", err)
				// }
				// go func() {
				// 	for {
				// 		_, message, err := co.ReadMessage()
				// 		if err != nil {
				// 			log.Println("read:", err)
				// 			return
				// 		}
				// 		log.Printf("recv: %s", message)
				// 	}
				// }()
				// select {}
				return nil
			},
		},
		{
			Name:    "account",
			Aliases: []string{"a"},
			Usage:   "Manager your wallet account",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create a new account",
					Action: func(c *cli.Context) error {
						pass, err := password.GetConfirmedPassword()
						if err != nil {
							return err
						}
						_, err = account.Create("wallet.dat", pass)
						if err != nil {
							return err
						}
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						os.Remove("wallet.dat")
						return nil
					},
				},
			},
		},
		{
			Name:    "tx_send",
			Aliases: []string{"tx"},
			Usage:   "send a transaction block",
			Action: func(c *cli.Context) error {
				// testCliArgs(c)
				toAddr := c.Args().Get(0)
				amount, err := strconv.ParseInt(c.Args().Get(1), 10, 64)
				if err != nil {
					panic(err)
				}
				seedHash := net.GetHashHexString(params.SeedString)
				hub := net.CreateConnections(seedHash)
				data := net.BuildMessage([]byte(toAddr), amount)
				hub.Boardcast(data)
				return nil
			},
		},
		{
			Name:    "vt_send",
			Aliases: []string{"vt"},
			Usage:   "send a vote block",
			Action: func(c *cli.Context) error {
				toAddr := c.Args().Get(0)
				amount, err := strconv.ParseInt(c.Args().Get(1), 10, 64)
				if err != nil {
					panic(err)
				}
				seedHash := net.GetHashHexString(params.SeedString)
				hub := net.CreateConnections(seedHash)
				data := net.BuildMessage([]byte(toAddr), amount)
				hub.Boardcast(data)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	crypto.SetAlg()
	mlog.Init()
	rand.Seed(time.Now().UnixNano())
	// rpc.Bootstrap()
}
