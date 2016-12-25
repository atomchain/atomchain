package bus

import (
	"bytes"
	"chain/account"
	"chain/common"
	"chain/consensus/transaction"
	"encoding/gob"

	log "github.com/sirupsen/logrus"
)

func DealAndVerifyTx(data []byte) *transaction.Transfer {
	tx := transaction.Unserialize(data)
	log.Infof("Received Transfer: %+v", tx)
	bpub, _ := common.HexStringToBytes(tx.Public)
	from, _ := common.HexStringToBytes(tx.From)
	pub := account.DePubkey(bpub)
	addr := account.PubToAddress(*pub)
	if !bytes.Equal(from, addr) {
		return nil
	}
	sign, _ := common.HexStringToBytes(tx.Sign)
	hash := transaction.GetHash(tx)
	if !transaction.Verify(bpub, sign, hash) {
		return nil
	}
	return tx
}

func DealSyncAsk(data []byte) *SyncAsk {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var sa SyncAsk
	dec.Decode(&sa)
	return &sa
}

func DealSyncResp(data []byte) *SyncResp {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var sp SyncResp
	dec.Decode(&sp)
	return &sp
}

func DealViewSync(data []byte) *ViewSync {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var vs ViewSync
	dec.Decode(&vs)
	return &vs
}

func DealViewResp(data []byte) *ViewResp {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var vr ViewResp
	dec.Decode(&vr)
	return &vr
}

func DealCommitSync(data []byte) *CommitSync {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var cs CommitSync
	dec.Decode(&cs)
	return &cs
}

func DealCommitResp(data []byte) *CommitResp {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var cr CommitResp
	dec.Decode(&cr)
	return &cr
}
