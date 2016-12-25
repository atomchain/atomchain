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

type Transfer struct {
	Txid   string
	From   string
	To     string
	Amount int64
	Sign   string
	Public string
}

func Serialize(trans Transfer) []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(trans); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func Unserialize(data []byte) *Transfer {
	var b bytes.Buffer
	b.Write(data)
	dec := gob.NewDecoder(&b)
	var t Transfer
	dec.Decode(&t)
	return &t
}

func GetHash(t *Transfer) []byte {
	s := t.From + t.To + strconv.FormatInt(t.Amount, 10)
	return mcrypto.Sha256([]byte(s))
}

func Sign(t *Transfer, account *accounts.Account, wallet *hdwallet.Wallet) []byte {
	s := t.From + t.To + strconv.FormatInt(t.Amount, 10)
	hash := mcrypto.Sha256([]byte(s))
	priv, _ := wallet.PrivateKey(*account)
	sign, _ := crypto.Sign(hash, priv)
	return sign[0:64]
}

func Verify(cpub []byte, sign []byte, hash []byte) bool {
	return crypto.VerifySignature(cpub, hash, sign)
}
