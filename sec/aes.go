package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const (
	CipherKey = "L9EhKyurl9vI28dRjjOECMh3dv09wkLe"
)

// func Encrypt(key []byte, message string) (encmess string, err error) {
func Encrypt(message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher([]byte(CipherKey))
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return "", nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return encmess, nil
}

func Decrypt(securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	block, err := aes.NewCipher([]byte(CipherKey))
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return "", nil
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return decodedmess, nil
}
