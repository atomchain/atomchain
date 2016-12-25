package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	DefaultConfigFilename = "./config.json"
)

var Version string

type Configuration struct {
	Magic           int64              `json:"Magic"`
	Version         int                `json:"Version"`
	SeedList        []string           `json:"SeedList"`
	HttpRestPort    int                `json:"HttpRestPort"`
	RestCertPath    string             `json:"RestCertPath"`
	RestKeyPath     string             `json:"RestKeyPath"`
	HttpInfoPort    uint16             `json:"HttpInfoPort"`
	HttpInfoStart   bool               `json:"HttpInfoStart"`
	HttpWsPort      int                `json:"HttpWsPort"`
	HttpJsonPort    int                `json:"HttpJsonPort"`
	OauthServerUrl  string             `json:"OauthServerUrl"`
	NoticeServerUrl string             `json:"NoticeServerUrl"`
	NodePort        int                `json:"NodePort"`
	NodeType        string             `json:"NodeType"`
	WebSocketPort   int                `json:"WebSocketPort"`
	PrintLevel      int                `json:"PrintLevel"`
	IsTLS           bool               `json:"IsTLS"`
	CertPath        string             `json:"CertPath"`
	KeyPath         string             `json:"KeyPath"`
	CAPath          string             `json:"CAPath"`
	GenBlockTime    uint               `json:"GenBlockTime"`
	MultiCoreNum    uint               `json:"MultiCoreNum"`
	EncryptAlg      string             `json:"EncryptAlg"`
	MaxLogSize      int64              `json:"MaxLogSize"`
	MaxTxInBlock    int                `json:"MaxTransactionInBlock"`
	MaxHdrSyncReqs  int                `json:"MaxConcurrentSyncHeaderReqs"`
	TransactionFee  map[string]float64 `json:"TransactionFee"`
	//
	WalletFile       string   `json:"WalletFile"`
	ServerAddr       string   `json:"ServerAddr"`
	SeedString       string   `json:"TestSeed"`
	TestnetBootnodes []string `json:"TestnetBootnodes"`
	Period           int64    `json:"Period"`
	AquaData         string   `json:"AquaData"`
}

type ConfigFile struct {
	ConfigFile Configuration `json:"Configuration"`
}

var Parameters *Configuration

func init() {
	file, e := ioutil.ReadFile(DefaultConfigFilename)
	if e != nil {
		log.Fatalf("File error: %v\n", e)
		os.Exit(1)
	}
	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	config := ConfigFile{}
	e = json.Unmarshal(file, &config)
	if e != nil {
		log.Fatalf("Unmarshal json file erro %v", e)
		os.Exit(1)
	}
	fmt.Println("Parameters init finished")
	Parameters = &(config.ConfigFile)
}

// Init2 global parameters
func Init2(configFilename string) {
	// fmt.Printf("reload config file: %s\n", configFilename)
	file, e := ioutil.ReadFile(configFilename)
	if e != nil {
		log.Fatalf("File error: %v\n", e)
		os.Exit(1)
	}
	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	config := ConfigFile{}
	e = json.Unmarshal(file, &config)
	if e != nil {
		log.Fatalf("Unmarshal json file erro %v", e)
		os.Exit(1)
	}
	Parameters = &(config.ConfigFile)
}
