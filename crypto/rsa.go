package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"go-test/config"
	"io/ioutil"
	"log"
	"math/big"
)

var pemData []byte
var privateKey *rsa.PrivateKey
var block *pem.Block

func LoadRsa() error {
	if err := loadPem(); err != nil {
		return err
	}

	if err := extractDataBlock(); err != nil {
		return err
	}

	if err := decodeRsa(); err != nil {
		return err
	}

	return nil
}

func loadPem() error {
	var err error

	pemData, err = ioutil.ReadFile(config.GetConfig().PemFile)
	if err != nil {
		return fmt.Errorf("read key file: %s", err)
	}

	return nil
}

func extractDataBlock() error {
	block, _ = pem.Decode(pemData)
	if block == nil {
		return fmt.Errorf("bad key data: %s", "not PEM-encoded")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return fmt.Errorf("unknown key type %q, want %q", got, want)
	}

	return nil
}

func decodeRsa() error {
	var err error

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("bad private key: %s", err)
	}

	return nil
}

func Decrypt(encryptedBytes []byte) []byte {
	if length := len(encryptedBytes); length != 128 {
		log.Fatalf("invalid buffer length: %d", length)
	}

	c := new(big.Int).SetBytes(encryptedBytes)
	plainText := c.Exp(c, privateKey.D, privateKey.N).Bytes()

	return plainText
}
