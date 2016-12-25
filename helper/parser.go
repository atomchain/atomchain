package helper

import (
	"chain/common"
	"chain/consensus/transaction"
	"encoding/json"
)

type Msg struct {
	T    string
	Data string
}

func GetX(datas [][]byte) ([]transaction.Transfer, []transaction.Vote) {
	var txs []transaction.Transfer
	var vxs []transaction.Vote
	for _, message := range datas {
		var msg Msg
		json.Unmarshal(message, &msg)
		data, _ := common.HexStringToBytes(msg.Data)
		switch msg.T {
		case "TX":
			tx := transaction.Unserialize(data)
			txs = append(txs, *tx)
			// fmt.Printf("TX %+v\n", tx)
			break
		case "VX":
			vx := transaction.VUnserialize(data)
			vxs = append(vxs, *vx)
			// fmt.Printf("VX %+v\n", vx)
			break
		}
	}
	return txs, vxs
}
