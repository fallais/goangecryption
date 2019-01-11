package pnghide

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

// padding
func padding(data []byte, blocklen int) ([]byte, error) {
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte("\x00"), padlen)

	return append(data, pad...), nil
}

// xorBytes
func xorBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}

	buf := make([]byte, len(a))

	for i := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}

// encryptCBC
func encryptCBC(key, iv, plaintext []byte) (ciphertext []byte, err error) {
	if len(plaintext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext = make([]byte, len(plaintext))
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, plaintext)

	return
}

// decryptECB
func decryptECB(data, key []byte) []byte {
	cipher, err := aes.NewCipher(key)
	if err == nil {
		cipher.Decrypt(data, data)
		return data
	}
	return nil
}

// decryptCBC
func decryptCBC(s []byte, keystring []byte) []byte {
	// Byte array of the string
	ciphertext := s

	// Key
	key := keystring

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext
}
