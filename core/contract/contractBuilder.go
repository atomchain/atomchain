package contract

import (
	"chain/common"
	. "chain/common"
	pg "chain/core/contract/program"
	"chain/crypto"
	. "chain/errors"
	"io"
	"math/big"
	"sort"
)

//SignableData describe the data need be signed.
type SignableData interface {
	//Get the the SignableData's program hashes
	GetProgramHashes() ([]common.Uint160, error)

	SetPrograms([]*pg.Program)

	GetPrograms() []*pg.Program

	//TODO: add SerializeUnsigned
	SerializeUnsigned(io.Writer) error
}

type ContractContext struct {
	Data          SignableData
	ProgramHashes []Uint160
	Codes         [][]byte
	Parameters    [][][]byte

	MultiPubkeyPara [][]PubkeyParameter

	//temp index for multi sig
	tempParaIndex int
}

//create a Single Singature contract for owner
func CreateSignatureContract(ownerPubKey *crypto.PubKey) (*Contract, error) {
	temp, err := ownerPubKey.EncodePoint(true)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureContract failed.")
	}
	signatureRedeemScript, err := CreateSignatureRedeemScript(ownerPubKey)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureContract failed.")
	}
	hash, err := ToCodeHash(temp)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureContract failed.")
	}
	signatureRedeemScriptHashToCodeHash, err := ToCodeHash(signatureRedeemScript)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureContract failed.")
	}
	return &Contract{
		Code:            signatureRedeemScript,
		Parameters:      []ContractParameterType{Signature},
		ProgramHash:     signatureRedeemScriptHashToCodeHash,
		OwnerPubkeyHash: hash,
	}, nil
}

func CreateSignatureRedeemScript(pubkey *crypto.PubKey) ([]byte, error) {
	temp, err := pubkey.EncodePoint(true)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureRedeemScript failed.")
	}
	sb := pg.NewProgramBuilder()
	sb.PushData(temp)
	sb.AddOp(pg.CHECKSIG)
	return sb.ToArray(), nil
}

//create a Multi Singature contract for owner  。
func CreateMultiSigContract(publicKeyHash Uint160, m int, publicKeys []*crypto.PubKey) (*Contract, error) {

	params := make([]ContractParameterType, m)
	for i, _ := range params {
		params[i] = Signature
	}
	MultiSigRedeemScript, err := CreateMultiSigRedeemScript(m, publicKeys)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureRedeemScript failed.")
	}
	signatureRedeemScriptHashToCodeHash, err := ToCodeHash(MultiSigRedeemScript)
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureContract failed.")
	}
	return &Contract{
		Code:            MultiSigRedeemScript,
		Parameters:      params,
		ProgramHash:     signatureRedeemScriptHashToCodeHash,
		OwnerPubkeyHash: publicKeyHash,
	}, nil
}

func CreateMultiSigRedeemScript(m int, pubkeys []*crypto.PubKey) ([]byte, error) {
	if !(m >= 1 && m <= len(pubkeys) && len(pubkeys) <= 24) {
		return nil, nil //TODO: add panic
	}

	sb := pg.NewProgramBuilder()
	sb.PushNumber(big.NewInt(int64(m)))

	//sort pubkey
	sort.Sort(crypto.PubKeySlice(pubkeys))

	for _, pubkey := range pubkeys {
		temp, err := pubkey.EncodePoint(true)
		if err != nil {
			return nil, NewDetailErr(err, ErrNoCode, "[Contract],CreateSignatureContract failed.")
		}
		sb.PushData(temp)
	}

	sb.PushNumber(big.NewInt(int64(len(pubkeys))))
	sb.AddOp(pg.CHECKMULTISIG)
	return sb.ToArray(), nil
}
