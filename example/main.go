package main

import (
	"flag"
	"fmt"
	"time"

	"goangecryption"

	"github.com/sirupsen/logrus"
)

var (
	source  = flag.String("source", "png/google.png", "Source")
	target  = flag.String("target", "png/google.png", "Target")
	result  = flag.String("result", "png/result.png", "Result")
	action  = flag.String("action", "hide", "Action (hide, reveal)")
	method  = flag.String("method", "png", "Method (png, jpg, flv, pdf)")
	key     = flag.String("key", "alpacaAndKoala!!", "Key")
	iv      = flag.String("iv", "2e45d76068c706e5a1acd82467fe431c", "IV")
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

	// Select the method
	switch *method {
	case "png":
		switch *action {
		case "hide":
			// Hide the image
			iv, err := ga.HidePNG(*source, *target, *result)
			if err != nil {
				logrus.Fatalln("Error while hidding :", err)
			}

			logrus.WithFields(logrus.Fields{
				"IV": fmt.Sprintf("%x", iv),
			}).Infoln("The source has been hidden in the target")
			break
		case "reveal":
			// Reveal the image
			err := ga.Reveal(*source, []byte(*iv), *result)
			if err != nil {
				logrus.Fatalln("Error while revealing the source image :", err)
			}

			logrus.Infoln("The PNG image has been revealed")
			break
		default:
			logrus.Fatalln("The action is not valid")
		}
	}
}
