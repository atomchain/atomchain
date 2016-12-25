package block

import (
	"bytes"
	"encoding/json"
	"fmt"

	"chain/crypto"
)

type BlockBody struct {
	Index        uint64
	Transactions [][]byte
}

//json encoding of body only
func (bb *BlockBody) Marshal() ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(bf)
	if err := enc.Encode(bb); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (bb *BlockBody) Unmarshal(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b) //will read from b
	if err := dec.Decode(bb); err != nil {
		return err
	}
	return nil
}

func (bb *BlockBody) Hash() ([]byte, error) {
	hashBytes, err := bb.Marshal()
	if err != nil {
		return nil, err
	}
	return crypto.Sha256(hashBytes), nil
}

//------------------------------------------------------------------------------

type BlockSignature struct {
	Validator string
	Index     uint64
	Signature []byte
}

func (bs *BlockSignature) ValidatorHex() string {
	return fmt.Sprintf("0x%X", bs.Validator)
}

func (bs *BlockSignature) Marshal() ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(bf)
	if err := enc.Encode(bs); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (bs *BlockSignature) Unmarshal(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b) //will read from b
	if err := dec.Decode(bs); err != nil {
		return err
	}
	return nil
}

func (bs *BlockSignature) ToWire() WireBlockSignature {
	return WireBlockSignature{
		Index:     bs.Index,
		Signature: bs.Signature,
	}
}

type WireBlockSignature struct {
	Index     uint64
	Signature []byte
}

//------------------------------------------------------------------------------

type Block struct {
	Body       BlockBody
	Packer     string
	Signatures map[string][]byte // [validator hex] => signature

	hash []byte
	hex  string
}

func NewBlock(blockIndex uint64, transactions [][]byte) Block {
	body := BlockBody{
		Index:        blockIndex,
		Transactions: transactions,
	}
	return Block{
		Body:       body,
		Signatures: make(map[string][]byte),
	}
}

func (b *Block) Index() uint64 {
	return b.Body.Index
}

func (b *Block) Transactions() [][]byte {
	return b.Body.Transactions
}

func (b *Block) GetSignature(validator string) (res BlockSignature, err error) {
	sig, ok := b.Signatures[validator]
	if !ok {
		return res, fmt.Errorf("signature not found")
	}
	return BlockSignature{
		Validator: "test",
		Index:     b.Index(),
		Signature: sig,
	}, nil
}

func (b *Block) AppendTransactions(txs [][]byte) {
	b.Body.Transactions = append(b.Body.Transactions, txs...)
}

func (b *Block) Marshal() ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(bf)
	if err := enc.Encode(b); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (b *Block) Unmarshal(data []byte) error {
	bf := bytes.NewBuffer(data)
	dec := json.NewDecoder(bf)
	if err := dec.Decode(b); err != nil {
		return err
	}
	return nil
}

func (b *Block) Hash() ([]byte, error) {
	if len(b.hash) == 0 {
		hashBytes, err := b.Marshal()
		if err != nil {
			return nil, err
		}
		b.hash = crypto.Sha256(hashBytes)
	}
	return b.hash, nil
}

func (b *Block) Hex() string {
	if b.hex == "" {
		hash, _ := b.Hash()
		b.hex = fmt.Sprintf("0x%X", hash)
	}
	return b.hex
}

func (b *Block) Sign(privKey []byte) (bs BlockSignature, err error) {

	signBytes, err := b.Body.Hash()
	if err != nil {
		return bs, err
	}
	R, err := crypto.Sign(privKey, signBytes)
	if err != nil {
		return bs, err
	}
	signature := BlockSignature{
		Index:     b.Index(),
		Signature: R,
	}

	return signature, nil
}

func (b *Block) SetSignature(bs BlockSignature) error {
	b.Signatures[bs.ValidatorHex()] = bs.Signature
	return nil
}

func (b *Block) Verify(sig BlockSignature) (bool, error) {

	return true, nil
}
