package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

// AesProvider wrapper over crypto/aes
type AesProvider struct {
	key []byte
}

// NewAes returns new AES based CryptoProvider
func NewAes(secret []byte, salt []byte) *AesProvider {
	argon2key := argon2.IDKey(secret, salt, 1, 64*1024, 4, 32)
	return &AesProvider{
		argon2key,
	}
}

// Encrypt takes plain data and returns encrypted data or failure
func (a *AesProvider) Encrypt(plainText []byte) ([]byte, error) {
	c, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	// What is GCM you say? --> https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// Create a nonce of suitable length
	nonce := make([]byte, gcm.NonceSize())
	// Randomize it
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and seal
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return cipherText, nil
}

// Decrypt take cipherText and returns plainText or failure
func (a *AesProvider) Decrypt(cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
