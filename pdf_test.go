package goangecryption

import (
	"testing"
	"io/ioutil"
	"bytes"
)

// HidePNG testing
func TestHidePDF(t *testing.T) {
	// Create the AC
	ga := NewGoAngecryption("IsAESbrokenYET ?")
	iv := []byte(")\x97\x89\xf9p\xf7\x15\xb9$\xe0TC\xd7\xbc\x10\xab")

	// Reveal the image
	err := ga.Reveal("example/pdf/hide.pdf", iv, "example/pdf/result2.pdf")
	if err != nil {
		t.Fatal("Error while revealing the PNG image :", err)
	}

	f1, err := ioutil.ReadFile("example/pdf/result2.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err := ioutil.ReadFile("example/pdf/reveal.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Revealed image is not correct")
	}
}