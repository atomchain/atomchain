package bus

import (
	"chain/account"

	"github.com/ethereum/go-ethereum/crypto"
)

var PIPE Pipe
var STATE State

type State struct {
	Index       uint64
	Root        string
	Client      *account.ClientImpl
	SyncNeed    bool
	InCommit    bool
	CommitIndex uint64
}

type Pipe struct {
	SP chan SyncResp
	SA chan SyncAsk
	VR chan ViewResp
	CR chan CommitResp
	CS chan CommitSync
}

func SetupPipe(index uint64, root string, client *account.ClientImpl) {
	PIPE = Pipe{
		SP: make(chan SyncResp, 1000),
		SA: make(chan SyncAsk, 1000),
		VR: make(chan ViewResp, 1000),
		CR: make(chan CommitResp, 1000),
		CS: make(chan CommitSync, 1000),
	}
	STATE = State{
		Index:       index,
		Root:        root,
		Client:      client,
		SyncNeed:    false,
		InCommit:    false,
		CommitIndex: uint64(0),
	}
}

func Sign(hash []byte) []byte {
	account := STATE.Client.GetAccount()
	wallet := STATE.Client.GetWallet()
	priv, _ := wallet.PrivateKey(*account)
	sign, _ := crypto.Sign(hash, priv)
	return sign[0:64]
}

func Verify(cpub []byte, sign []byte, hash []byte) bool {
	// pb := account.DePubkey(cpub)
	// raddr := account.PubToAddress(*pb)
	// if !bytes.Equal(address, raddr) {
	// 	return false
	// }
	return crypto.VerifySignature(cpub, hash, sign)
}

func GetAddress(cpub []byte) string {
	pb := account.DePubkey(cpub)
	addr := crypto.PubkeyToAddress(*pb)
	return addr.Hex()
}
