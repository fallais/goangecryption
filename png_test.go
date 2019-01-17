package goangecryption

import (
	"testing"
	"io/ioutil"
	"bytes"
	"fmt"
)

// HidePNG testing
func TestHidePNG(t *testing.T) {
	// Create the AC
	ga := NewGoAngecryption("AngeCryptionKey!")

	// Hide the image
	iv, err := ga.HidePNG("example/png/google.png", "example/png/duckduckgo.png", "example/png/result1.png")
	if err != nil {
		t.Fatal("Error while hidding :", err)
	}

	fmt.Println("IV is :", fmt.Sprintf("%x", iv))
	t.Log("IV is :", fmt.Sprintf("%x", iv))

	f1, err := ioutil.ReadFile("example/png/result1.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err := ioutil.ReadFile("example/png/hide.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Hidden image is not correct")
	}

	// Reveal the image
	err = ga.Reveal("example/png/hide.png", iv, "example/png/result2.png")
	if err != nil {
		t.Fatal("Error while revealing the PNG image :", err)
	}

	f1, err = ioutil.ReadFile("example/png/result2.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err = ioutil.ReadFile("example/png/reveal.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Revealed image is not correct")
	}

	/////////////////////////////////////////////////////

	// Create the AC
	ga2 := NewGoAngecryption("alpacaAndKoala!!")

	// Hide the image
	iv, err = ga2.HidePNG("example/png2/alpaca.png", "example/png2/koala.png", "example/png2/result1.png")
	if err != nil {
		t.Fatal("Error while hidding :", err)
	}

	f1, err = ioutil.ReadFile("example/png2/result1.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err = ioutil.ReadFile("example/png2/hide.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Hidden image is not correct")
	}

	// Reveal the image
	err = ga2.Reveal("example/png2/hide.png", iv, "example/png2/result2.png")
	if err != nil {
		t.Fatal("Error while revealing the PNG image :", err)
	}

	f1, err = ioutil.ReadFile("example/png2/result2.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	f2, err = ioutil.ReadFile("example/png2/reveal.png")
	if err != nil {
		t.Fatal("Error while opening file :", err)
	}

	if !bytes.Equal(f1, f2) {
		t.Fatal("Revealed image is not correct")
	}
}