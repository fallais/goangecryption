package goangecryption

import (
	"fmt"
	"io/ioutil"
)

// Reveal reveals the hidden data in the src with given IV.
func (p *GoAngecryption) Reveal(src string, iv []byte, dst string) error {
	// Read the img
	file, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("Error while reading the file : %s", err)
	}

	// Encrypted the file with AES-CBC
	fileEncrypted, err := encryptCBC([]byte(p.Key), iv, file)
	if err != nil {
		return fmt.Errorf("Error while encrypting the file with AES-CBC : %s", err)
	}

	// Write the encrypted file in dst
	err = ioutil.WriteFile(dst, fileEncrypted, 0644)
	if err != nil {
		return fmt.Errorf("Error while writing the file : %s", err)
	}

	return nil
}