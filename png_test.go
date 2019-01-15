package goangecryption

import (
	"testing"
	"io/ioutil"
	"bytes"
)

// HidePNG testing
func TestHidePNG(t *testing.T) {
	// Create the AC
	ga := NewGoAngecryption("AngeCryptionKey!")

	// Hide the image
	iv, err := ga.HidePNG("example/googleAndDuck/google.png", "example/googleAndDuck/duckduckgo.png", "example/googleAndDuck/result1.png")
	if err != nil {
		t.Fatal("Error while hidding :", err)
	}

	f1, err := ioutil.ReadFile("example/googleAndDuck/result1.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err := ioutil.ReadFile("example/googleAndDuck/hide.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Hidden image is not correct")
	}

	// Reveal the image
	err = ga.Reveal("example/googleAndDuck/hide.png", iv, "example/googleAndDuck/result2.png")
	if err != nil {
		t.Fatal("Error while revealing the PNG image :", err)
	}

	f1, err = ioutil.ReadFile("example/googleAndDuck/result2.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err = ioutil.ReadFile("example/googleAndDuck/reveal.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Revealed image is not correct")
	}
}