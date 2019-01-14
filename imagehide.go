package goangecryption

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
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

// PNGHeader is the PNG magic header.
const PNGHeader = "\x89PNG\r\n\x1a\n"

// FakeChunkType is the fake PNG chunk type.
const FakeChunkType = "ilym"

// BlockSize is the size of the AES block.
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
	// Read the img1
	file1, err := ioutil.ReadFile(img1)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file1 : %s", err)
	}

	// Right padding of the img1
	file1Padded, err := padding(file1, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the file1 : %s", err)
	}

	// Process the size
	fmt.Println("The decimal size of the file is :", len(file1Padded))
	size := len(file1Padded) - BlockSize
	fmt.Println("The decimal size of the file after substracting block size is :", size)
	fmt.Println("The hex size of the file is :", strconv.FormatInt(int64(size), 16))

	// Create C1
	var c1 bytes.Buffer
	c1.WriteString(PNGHeader)
	u := uint32(size)
	c1.WriteByte(uint8(u >> 24))
	c1.WriteByte(uint8(u >> 16))
	c1.WriteByte(uint8(u >> 8))
	c1.WriteByte(uint8(u >> 0))
	c1.WriteString(FakeChunkType)

	// Decrypt C1 with AES-ECB
	c1Decrypted := decryptECB(c1.Bytes(), []byte(p.Key))

	// Create P1
	p1 := file1[:BlockSize]

	// XOR P2 with P1
	iv, err := xorBytes(c1Decrypted, []byte(p1))
	if err != nil {
		return nil, fmt.Errorf("Error while xoring the arrays : %s", err)
	}
	fmt.Println("IV is :", fmt.Sprintf("%x", iv))
	fmt.Println("Key is :", p.Key)

	// Encrypt file1 with AES-CBC
	file1Encrypted, err := encryptCBC([]byte(p.Key), iv, file1Padded)
	if err != nil {
		return nil, fmt.Errorf("Error while encrypting the file1 : %s", err)
	}

	// Calculate the CRC32
	crc := crc32.NewIEEE()
	crc.Write(file1Encrypted[12:])
	var bufi bytes.Buffer
	err = binary.Write(&bufi, binary.BigEndian, int32(crc.Sum32()))
	if err != nil {
		return nil, fmt.Errorf("Error while writing the CRC32 in binary buffer : %s", err)
	}
	file1Encrypted = append(file1Encrypted, bufi.Bytes()...)

	// Read the img2
	file2, err := ioutil.ReadFile(img2)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file2 : %s", err)
	}

	// Right padding of the img1
	file2Padded, err := padding(file2, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the file2 : %s", err)
	}

	// Append the img2
	file1Encrypted = append(file1Encrypted, file2Padded[8:]...)

	// Right padding of the result
	finalPadded, err := padding(file1Encrypted, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the final file : %s", err)
	}

	// Decrypt the result with AES-CBC
	final := decryptCBC(finalPadded, []byte(p.Key), iv)

	// Write the result file
	err = ioutil.WriteFile("final.png", final, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error while writing the final file : %s", err)
	}

	// Write the encrypted file
	finalEncrypt, _ := encryptCBC([]byte(p.Key), iv, final)
	err = ioutil.WriteFile("final_enc.png", finalEncrypt, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error while writing the final file : %s", err)
	}

	return nil, nil
}
