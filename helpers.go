package pnghide

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"hash/crc32"
	"io"
)

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

func encryptCBC(key, plaintext []byte) (ciphertext []byte, err error) {
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return
}

// Appends padding.
func padding(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}

	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(data, pad...), nil
}

func XORBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}

	buf := make([]byte, len(a))

	for i, _ := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}

func decryptAes128Ecb(data, key []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func decrypt(s []byte, keystring []byte) []byte {
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

// EncryptDecrypt runs a XOR encryption on the input string, encrypting it if it hasn't already been,
// and decrypting it if it has, using the key provided.
func EncryptDecrypt(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}

	return output
}

// Big-endian.
func writeUint32(b []uint8, u uint32) {
	b[0] = uint8(u >> 24)
	b[1] = uint8(u >> 16)
	b[2] = uint8(u >> 8)
	b[3] = uint8(u >> 0)
}

func createChunk(w io.Writer, b []byte, name string) error {
	var header [8]byte
	var footer [4]byte

	// Calculate the length
	n := uint32(len(b))
	if int(n) != len(b) {
		return fmt.Errorf("Error with the length")
	}

	// Write the length
	writeUint32(header[:4], n)

	// Wirte the type
	header[4] = name[0]
	header[5] = name[1]
	header[6] = name[2]
	header[7] = name[3]

	// Calculate the CRC32
	crc := crc32.NewIEEE()
	crc.Write(header[4:8])
	crc.Write(b)

	// Write the CRC32
	writeUint32(footer[:4], crc.Sum32())

	// Write the chunk
	_, err := w.Write(header[:8])
	if err != nil {
		return fmt.Errorf("Error while writing the header : %s", err)
	}
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("Error while writing the header : %s", err)
	}
	_, err = w.Write(footer[:4])
	if err != nil {
		return fmt.Errorf("Error while writing the header : %s", err)
	}

	return nil
}
