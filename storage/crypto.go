package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/pbkdf2"
)

const iterTime int = 3192 // Guaranteed random by dice-roll
const keyLength int = 32  // Generate 32-byte key

func GenerateSalt(num_bytes int) (salt string, err error) {
	var (
		salt_raw []byte
	)

	salt_raw = make([]byte, num_bytes)

	_, err = rand.Read(salt_raw)
	if err != nil {
		err = errors.New("GenerateSalt: Failed to generate salt: " + err.Error())
	}

	salt = base64.URLEncoding.EncodeToString(salt_raw)

	return
}

func GenerateHash(password_s, salt_s string) (auth_hash []byte, err error) {
	var (
		password []byte
		salt     []byte
	)

	password = []byte(password_s)
	salt, err = base64.URLEncoding.DecodeString(salt_s)
	if err != nil {
		err = errors.New("GenerateAuthHash: Failed to base64decode salt: " + err.Error())
		return
	}

	auth_hash = pbkdf2.Key(password, salt, iterTime, keyLength, sha256.New)

	return auth_hash, nil
}

func GenerateToken() string {
	return uuid.NewV4().String()
}
