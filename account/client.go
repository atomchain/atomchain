package account

import (
	"bytes"
	. "chain/common"
	"chain/common/log"
	"chain/core/contract"

	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	syslog "log"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

type ClientImpl struct {
	mu          sync.Mutex
	iv          []byte
	masterKey   []byte
	path        string
	mainAccount *accounts.Account
	accounts    map[Uint160]*Account
	contracts   map[Uint160]*contract.Contract
	wallet      *hdwallet.Wallet

	FileStore
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func getAESKey(pass []byte) []byte {
	passhash := sha256.Sum256(pass)
	passhash2 := sha256.Sum256(passhash[:])
	return passhash2[:]
}

func AesEncrypt(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

func AesDecrypt(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("invalid decrypt key")
	}

	blockModel := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}

func CreateClient(path string, password []byte) *ClientImpl {
	client := &ClientImpl{
		path:      path,
		FileStore: FileStore{path: path},
		accounts:  map[Uint160]*Account{},
		contracts: map[Uint160]*contract.Contract{},
	}
	entropy, _ := bip39.NewEntropy(128)
	memos, _ := bip39.NewMnemonic(entropy)
	fmt.Print("please write the words down:\n")
	color.Red(memos)
	wallet, err := hdwallet.NewFromMnemonic(memos)
	if err != nil {
		syslog.Fatal(err)
	}
	walletPath := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(walletPath, false)
	if err != nil {
		log.Fatal(err)
	}
	address, _ := wallet.AddressHex(account)
	color.Green("Your Master Key is stored")
	color.Green("Your Address is %s", address)

	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// client.iv = make([]byte, 16)
	// client.masterKey = make([]byte, 32)
	// for i := 0; i < 16; i++ {
	// 	client.iv[i] = byte(r.Intn(256))
	// }
	// for i := 0; i < 32; i++ {
	// 	client.masterKey[i] = byte(r.Intn(256))
	// }
	client.BuildDatabase(path)
	if err := client.SaveStoredData("Version", []byte(WalletStoreVersion)); err != nil {
		log.Error(err)
		return nil
	}
	seed, _ := hdwallet.NewSeedFromMnemonic(memos)
	if err := client.SaveStoredData("Seed", seed); err != nil {
		log.Error(err)
		return nil
	}
	client.wallet = wallet
	client.mainAccount = &account
	return client
}

func OpenClient(path string) *ClientImpl {
	client := &ClientImpl{
		path:      path,
		FileStore: FileStore{path: path},
		accounts:  map[Uint160]*Account{},
		contracts: map[Uint160]*contract.Contract{},
	}
	seed, err := client.LoadStoredData("Seed")
	if err != nil {
		log.Warn(err)
		return nil
	}
	wallet, _ := hdwallet.NewFromSeed(seed)
	walletPath := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, _ := wallet.Derive(walletPath, false)
	// priv, _ := wallet.PrivateKey(account)
	// pub, _ := wallet.PublicKey(account)
	// cpub := crypto.CompressPubkey(pub)
	// data, err := crypto.Sign([]byte("asadasadasadasadasadasadasadasad"), priv)
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// fmt.Print(len(data[0:64]))
	// is := crypto.VerifySignature(cpub, []byte("asadasadasadasadasadasadasadasad"), data[0:64])
	// fmt.Print(is)
	client.wallet = wallet
	client.mainAccount = &account
	return client
}

func (cl *ClientImpl) CreateAccount() (*Account, error) {
	account, err := NewAccount()
	if err != nil {
		return nil, err
	}
	if err := cl.SaveAccount(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (cl *ClientImpl) EncryptPrivateKey(prikey []byte) ([]byte, error) {
	enc, err := AesEncrypt(prikey, cl.masterKey, cl.iv)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

func (cl *ClientImpl) SaveAccount(ac *Account) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	// save Account to memory
	programHash := ac.ProgramHash
	cl.accounts[programHash] = ac

	decryptedPrivateKey := make([]byte, 96)
	temp, err := ac.PublicKey.EncodePoint(false)
	if err != nil {
		return err
	}
	for i := 1; i <= 64; i++ {
		decryptedPrivateKey[i-1] = temp[i]
	}
	for i := len(ac.PrivateKey) - 1; i >= 0; i-- {
		decryptedPrivateKey[96+i-len(ac.PrivateKey)] = ac.PrivateKey[i]
	}
	encryptedPrivateKey, err := cl.EncryptPrivateKey(decryptedPrivateKey)
	if err != nil {
		return err
	}
	ClearBytes(decryptedPrivateKey, 96)

	// save Account keys to db
	err = cl.SaveAccountData(programHash.ToArray(), encryptedPrivateKey)
	if err != nil {
		return err
	}

	return nil
}

func (cl *ClientImpl) CreateContract(account *Account) error {
	contract, err := contract.CreateSignatureContract(account.PubKey())
	if err != nil {
		return err
	}
	if err := cl.SaveContract(contract); err != nil {
		return err
	}
	return nil
}

func (cl *ClientImpl) SaveContract(ct *contract.Contract) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	// save contract to memory
	cl.contracts[ct.ProgramHash] = ct

	// save contract to db
	return cl.SaveContractData(ct)
}

func (cl *ClientImpl) GetPath() string {
	return cl.path
}

func (cl *ClientImpl) GetWallet() *hdwallet.Wallet {
	return cl.wallet
}

func (cl *ClientImpl) GetAccount() *accounts.Account {
	return cl.mainAccount
}

func (cl *ClientImpl) GetAddress() string {
	baddr, _ := cl.wallet.AddressBytes(*cl.mainAccount)
	addr := BytesToHexString(baddr)
	return addr
}

func (cl *ClientImpl) Getcpub() []byte {
	pub, _ := cl.wallet.PublicKey(*cl.mainAccount)
	return crypto.CompressPubkey(pub)
}

func Create(path string, passwordKey []byte) (*ClientImpl, error) {
	client := CreateClient(path, passwordKey)
	if client == nil {
		return nil, errors.New("client nil")
	}
	// account, err := client.CreateAccount()
	// if err != nil {
	// 	return nil, err
	// }
	// if err := client.CreateContract(account); err != nil {
	// 	return nil, err
	// }

	// client.mainAccount = account.ProgramHash
	return client, nil
}

func Open(path string) (*ClientImpl, error) {
	client := OpenClient(path)
	if client == nil {
		return nil, errors.New("You have to create an account")
	}
	return client, nil
}

func PubToAddress(p ecdsa.PublicKey) []byte {
	address := crypto.PubkeyToAddress(p)
	return address.Bytes()
}

func DePubkey(p []byte) *ecdsa.PublicKey {
	pk, _ := crypto.DecompressPubkey(p)
	return pk
}
