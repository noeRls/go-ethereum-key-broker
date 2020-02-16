package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"

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
	println("Key:")
	println("  private: ", hex.EncodeToString(key.private))
	println("  public: ", hex.EncodeToString(key.public))
	println("  address: ", key.address)
	println("  amount: ", key.value)
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
