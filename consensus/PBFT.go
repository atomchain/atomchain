package consensus

import (
	"chain/bus"
	"chain/common"
	"chain/net"
	"chain/storage"
	"encoding/binary"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	PBFT_NUMBER_OF_NODE = 3 // shoule 11
	PBFT_NUMBER_OF_BP   = 2 * PBFT_NUMBER_OF_NODE
)

const (
	NONE    = 0
	PREPRAE = 1
	COMMIT  = 2
	DONE    = 3
)

const (
	FALSE = 0
	TRUE  = 1
)

const (
	VIEW_INDEX         = "SYS_VIEW_INDEX"
	HEIGHT             = "SYS_HEIGHT"
	PREPARE_MESSAGE    = "SYS_PREPARE_MESSAGE"
	PREPARE_DONE_COUNT = "SYS_PREPARE_DONE_COUNT"
	PREPARE_DONE       = "SYS_PREPARE_DONE"
	COMMIT_DONE_COUNT  = "SYS_COMMIT_DONE_COUNT"
	COMMIT_DONE        = "SYS_COMMIT_DONE"
)

// 每一个候选区块
type pbft_round struct {
	state int
	view  uint64
	node  *net.Node
}

func encodeUint(num uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, num)
	return b
}

// func Prepare(pround *pbft_round) {
// 	val, _ := pround.db.Get([]byte(PREPARE_DONE_COUNT))
// 	count := decodeUint(val)
// 	if count > (PBFT_NUMBER_OF_NODE-1)*2/3 {
// 		pround.db.Set([]byte(PREPARE_DONE), encodeUint(uint64(TRUE)))
// 	}
// 	pround.db.Set([]byte(COMMIT_DONE), encodeUint(uint64(FALSE)))
// 	pround.db.Set([]byte(COMMIT_DONE_COUNT), encodeUint(uint64(1)))

// }

// func PrePrepare(pround *pbft_round) {
// 	var b bytes.Buffer
// 	enc := gob.NewEncoder(&b)
// 	if err := enc.Encode(pround.candidate_block); err != nil {
// 		panic(err)
// 	}
// 	pround.db.Set([]byte(PREPARE_MESSAGE), b.Bytes())
// 	pround.db.Set([]byte(PREPARE_DONE_COUNT), encodeUint(uint64(1)))
// 	pround.db.Set([]byte(PREPARE_DONE), encodeUint(uint64(FALSE)))
// }

func (p *pbft_round) CheckViewNumber() {
	data := bus.PackViewSync(GetViewNumber())
	p.node.Shub.Broadcast(data)
	scount := 0
	fcount := 0
	N := PBFT_NUMBER_OF_NODE
	shortTicker := time.NewTicker(time.Duration(2) * time.Second)
	for {
		select {
		case vr := <-bus.PIPE.VR:
			if vr.ViewNumber == GetViewNumber() {
				scount += 1
			} else {
				fcount += 1
			}
			break
		case <-shortTicker.C:
			p.node.Shub.Broadcast(data)
		}
		if scount >= N {
			break
		}
	}
	log.Infof("Check ViewNumber Success")
}

func (p *pbft_round) Commit(height uint64, bdata []byte, master bool) uint64 {
	// fmt.Printf("In Commit %d\n", height)
	defer func() {
		bus.STATE.InCommit = false
	}()
	BPs := GetFixedBP()
	bus.STATE.InCommit = true
	bus.STATE.CommitIndex = height
	var data []byte
	if master {
		data = bus.PackCommitSync(height, bdata)
		p.node.Shub.Broadcast(data)
	}
	scount := 0
	fcount := 0
	shortTicker := time.NewTicker(time.Duration(1) * time.Second)
	N := PBFT_NUMBER_OF_NODE
	for {
		select {
		case cr := <-bus.PIPE.CR:
			if !master {
				if cr.Index == height {
					scount = N
				}
			} else {
				ok, i := common.ContainsString(BPs, cr.From)
				if ok && cr.Index == height {
					common.StringSliceRemove(BPs, i)
					scount += 1
				} else if ok {
					common.StringSliceRemove(BPs, i)
					fcount += 1
				}
			}
			break
		case cs := <-bus.PIPE.CS:
			if !master {
				height = cs.Index
			}
			break
		case <-shortTicker.C:
			if master {
				p.node.Shub.Broadcast(data)
			}
			break
		}
		if scount >= N {
			break
		}
	}
	if master {
		fmt.Println("Broadcast CR")
		cr := bus.PackCommitResp(height)
		p.node.Shub.Broadcast(cr)
	}
	bus.STATE.InCommit = false
	// fmt.Printf("Out Of Commit %d\n", height)
	return height
}

func decodeUint(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func getViewIndex(db *storage.DB) uint64 {
	val, _ := db.Get([]byte(VIEW_INDEX))
	return decodeUint(val)
}

// func UpdateView(pround *pbft_round) {
// 	log.Info(fmt.Sprintf("The current view index is %d", getViewIndex(pround.db)))
// }

func CheckView(pround *pbft_round) {
	log.Info("Check if Master is Online")
}

func NewRound(shareState *bus.State, node *net.Node) *pbft_round {
	log.Info(fmt.Sprintf("Entering into an new PBFT consensus round %d", GetViewNumber()))
	pround := pbft_round{
		state: NONE,
		view:  GetViewNumber(),
		node:  node,
	}

	return &pround
}
