package utils

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"os"
)

func HostName(name, id string) string {
	host := os.Getenv("CATEGORY_SERVICE_URL")
	return host + "/" + (name) + "/" + id
}

func Href(host, name, id string) string {
	return fmt.Sprintf("%s/%s/%s", host, name, id)
}

func EncryptAES(plaintext string) string {
	aesKey := os.Getenv("AES_KEY")
	c, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		panic(err)
	}

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func DecryptAES(ct string) string {
	ciphertext, _ := hex.DecodeString(ct)
	aesKey := os.Getenv("AES_KEY")
	c, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		panic(err)
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	fmt.Println("DECRYPTED:", s)
	return s
}
