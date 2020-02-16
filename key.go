package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
)

// Key is a struct that hold ethereum keys address
type Key struct {
	private []byte
	public  []byte
	address string
	value   uint64
}

// Debug a key
func (key *Key) Debug() {
	println(key.getJSON())
}

func (key *Key) getJSON() string {
	result := map[string]string{
		"private": hex.EncodeToString(key.private),
		"public":  hex.EncodeToString(key.public),
		"address": key.address,
		"amount":  string(key.value),
	}
	keyjson, _ := json.Marshal(result)
	return string(keyjson)
}

// Save the key at the desired path
func (key *Key) Save(keydir string) error {
	path := filepath.Join(keydir, "0x"+key.address+".key")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	keyjson := key.getJSON()
	if _, err := writer.WriteString(string(keyjson)); err != nil {
		return err
	}
	writer.Flush()
	return nil
}

// GenerateKey generate a random ethereum key
func GenerateKey() *Key {
	key := Key{}
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	key.private = crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	key.public = crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	key.address = address.Hex()[2:]
	return &key
}
