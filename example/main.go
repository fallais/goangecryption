package main

import (
	"flag"
	"fmt"
	"time"

	"goangecryption"

	"github.com/sirupsen/logrus"
)

var (
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
	ga := goangecryption.NewGoAngecryption("alpacaAndKoala!!")

	// Hide the image
	iv, err := ga.HidePNG("koala.png", "alpaca.png", "hide.png")
	if err != nil {
		logrus.Fatalln("Error while hidding :", err)
	}

	logrus.WithFields(logrus.Fields{
    "IV": fmt.Sprintf("%x", iv),
  }).Infoln("The PNG image has been hidden")

	// Reveal the image
	err = ga.Reveal("hide.png", iv, "reveal.txt")
	if err != nil {
		logrus.Fatalln("Error while revealing the PNG image :", err)
	}

	logrus.Infoln("The PNG image has been revealed")
}
