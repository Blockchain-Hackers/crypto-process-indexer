package database

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	// "encoding/json"
	// "flag"
	// "fmt"
	// "os"
)

// OperationResult represents the result of an encryption or decryption operation
type OperationResult struct {
	Result string `json:"result"`
	Status bool   `json:"status"`
	Error  string `json:"error,omitempty"`
}

// GetAESDecrypted decrypts given text in AES 256 CBC
func GetAESDecrypted(encrypted string, key string, iv string) (*OperationResult, *OperationResult) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, &OperationResult{Status: false, Error: err.Error()}
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, &OperationResult{Status: false, Error: err.Error()}
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, &OperationResult{Status: false, Error: "block size must be a multiple of 16"}
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return &OperationResult{Result: string(ciphertext), Status: true}, nil
}

// PKCS5UnPadding pads a certain blob of data with necessary data to be used in AES block cipher
// func PKCS5UnPadding(src []byte) []byte {
// 	length := len(src)
// 	unpadding := int(src[length-1])
// 	return src[:(length - unpadding)]
// }

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)

	if length == 0 {
		return src
	}

	unpadding := int(src[length-1])

	if unpadding > length {
		// Invalid unpadding value, return original slice
		return src
	}

	return src[:(length - unpadding)]
}


// GetAESEncrypted encrypts given text in AES 256 CBC
func GetAESEncrypted(plaintext string, key string, iv string) (*OperationResult, *OperationResult) {
	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, &OperationResult{Status: false, Error: err.Error()}
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	result := base64.StdEncoding.EncodeToString(ciphertext)

	return &OperationResult{Result: result, Status: true}, nil
}

