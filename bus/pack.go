package bus

import (
	"bytes"
	"strconv"
	"strings"

	. "chain/common"
	"chain/crypto"
	"encoding/gob"
	"encoding/json"
)

func PackSyncAsk(mtr string, idx uint64) []byte {
	sa := SyncAsk{
		Index: idx,
		Root:  mtr,
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(sa); err != nil {
		panic(err)
	}
	saData := b.Bytes()
	m := Msg{
		T:    "SA",
		Data: BytesToHexString(saData),
	}
	data, _ := json.Marshal(m)
	return data
}

func PackSyncResp(mtr string, idx uint64) []byte {
	sp := SyncResp{
		Index: idx,
		Root:  mtr,
	}
	tohashstirng := strconv.FormatUint(STATE.Index, 10) + mtr
	hash := crypto.Sha256([]byte(tohashstirng))
	sign := Sign(hash)
	sp.Signature = BytesToHexString(sign)
	sp.Pubkey = BytesToHexString(STATE.Client.Getcpub())
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(sp); err != nil {
		panic(err)
	}
	saData := b.Bytes()
	m := Msg{
		T:    "SP",
		Data: BytesToHexString(saData),
	}
	data, _ := json.Marshal(m)
	return data
}

func PackViewSync(vn uint64) []byte {
	vs := ViewSync{
		ViewNumber: vn,
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(vs); err != nil {
		panic(err)
	}
	vsData := b.Bytes()
	m := Msg{
		T:    "VS",
		Data: BytesToHexString(vsData),
	}
	data, _ := json.Marshal(m)
	return data
}

func PackViewResp(vn uint64) []byte {
	vr := ViewResp{
		ViewNumber: vn,
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(vr); err != nil {
		panic(err)
	}
	vrData := b.Bytes()
	m := Msg{
		T:    "VR",
		Data: BytesToHexString(vrData),
	}
	data, _ := json.Marshal(m)
	return data
}

func PackCommitSync(idx uint64, bdata []byte) []byte {
	cs := CommitSync{
		Index:     idx,
		BlockData: bdata,
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(cs); err != nil {
		panic(err)
	}
	csData := b.Bytes()
	m := Msg{
		T:    "CS",
		Data: BytesToHexString(csData),
	}
	data, _ := json.Marshal(m)
	return data
}

func PackCommitResp(idx uint64) []byte {
	addr := "0x" + strings.ToLower(STATE.Client.GetAddress())
	cr := CommitResp{
		Index: idx,
		From:  addr,
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(cr); err != nil {
		panic(err)
	}
	crData := b.Bytes()
	m := Msg{
		T:    "CR",
		Data: BytesToHexString(crData),
	}
	data, _ := json.Marshal(m)
	return data
}
