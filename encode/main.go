package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

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
func EncryptAes128Ecb(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func main() {
	origData := "{\\\"userName\\\":\\\"datatom\\\",\\\"nickName\\\":\\\"德拓\\\", \\\"phone\\\":\\\"111\\\"}"
	val, err := EncryptAes128Ecb([]byte(origData), []byte("MTY4NjgxOTI2Mw=="))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(val))
}
