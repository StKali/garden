package util

import (
	"fmt"
	"github.com/stkali/errors"
	"golang.org/x/crypto/bcrypt"
	"os"
)

const START_FAILED = 1

// CheckError prints the msg with the prefix and exits with error code 1.
// If the msg is nil, it does nothing.
func CheckError(text string, err error) {
	if err == nil {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "error: %s, err: %s\n", text, err)
	os.Exit(START_FAILED)
}

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password, err: %s", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword checks the password is correct or not
func VerifyPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
