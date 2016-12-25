package helper

import (
	"chain/account"
	. "chain/common"
	"chain/common/config"
	"chain/consensus/transaction"
	"chain/net"
	"encoding/json"
)

type PingMsg net.PingMsg

func GetTrans() [][]byte {
	var transactions [][]byte
	client, _ := account.Open(config.Parameters.WalletFile)
	account := client.GetAccount()
	wallet := client.GetWallet()
	cpub := client.Getcpub()
	addr, _ := wallet.AddressBytes(*account)

	tx := transaction.Transfer{
		From:   BytesToHexString(addr),
		To:     "56088521d9eca9fdb5f7bae94a13f183c3246d3d",
		Amount: int64(120000),
	}
	tx.Sign = BytesToHexString(transaction.Sign(&tx, account, wallet))
	tx.Public = BytesToHexString(cpub)

	tx2 := transaction.Transfer{
		From:   BytesToHexString(addr),
		To:     "c97be97e26cd532724cdd2ede3af8548bfac3b68",
		Amount: int64(120000),
	}
	tx2.Sign = BytesToHexString(transaction.Sign(&tx2, account, wallet))
	tx2.Public = BytesToHexString(cpub)

	vx := transaction.Vote{
		From:   BytesToHexString(addr),
		To:     "56088521d9eca9fdb5f7bae94a13f183c3246d3d",
		Amount: int64(120000),
	}
	vx.Sign = BytesToHexString(transaction.VSign(&vx, account, wallet))
	vx.Public = BytesToHexString(cpub)

	vx2 := transaction.Vote{
		From:   BytesToHexString(addr),
		To:     "c97be97e26cd532724cdd2ede3af8548bfac3b68",
		Amount: int64(120000),
	}
	vx2.Sign = BytesToHexString(transaction.VSign(&vx2, account, wallet))
	vx2.Public = BytesToHexString(cpub)

	vx3 := transaction.Vote{
		From:   BytesToHexString(addr),
		To:     "c97be97e26cd532724cdd2ede3af8548bfac3b68",
		Amount: int64(120000),
	}
	vx3.Sign = BytesToHexString(transaction.VSign(&vx2, account, wallet))
	vx3.Public = BytesToHexString(cpub)

	msg := PingMsg{
		T:    "TX",
		Data: BytesToHexString(transaction.Serialize(tx)),
	}
	data, _ := json.Marshal(msg)
	transactions = append(transactions, data)

	msg = PingMsg{
		T:    "TX",
		Data: BytesToHexString(transaction.Serialize(tx2)),
	}
	data, _ = json.Marshal(msg)
	transactions = append(transactions, data)

	msg = PingMsg{
		T:    "VX",
		Data: BytesToHexString(transaction.VSerialize(vx)),
	}
	data, _ = json.Marshal(msg)
	transactions = append(transactions, data)

	msg = PingMsg{
		T:    "VX",
		Data: BytesToHexString(transaction.VSerialize(vx2)),
	}
	data, _ = json.Marshal(msg)
	transactions = append(transactions, data)

	msg = PingMsg{
		T:    "VX",
		Data: BytesToHexString(transaction.VSerialize(vx3)),
	}
	data, _ = json.Marshal(msg)
	transactions = append(transactions, data)
	return transactions
}
