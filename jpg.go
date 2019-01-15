package goangecryption

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

//------------------------------------------------------------------------------
// Functions
//------------------------------------------------------------------------------

// HideJPG hides the image1 into the image2.
func (p *GoAngecryption) HideJPG(img1, img2, dst string) ([]byte, error) {
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
	c1.WriteString(JPGHeader)
	c1.WriteString("\xFF\xFE")
	u := uint16(size)
	c1.WriteByte(uint8(u >> 8))
	c1.WriteByte(uint8(u >> 0))
	c1.WriteString(strings.Repeat("\x00", 10))

	// Decrypt C1 with AES-ECB
	c1Decrypted, err := decryptECB(c1.Bytes(), []byte(p.Key))
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting the file with AES-ECB : %s", err)
	}

	// XOR C1 with P1
	iv, err := xorBytes(c1Decrypted, []byte(file1[:BlockSize]))
	if err != nil {
		return nil, fmt.Errorf("Error while xoring the arrays : %s", err)
	}

	// Encrypt
	result, err := encryptCBC([]byte(p.Key), iv, file1Padded)
	if err != nil {
		return nil, fmt.Errorf("Error while encrypting the file : %s", err)
	}

	// Read the img2
	file2, err := ioutil.ReadFile(img2)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file2 : %s", err)
	}

	// Right padding of the img2
	file2Padded, err := padding(file2, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the file2 : %s", err)
	}

	// Append
	result = append(result, file2Padded[2:]...)

	// Right padding of the result
	resultPadded, err := padding(result, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the result : %s", err)
	}

	// Decrypt the result with AES-CBC
	final, err := decryptCBC(resultPadded, []byte(p.Key), iv)
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
