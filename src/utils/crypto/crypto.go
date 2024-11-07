package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
	"golang.org/x/crypto/pbkdf2"
)

// Decrypt decrypts ciphertext using the passphrase
func Decrypt(passphrase, ciphertext string) string {
	return string(DecryptToBytes(passphrase, ciphertext))
}

// Decrypt decrypts ciphertext using the passphrase. The output is a byte array
func DecryptToBytes(passphrase, ciphertext string) []byte {

	if len(ciphertext) < 50 {
		return []byte{}
	}

	salt, err := hex.DecodeString(ciphertext[0:24])
	if err != nil {
		log.Errorf("Error decoding salt: %+v", err)
	}
	iv, err := hex.DecodeString(ciphertext[24:48])
	if err != nil {
		log.Errorf("Error decoding iv: %+v", err)
	}
	data, err := hex.DecodeString(ciphertext[48:])
	if err != nil {
		log.Errorf("Error decoding ciphertext: %+v", err)
	}
	key, _ := DeriveKey([]byte(passphrase), salt)
	b, err := aes.NewCipher(key)
	if err != nil {
		log.Errorf("Error creating New Cipher: %+v", err)
	}
	aesgcm, err := cipher.NewGCM(b)
	if err != nil {
		log.Errorf("Error creating new gcm: %+v", err)
	}
	data, err = aesgcm.Open(nil, iv, data, nil)
	if err != nil {
		log.Errorf("Error executing final decoding: %+v", err)
	}
	return data
}

func DeriveKey(passphrase []byte, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 12)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}
