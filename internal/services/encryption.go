package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(password string) (string, error) {
	key := os.Getenv("KEY")
	block, err := aes.NewCipher([]byte(createHash(key)))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(password), nil)
	return string(encrypted), nil
}

func decrypt(crypted string) (string, error) {
	key := os.Getenv("KEY")
	block, err := aes.NewCipher([]byte(createHash(key)))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(crypted) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := crypted[:nonceSize], crypted[nonceSize:]
	decrypted, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", nil
	}

	return string(decrypted), nil
}
