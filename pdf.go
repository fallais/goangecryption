package goangecryption

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

//------------------------------------------------------------------------------
// Functions
//------------------------------------------------------------------------------

// HideInPDF hides the source into the target and writes the result file.
func (p *GoAngecryption) HideInPDF(s, t, r string) ([]byte, error) {
	// Read the source
	source, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file1 : %s", err)
	}

	// Right padding of the source
	sourcePadded, err := padding(source, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the file1 : %s", err)
	}

	// Read the target
	target, err := ioutil.ReadFile(t)
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file2 : %s", err)
	}

	// Right padding of the target
	targetPadded, err := padding(target, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the file2 : %s", err)
	}

	// Create C1
	var c1 bytes.Buffer
	c1.WriteString("%PDF-\x00obj\nstream")

	// Decrypt C1 with AES-ECB
	c1Decrypted, err := decryptECB(c1.Bytes(), []byte(p.Key))
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting the file with AES-ECB : %s", err)
	}

	// XOR C1 with P1
	iv, err := xorBytes(c1Decrypted, []byte(source[:BlockSize]))
	if err != nil {
		return nil, fmt.Errorf("Error while xoring the arrays : %s", err)
	}

	// Encrypt
	result, err := encryptCBC([]byte(p.Key), iv, sourcePadded)
	if err != nil {
		return nil, fmt.Errorf("Error while encrypting the file : %s", err)
	}

	// Append the file2
	result = append(result, []byte("\nendstream\nendobj\n")...)
	result = append(result, targetPadded...)

	// Right padding of the result
	resultPadded, err := padding(result, 16)
	if err != nil {
		return nil, fmt.Errorf("Error while padding the result : %s", err)
	}

	// Decrypt the result with AES-CBC
	final, err := decryptCBC(resultPadded, []byte(p.Key), iv)
	if err != nil {
		return nil, fmt.Errorf("Error while decrypting the final file with AES-CBC : %s", err)
	}

	// Write the result file
	err = ioutil.WriteFile(r, final, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error while writing the result file : %s", err)
	}

	return iv, nil
}