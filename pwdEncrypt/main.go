package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
)

var rsaCrypto *RsaCrypto

type RsaCrypto struct {
	publicKey *rsa.PublicKey
}

const publicKeyPath = `-----BEGIN PUBLIC KEY-----
xxxxxxx
-----END PUBLIC KEY-----`

func init() {
	pubKey := []byte(publicKeyPath)
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	rsaCrypto = &RsaCrypto{
		publicKey: pubInterface.(*rsa.PublicKey),
	}
}

// PublicKeyEncrypt rsa公钥加密
func (rc *RsaCrypto) PublicKeyEncrypt(data []byte) ([]byte, error) {
	encode, err := rsa.EncryptPKCS1v15(rand.Reader, rc.publicKey, data)
	if err != nil {
		return nil, err
	}
	return []byte(base64.StdEncoding.EncodeToString(encode)), nil
}

// RsaEncrypt 加密
func RsaEncrypt(data []byte) ([]byte, error) {
	return rsaCrypto.PublicKeyEncrypt(data)
}

func main() {
	pwd := "Aa123456789+"
	hash := md5.Sum([]byte(pwd))
	pwd = fmt.Sprintf("%x", hash)
	newPwd, err := RsaEncrypt([]byte(pwd))
	if err != nil {
		log.Println("Error:", err)
	}
	fmt.Println("最终的加密结果", string(newPwd))
}
