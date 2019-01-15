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
	size := len(file1) - BlockSize

	// Create C1
	var c1 bytes.Buffer
	c1.WriteString(JPGHeader)
	c1.WriteString("\xFF\xFE")
	u := uint16(size)
	c1.WriteByte(uint8(u >> 8))
	c1.WriteByte(uint8(u >> 0))
	c1.WriteString(strings.Repeat("\x00", 10))

	// Decrypt C1 with AES-ECB
	iv, err := decryptECB(c1.Bytes(), []byte(p.Key))
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting the file with AES-ECB : %s", err)
	}

	// Encrypt
	result, err := encryptCBC([]byte(p.Key), iv, file1Padded)
	if err != nil {
		return nil, fmt.Errorf("Error while encrypting the file : %s", err)
	}

	// Read the img2
	file2, err := ioutil.ReadFile(img2)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file1 : %s", err)
	}

	// Append
	result = append(result, file2[2:]...)

	// Write the result file
	err = ioutil.WriteFile(dst, result, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error while writing the final file : %s", err)
	}

	return iv, nil
}

// RevealJPG reveals the image hidden in the image.
func (p *GoAngecryption) RevealJPG(img1 string, iv []byte) ([]byte, error) {
	return nil, nil
}
