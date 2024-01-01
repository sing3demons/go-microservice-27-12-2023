package utils

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func RunCMD(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
			panic("some error found")
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	return out
}
