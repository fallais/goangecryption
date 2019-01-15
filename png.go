package goangecryption

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
)

//------------------------------------------------------------------------------
// Functions
//------------------------------------------------------------------------------

// HidePNG hides the image1 into the image2.
func (p *GoAngecryption) HidePNG(img1, img2, dst string) ([]byte, error) {
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
	size := len(file1Padded) - BlockSize

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
	c1Decrypted, err := decryptECB(c1.Bytes(), []byte(p.Key))
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting the file with AES-ECB : %s", err)
	}

	// Create P1
	p1 := file1[:BlockSize]

	// XOR C1 with P1
	iv, err := xorBytes(c1Decrypted, []byte(p1))
	if err != nil {
		return nil, fmt.Errorf("Error while xoring the arrays : %s", err)
	}

	// Encrypt file1 with AES-CBC
	file1Encrypted, err := encryptCBC([]byte(p.Key), iv, file1Padded)
	if err != nil {
		return nil, fmt.Errorf("Error while encrypting the file1 : %s", err)
	}

	// Calculate the CRC32
	crc := crc32.NewIEEE()
	crc.Write(file1Encrypted[12:])
	var bu bytes.Buffer
	err = binary.Write(&bu, binary.BigEndian, int32(crc.Sum32()))
	if err != nil {
		return nil, fmt.Errorf("Error while writing the CRC32 in binary buffer : %s", err)
	}
	file1Encrypted = append(file1Encrypted, bu.Bytes()...)

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
	final, err := decryptCBC(finalPadded, []byte(p.Key), iv)
	if err != nil {
		return nil, fmt.Errorf("Error while decrupting the final file with AES-CBC : %s", err)
	}

	// Write the result file
	err = ioutil.WriteFile(dst, final, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error while writing the final file : %s", err)
	}

	return iv, nil
}

// RevealPNG reveals the image hidden in the image.
func (p *GoAngecryption) RevealPNG(img string, iv []byte, dst string) error {
	// Read the img
	file, err := ioutil.ReadFile(img)
	if err != nil {
		return fmt.Errorf("Error while reading the file : %s", err)
	}

	// Write the encrypted file
	fileEncrypted, err := encryptCBC([]byte(p.Key), iv, file)
	if err != nil {
		return fmt.Errorf("Error while encrypting the file : %s", err)
	}

	err = ioutil.WriteFile(dst, fileEncrypted, 0644)
	if err != nil {
		return fmt.Errorf("Error while writing the final file : %s", err)
	}

	return nil
}