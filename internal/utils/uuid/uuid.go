package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// GenerateUUID UUIDを生成する
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(uuid), nil
}

// HashPassword パスワードをハッシュ化する
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword パスワードとハッシュを比較検証する
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
