package goangecryption

import (
	"bytes"
	"io/ioutil"
	"testing"
)

// HideJPG testing
func TestHideJPG(t *testing.T) {
	// Create the AC
	ga := NewGoAngecryption("AngeCryptionKey!")

	// Hide the image
	iv, err := ga.HideJPG("example/jpg/alpaca.jpg", "example/jpg/koala.jpg", "example/jpg/result1.jpg")
	if err != nil {
		t.Fatal("Error while hidding :", err)
	}

	f1, err := ioutil.ReadFile("example/jpg/result1.jpg")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err := ioutil.ReadFile("example/jpg/hide.jpg")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Hidden image is not correct")
	}

	// Reveal the image
	err = ga.Reveal("example/jpg/hide.jpg", iv, "example/jpg/result2.jpg")
	if err != nil {
		t.Fatal("Error while revealing the JPG image :", err)
	}

	f1, err = ioutil.ReadFile("example/jpg/result2.jpg")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err = ioutil.ReadFile("example/jpg/reveal.jpg")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Revealed image is not correct")
	}
}
