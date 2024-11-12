package utils

import "crypto/rand"

const (
	wide string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-$#!_"
)

// GenerateRandomString returns a securely generated random string of length n.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue
func GenerateRandomString(n int) (*string, error) {
	return generateRandomStringCode(n, wide)
}

func generateRandomStringCode(n int, letters string) (*string, error) {
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return nil, err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	result := string(bytes)
	return &result, nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
