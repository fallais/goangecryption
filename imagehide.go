package pnghide

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

var (
	// ErrBinaryTooLong is raised when binary representation of character is too long.
	ErrBinaryTooLong = errors.New("The binary representation of character is too long")

	// ErrParsingHexaToDecimal is raised when hexa parsing fails.
	ErrParsingHexaToDecimal = errors.New("Error while parsing hexa to decimal")

	// ErrParsingBinaryToDecimal is raised when binary parsing fails.
	ErrParsingBinaryToDecimal = errors.New("Error while parsing binary to decimal")

	// ErrInvalidCharacter is raised when an invalid character is used.
	ErrInvalidCharacter = errors.New("Invalid character")
)

const pngHeader = "\x89PNG\r\n\x1a\n"
const fakeType = "rmll"

// BlockSize ...
const BlockSize = 16

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// PNGHide ...
type PNGHide struct {
	// Key is the key.
	Key string
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// NewPNGHide returns a new PNGHide.
func NewPNGHide(key string) *PNGHide {
	p := &PNGHide{
		Key: key,
	}

	return p
}

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// Hide the image1 in the image2.
func (p *PNGHide) Hide(img1, img2 string) ([]byte, error) {
	// Padding of I1
	file1, err := ioutil.ReadFile(img1)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file1 : %s", err)
	}

	// Process the size
	size := len(file1) - BlockSize
	fmt.Println("The hex size of the file is :", strconv.FormatInt(int64(size), 16))

	// Create C1
	var c1 bytes.Buffer
	c1.WriteString(pngHeader)
	u := uint32(size)
	c1.WriteByte(uint8(u >> 24))
	c1.WriteByte(uint8(u >> 16))
	c1.WriteByte(uint8(u >> 8))
	c1.WriteByte(uint8(u >> 0))
	c1.WriteString(fakeType)
	fmt.Println("C1 is :", c1.String())

	// Decrypt
	c1Decrypted := decryptAes128Ecb(c1.Bytes(), []byte(p.Key))

	fmt.Println("P2 decrypted is :", c1Decrypted)

	// Create P1
	p1 := file1[:BlockSize]
	fmt.Println("P1 is :", string(p1))

	// XOR P2 with P1
	iv, err := XORBytes(c1Decrypted, []byte(p1))
	if err != nil {
		return nil, fmt.Errorf("Error while xoring the arrays : %s", err)
	}
	fmt.Println("xored is :", iv)

	// Pad file1
	file1Padded, _ := padding(file1, 16)
	fmt.Println(len(file1), len(file1Padded))

	// AES encrypt file1
	_, err = encryptCBC([]byte(p.Key), file1Padded)
	//mt.Println(string(file1Encrypted))

	return nil, nil
}
