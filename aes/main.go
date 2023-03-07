// Copyright 2023 DATATOM Authors. All rights reserved.
//AES, encrypt and decrypt

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {
	text := []byte("加密内容")
	ciphertext, _ := Base64AesCBCEncrypt(text)
	fmt.Println(string(ciphertext))
	plaintext, _ := Base64AesCBCDecrypt(ciphertext)
	fmt.Println(string(plaintext))
}

// AES加密中需要指定一个16、24或32字节的密钥，用于选择AES-128、AES-192或AES-256加密算法。
const AecKey = "512aa7be16d35e44"

// 需要生成一个随机的初始化向量（IV）作为初始状态，以便在加密过程中增加不可预测性和安全性。
const AecIV = "cdccB3uiWDu7mcxw"

var key = []byte(AecKey)
var iv = []byte(AecIV)

func Base64AesCBCEncrypt(text []byte) (string, error) {
	ciphertext, err := AesCBCEncrypt(text)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AesCBC加密
func AesCBCEncrypt(data []byte) ([]byte, error) {
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 填充内容，如果不足16位字符
	oriData := pad(data, block.BlockSize())
	// 加密，输出到[]byte数组
	cipherData := make([]byte, len(oriData))
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, oriData)
	return cipherData, nil
}

func pad(data []byte, blockSize int) []byte {
	pad := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(data, padtext...)
}

func Base64AesCBCDecrypt(text string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}
	plaintext, err := AesCBCDecrypt(data)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// AesCBC解密
func AesCBCDecrypt(data []byte) ([]byte, error) {
	// 生成密码数据块cipher.Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 输出到[]byte数组
	plainData := make([]byte, len(data))
	// 解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plainData, data)
	// 去除填充,并返回
	return unpad(plainData), nil
}

func unpad(data []byte) []byte {
	length := len(data)
	//去掉最后一次的padding
	paddLen := int(data[length-1])
	return data[:(length - paddLen)]
}
