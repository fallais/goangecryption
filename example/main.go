package main

import (
	"flag"
	"time"

	"goangecryption"

	"github.com/sirupsen/logrus"
)

var (
	img1    = flag.String("img1", "koala.png", "First image path")
	img2    = flag.String("img2", "alpaca.png", "Second image path")
	key = flag.String("key", "alpacaAndKoala!!", "Key")
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
	ph := goangecryption.NewPNGHide(*key)
	_, err := ph.Hide(*img1, *img2)
	if err != nil {
		logrus.Fatalln("Error while hidding :", err)
	}
}
