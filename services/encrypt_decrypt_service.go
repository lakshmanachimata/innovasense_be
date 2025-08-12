package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// EncryptDecryptService provides encryption and decryption functionality
type EncryptDecryptService struct {
	key string
	iv  string
}

// NewEncryptDecryptService creates a new instance of the service
func NewEncryptDecryptService() *EncryptDecryptService {
	return &EncryptDecryptService{
		key: "innovosens2022ma",
		iv:  "smashapp01012022",
	}
}

// NewEncryptDecryptServiceWithKeys creates a new instance with custom key and IV
func NewEncryptDecryptServiceWithKeys(key, iv string) *EncryptDecryptService {
	return &EncryptDecryptService{
		key: key,
		iv:  iv,
	}
}

// set encrypts the value using AES encryption
func (e *EncryptDecryptService) set(keys string, value string) (string, error) {
	// Parse key and IV
	key := []byte(keys)
	iv := []byte(e.iv)

	// Ensure key and IV are the correct length (AES-128 requires 16 bytes)
	if len(key) < 16 {
		// Pad key to 16 bytes if shorter
		paddedKey := make([]byte, 16)
		copy(paddedKey, key)
		key = paddedKey
	} else if len(key) > 16 {
		// Truncate key to 16 bytes if longer
		key = key[:16]
	}

	if len(iv) < 16 {
		// Pad IV to 16 bytes if shorter
		paddedIV := make([]byte, 16)
		copy(paddedIV, iv)
		iv = paddedIV
	} else if len(iv) > 16 {
		// Truncate IV to 16 bytes if longer
		iv = iv[:16]
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Pad the plaintext to be a multiple of block size
	paddedValue := e.pkcs7Pad([]byte(value), aes.BlockSize)

	// Create CBC encrypter
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt the data
	ciphertext := make([]byte, len(paddedValue))
	mode.CryptBlocks(ciphertext, paddedValue)

	// Return base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// get decrypts the value using AES decryption
func (e *EncryptDecryptService) get(keys string, value string) (string, error) {
	// Parse key and IV
	key := []byte(keys)
	iv := []byte(e.iv)

	// Ensure key and IV are the correct length
	if len(key) < 16 {
		// Pad key to 16 bytes if shorter
		paddedKey := make([]byte, 16)
		copy(paddedKey, key)
		key = paddedKey
	} else if len(key) > 16 {
		// Truncate key to 16 bytes if longer
		key = key[:16]
	}

	if len(iv) < 16 {
		// Pad IV to 16 bytes if shorter
		paddedIV := make([]byte, 16)
		copy(paddedIV, iv)
		iv = paddedIV
	} else if len(iv) > 16 {
		// Truncate IV to 16 bytes if longer
		iv = iv[:16]
	}

	// Decode base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %v", err)
	}

	// Create CBC decrypter
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the data
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove padding
	unpaddedPlaintext, err := e.pkcs7Unpad(plaintext)
	if err != nil {
		return "", fmt.Errorf("failed to remove padding: %v", err)
	}

	return string(unpaddedPlaintext), nil
}

// GetEncryptData encrypts data using the default key
func (e *EncryptDecryptService) GetEncryptData(data interface{}) (string, error) {
	// Convert data to string
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	case []byte:
		dataStr = string(v)
	default:
		dataStr = fmt.Sprintf("%v", v)
	}

	return e.set(e.key, dataStr)
}

// GetDecryptData decrypts data using the default key
func (e *EncryptDecryptService) GetDecryptData(data interface{}) (string, error) {
	// Convert data to string
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	case []byte:
		dataStr = string(v)
	default:
		dataStr = fmt.Sprintf("%v", v)
	}

	return e.get(e.key, dataStr)
}

// EncryptWithKey encrypts data using a custom key
func (e *EncryptDecryptService) EncryptWithKey(key string, data interface{}) (string, error) {
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	case []byte:
		dataStr = string(v)
	default:
		dataStr = fmt.Sprintf("%v", v)
	}

	return e.set(key, dataStr)
}

// DecryptWithKey decrypts data using a custom key
func (e *EncryptDecryptService) DecryptWithKey(key string, data interface{}) (string, error) {
	var dataStr string
	switch v := data.(type) {
	case string:
		dataStr = v
	case []byte:
		dataStr = string(v)
	default:
		dataStr = fmt.Sprintf("%v", v)
	}

	return e.get(key, dataStr)
}

// pkcs7Pad adds PKCS7 padding to the data
func (e *EncryptDecryptService) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad removes PKCS7 padding from the data
func (e *EncryptDecryptService) pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding")
	}

	// Verify padding
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:length-padding], nil
}
