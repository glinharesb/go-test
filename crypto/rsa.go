package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"go-test/config"
	"io/ioutil"
	"log"
	"math/big"
)

var pemData []byte
var privateKey *rsa.PrivateKey
var block *pem.Block

func LoadRsa() {
	loadPem()
	extractDataBlock()
	decodeRsa()
}

func loadPem() {
	var err error

	pemData, err = ioutil.ReadFile(config.GetConfig().PemFile)
	if err != nil {
		log.Fatalf("read key file: %s", err)
	}
}

func extractDataBlock() {
	block, _ = pem.Decode(pemData)
	if block == nil {
		log.Fatalf("bad key data: %s", "not PEM-encoded")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Fatalf("unknown key type %q, want %q", got, want)
	}
}

func decodeRsa() {
	var err error

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("bad private key: %s", err)
	}
}

func Decrypt(encryptedBytes []byte) []byte {
	if length := len(encryptedBytes); length != 128 {
		log.Fatalf("invalid buffer length: %d", length)
	}

	c := new(big.Int).SetBytes(encryptedBytes)
	plainText := c.Exp(c, privateKey.D, privateKey.N).Bytes()

	return plainText
}
