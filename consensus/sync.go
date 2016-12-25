package consensus

import (
	"chain/bus"
	. "chain/common"
	"chain/crypto"
	"chain/net"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func NeedSyncBlock(mtr string, idx uint64, node *net.Node) bool {
	defer func() { bus.STATE.SyncNeed = false }()
	bus.STATE.SyncNeed = true
	RetFlag := false
	BPs := GetFixedBP()
	N := 3
	BPCount := 0
	NBPCount := 0
	saMsg := bus.PackSyncAsk(mtr, idx)
	node.Shub.Broadcast(saMsg)
	shortTicker := time.NewTicker(time.Duration(2) * time.Second)
	for {
		select {
		case sp := <-bus.PIPE.SP:
			cpub, _ := HexStringToBytes(sp.Pubkey)
			sign, _ := HexStringToBytes(sp.Signature)
			addr := bus.GetAddress(cpub)
			addr = strings.ToLower(addr)
			ok, i := ContainsString(BPs, addr)
			s := strconv.FormatUint(sp.Index, 10) + sp.Root
			hash := crypto.Sha256([]byte(s))
			if ok {
				if bus.Verify(cpub, sign, hash) {
					StringSliceRemove(BPs, i)
					if sp.Index == bus.STATE.Index && sp.Root == bus.STATE.Root {
						BPCount = BPCount + 1
					} else {
						NBPCount = NBPCount + 1
					}
				}
			}
			break
		case <-shortTicker.C:
			log.Info("Collect Infomation Not Complete")
			node.Shub.Broadcast(saMsg)
		}
		if BPCount >= N {
			RetFlag = false
			break
		} else if NBPCount >= N {
			RetFlag = true
			break
		}
	}
	return RetFlag
}

func SyncBlock() {
	select {}
}

func UpdateBP(height uint64) {
	if height%(7*24*60*6) == 0 {
		// for _, vx := range vxs {
		// 	addr, _ := common.HexStringToBytes(vx.To)
		// 	consensus.UpdateChart(vx.Amount, addr) // Only need Update Every 7 * 24 * 60 * 6 Blocks
		// }
	}
}
