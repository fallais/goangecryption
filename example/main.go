package main

import (
	"flag"
	"time"
	"fmt"

	"goangecryption"

	"github.com/sirupsen/logrus"
)

var (
	img1    = flag.String("img1", "koala.png", "First image path")
	img2    = flag.String("img2", "alpaca.png", "Second image path")
	key     = flag.String("key", "alpacaAndKoala!!", "Key")
	logging = flag.String("logging", "info", "Logging level")
)

func init() {
	// Parse the flags
	flag.Parse()

	// Set localtime to UTC
	time.Local = time.UTC

	// Set the logging level
	level, err := logrus.ParseLevel(*logging)
	if err != nil {
		logrus.Fatalln("Invalid log level ! (panic, fatal, error, warn, info, debug)")
	}
	logrus.SetLevel(level)

	// Set the TextFormatter
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})
}

func main() {
	// Create the AC
	ga := goangecryption.NewGoAngecryption(*key)

	// Hide the image
	iv, err := ga.HidePNG(*img1, *img2, "hide.png")
	if err != nil {
		logrus.Fatalln("Error while hidding :", err)
	}

	logrus.Infoln("The image has been hidden, the IV is :", fmt.Sprintf("%x", iv))

	// Reveal the image
	err = ga.RevealPNG("hide.png", iv, "reveal.png")
	if err != nil {
		logrus.Fatalln("Error while revealing :", err)
	}

	logrus.Infoln("The image has been revealed")
}
