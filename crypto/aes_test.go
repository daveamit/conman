package crypto_test

import (
	"conman/crypto"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	pswd := []byte("password")
	cipher := crypto.NewAes(pswd, pswd)

	plainText := "some-data"

	cipherText, err := cipher.Encrypt([]byte(plainText))
	if err != nil {
		t.Error(err)
	}

	text, err := cipher.Decrypt(cipherText)
	if err != nil {
		t.Error(err)
	}

	if string(text) != plainText {
		t.Error("text and plain text don't match")
	}
}
