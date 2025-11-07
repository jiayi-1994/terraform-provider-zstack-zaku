package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"

	"github.com/forgoer/openssl"
)

const keyLength = 16

var defaultKeyGenerator KeyGenerator = fixedGenerator{}

// KeyGenerator will generate key for util.AESHandler
type KeyGenerator interface {
	KeyGenerate(key string) []byte
}

type fixedGenerator struct {
	// empty struct
}

func (generator fixedGenerator) KeyGenerate(key string) []byte {
	keyBytes := []byte(key)
	switch l := len(keyBytes); {
	case l < keyLength:
		keyBytes = append(keyBytes, make([]byte, keyLength-l)...)
	case l > keyLength:
		keyBytes = keyBytes[:keyLength]
	}
	return keyBytes
}

// AESHandler is a helper for handle aes encrypt/decrypt etc.
type AESHandler struct {
	KeyGenerator KeyGenerator
}

func (handler *AESHandler) encrypt(key string, origData []byte) ([]byte, error) {
	handler.checkGenerator()
	keyBytes := handler.KeyGenerator.KeyGenerate(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func (handler *AESHandler) decrypt(key string, crypted []byte) ([]byte, error) {
	handler.checkGenerator()
	keyBytes := handler.KeyGenerator.KeyGenerate(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyBytes[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5Padding pad with bytes
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS5UnPadding remove padding bytes
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	uLen := int(origData[length-1])
	return origData[:(length - uLen)]
}

// AesEncrypt encrypt a 'val' with aes then encode with base64
func (handler *AESHandler) AesEncrypt(key string, val string) (string, error) {
	origData := []byte(val)
	encrypted, err := handler.encrypt(key, origData)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// AesDecrypt decode 'val' with base64 then decrypt with aes.
func (handler *AESHandler) AesDecrypt(key string, val string) ([]byte, error) {
	encryptedData, err := base64.URLEncoding.DecodeString(val)
	if err != nil {
		return nil, err
	}

	return handler.decrypt(key, encryptedData)
}

func (handler *AESHandler) checkGenerator() {
	if handler.KeyGenerator == nil {
		handler.KeyGenerator = defaultKeyGenerator
	}
}

func EncryptByAccessKey(key, value string) string {
	sum := md5.Sum([]byte(key))
	sumStr := fmt.Sprintf("%x", sum)
	fmt.Println("sumStr: ", sumStr)
	iv := []byte(sumStr)[:16]
	fmt.Println("iv: ", iv)
	cbcDecrypt, _ := openssl.AesCBCEncrypt([]byte(value), []byte(sumStr), iv, openssl.PKCS7_PADDING)
	toString := base64.StdEncoding.EncodeToString(cbcDecrypt)
	return toString
}
