package dht

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	cid "github.com/ipfs/go-cid"
	iaddr "github.com/ipfs/go-ipfs-addr"
	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	inet "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	mh "github.com/multiformats/go-multihash"
)

var bootstrapPeers = []string{}

var rendezvous = "atomchain"

func handleStream(stream inet.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}

func discovery() {
	rendezvousString := flag.String("r", rendezvous, "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.Parse()

	ctx := context.Background()

	host, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}

	host.SetStreamHandler("/atomlchain/0.0.1", handleStream)

	kadDht, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	for _, peerAddr := range bootstrapPeers {
		addr, _ := iaddr.ParseString(peerAddr)
		peerinfo, _ := pstore.InfoFromP2pAddr(addr.Multiaddr())

		if err := host.Connect(ctx, *peerinfo); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Connection established with bootstrap node: ", *peerinfo)
		}
	}

	v1b := cid.V1Builder{Codec: cid.Raw, MhType: mh.SHA2_256}
	rendezvousPoint, _ := v1b.Sum([]byte(*rendezvousString))

	fmt.Println("announcing ourselves...")
	tctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := kadDht.Provide(tctx, rendezvousPoint, true); err != nil {
		panic(err)
	}

	fmt.Println("searching for other peers...")
	tctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	peers, err := kadDht.FindProviders(tctx, rendezvousPoint)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d peers!\n", len(peers))

	for _, p := range peers {
		fmt.Println("Peer: ", p)
	}

	for _, p := range peers {
		if p.ID == host.ID() || len(p.Addrs) == 0 {
			// No sense connecting to ourselves or if addrs are not available
			continue
		}
		print(p.ID)
		stream, err := host.NewStream(ctx, p.ID, "/atomlchain/0.0.1")

		if err != nil {
			fmt.Println(err)
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}

		fmt.Println("Connected to: ", p)
	}

	select {}
}
