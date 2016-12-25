package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	. "chain/common"
	"chain/core/contract"
	. "chain/errors"
)

const (
	WalletStoreVersion = "0.0.1"
)

type WalletData struct {
	PasswordHash string
	IV           string
	MasterKey    string
	Height       uint32
	Version      string
	Seed         string
}

type AccountData struct {
	Address             string
	ProgramHash         string
	PrivateKeyEncrypted string
	Type                string
}

type ContractData struct {
	ProgramHash string
	RawData     string
}

type CoinData string

type FileData struct {
	WalletData
	Account  []AccountData
	Contract []ContractData
	Coins    CoinData
}

type FileStore struct {
	// this lock could be hold by readDB, writeDB and interrupt signals.
	sync.Mutex

	data FileData
	file *os.File
	path string
}

// Caller holds the lock and reads bytes from DB, then close the DB and release the lock
func (cs *FileStore) readDB() ([]byte, error) {
	cs.Lock()
	defer cs.Unlock()
	defer cs.closeDB()

	var err error
	cs.file, err = os.OpenFile(cs.path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if cs.file != nil {
		data, err := ioutil.ReadAll(cs.file)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, NewDetailErr(errors.New("[readDB] file handle is nil"), ErrNoCode, "")
	}
}

// Caller holds the lock and writes bytes to DB, then close the DB and release the lock
func (cs *FileStore) writeDB(data []byte) error {
	cs.Lock()
	defer cs.Unlock()
	defer cs.closeDB()

	var err error
	cs.file, err = os.OpenFile(cs.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	if cs.file != nil {
		cs.file.Write(data)
	}

	return nil
}

func (cs *FileStore) closeDB() {
	if cs.file != nil {
		cs.file.Close()
		cs.file = nil
	}
}

func (cs *FileStore) BuildDatabase(path string) {
	os.Remove(path)
	jsonBlob, err := json.Marshal(cs.data)
	if err != nil {
		fmt.Println("Build DataBase Error")
		os.Exit(1)
	}
	cs.writeDB(jsonBlob)
}

func (cs *FileStore) SaveAccountData(programHash []byte, encryptedPrivateKey []byte) error {
	JSONData, err := cs.readDB()
	if err != nil {
		return errors.New("error: reading db")
	}
	if err := json.Unmarshal(JSONData, &cs.data); err != nil {
		return errors.New("error: unmarshal db")
	}

	var accountType string
	if len(cs.data.Account) == 0 {
		accountType = MAINACCOUNT
	} else {
		accountType = SUBACCOUNT
	}

	pHash, err := Uint160ParseFromBytes(programHash)
	if err != nil {
		return errors.New("invalid program hash")
	}
	addr, err := pHash.ToAddress()
	if err != nil {
		return errors.New("invalid address")
	}
	a := AccountData{
		Address:             addr,
		ProgramHash:         BytesToHexString(programHash),
		PrivateKeyEncrypted: BytesToHexString(encryptedPrivateKey),
		Type:                accountType,
	}
	cs.data.Account = append(cs.data.Account, a)

	JSONBlob, err := json.Marshal(cs.data)
	if err != nil {
		return errors.New("error: marshal db")
	}
	cs.writeDB(JSONBlob)

	return nil
}

func (cs *FileStore) DeleteAccountData(programHash string) error {
	JSONData, err := cs.readDB()
	if err != nil {
		return errors.New("error: reading db")
	}
	if err := json.Unmarshal(JSONData, &cs.data); err != nil {
		return errors.New("error: unmarshal db")
	}

	for i, v := range cs.data.Account {
		if programHash == v.ProgramHash {
			if v.Type == MAINACCOUNT {
				return errors.New("Can't remove main account")
			}
			cs.data.Account = append(cs.data.Account[:i], cs.data.Account[i+1:]...)
		}
	}

	JSONBlob, err := json.Marshal(cs.data)
	if err != nil {
		return errors.New("error: marshal db")
	}
	cs.writeDB(JSONBlob)

	return nil
}

func (cs *FileStore) LoadAccountData() ([]AccountData, error) {
	JSONData, err := cs.readDB()
	if err != nil {
		return nil, errors.New("error: reading db")
	}
	if err := json.Unmarshal(JSONData, &cs.data); err != nil {
		return nil, errors.New("error: unmarshal db")
	}
	return cs.data.Account, nil
}

func (cs *FileStore) SaveStoredData(name string, value []byte) error {
	JSONData, err := cs.readDB()
	if err != nil {
		return errors.New("error: reading db")
	}
	if err := json.Unmarshal(JSONData, &cs.data); err != nil {
		return errors.New("error: unmarshal db")
	}

	hexValue := BytesToHexString(value)
	switch name {
	case "Version":
		cs.data.Version = string(value)
	case "Seed":
		cs.data.Seed = hexValue
	}
	JSONBlob, err := json.Marshal(cs.data)
	if err != nil {
		return errors.New("error: marshal db")
	}
	cs.writeDB(JSONBlob)

	return nil
}

func (cs *FileStore) LoadStoredData(name string) ([]byte, error) {
	JSONData, err := cs.readDB()
	if err != nil {
		return nil, errors.New("error: reading db")
	}
	if err := json.Unmarshal(JSONData, &cs.data); err != nil {
		return nil, errors.New("error: unmarshal db")
	}
	switch name {
	case "Version":
		return []byte(cs.data.Version), nil
	case "Seed":
		return HexStringToBytes(cs.data.Seed)
	}

	return nil, errors.New("Can't find the key: " + name)
}

func (cs *FileStore) SaveContractData(ct *contract.Contract) error {
	JSONData, err := cs.readDB()
	if err != nil {
		return errors.New("error: reading db")
	}
	if err := json.Unmarshal(JSONData, &cs.data); err != nil {
		return errors.New("error: unmarshal db")
	}
	c := ContractData{
		ProgramHash: BytesToHexString(ct.ProgramHash.ToArray()),
		RawData:     BytesToHexString(ct.ToArray()),
	}
	cs.data.Contract = append(cs.data.Contract, c)

	JSONBlob, err := json.Marshal(cs.data)
	if err != nil {
		return errors.New("error: marshal db")
	}
	cs.writeDB(JSONBlob)

	return nil
}
