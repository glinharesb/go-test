package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
)

type Rsa struct {
	pemFile    string
	pemData    []byte
	privateKey *rsa.PrivateKey
	block      *pem.Block
}

var RsaInstance *Rsa = &Rsa{
	pemFile: "key.pem",
}

func (r *Rsa) Load() {
	r.loadPem()
	r.extractDataBlock()
	r.decodeRsa()
}

func (r *Rsa) loadPem() {
	pemData, err := ioutil.ReadFile(r.pemFile)
	if err != nil {
		log.Fatalf("read key file: %s", err)
	}

	RsaInstance.pemData = pemData
}

func (r *Rsa) extractDataBlock() {
	block, _ := pem.Decode(r.pemData)
	if block == nil {
		log.Fatalf("bad key data: %s", "not PEM-encoded")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Fatalf("unknown key type %q, want %q", got, want)
	}

	RsaInstance.block = block
}

func (r *Rsa) decodeRsa() {
	priv, err := x509.ParsePKCS1PrivateKey(r.block.Bytes)
	if err != nil {
		log.Fatalf("bad private key: %s", err)
	}

	r.privateKey = priv
}

func (r *Rsa) Decrypt(encryptedBytes []byte) []byte {
	if length := len(encryptedBytes); length != 128 {
		log.Fatalf("invalid buffer length: %d", length)
	}

	c := new(big.Int).SetBytes(encryptedBytes)
	plainText := c.Exp(c, r.privateKey.D, r.privateKey.N).Bytes()

	return plainText
}
