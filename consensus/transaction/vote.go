package transaction

import (
	"bytes"
	"encoding/gob"
	"strconv"

	mcrypto "chain/crypto"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type Vote struct {
	From   string
	To     string
	Amount int64
	Sign   string
	Void   string
	Public string
}

func VSerialize(v Vote) []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(v); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func VUnserialize(data []byte) *Vote {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var t Vote
	dec.Decode(&t)
	return &t
}

func GetVHash(v *Vote) []byte {
	s := v.From + v.To + strconv.FormatInt(v.Amount, 10)
	return mcrypto.Sha256([]byte(s))
}

func VSign(v *Vote, account *accounts.Account, wallet *hdwallet.Wallet) []byte {
	s := v.From + v.To + strconv.FormatInt(v.Amount, 10)
	hash := mcrypto.Sha256([]byte(s))
	priv, _ := wallet.PrivateKey(*account)
	sign, _ := crypto.Sign(hash, priv)
	return sign[0:64]
}

func VVerify(cpub []byte, sign []byte, hash []byte) bool {
	return crypto.VerifySignature(cpub, hash, sign)
}
