# goangecryption

This library is the **Golang** version of the great work of [Ange Albertini](https://corkami.github.io/).  

> It could be useful to read the [PNG specifications](http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html) before to start.

## How it works ? (PNG)

### Prerequisites

**Key** (length must be *16*) : `AngeCryptionKey!`

**FakeChunkType** (whatever the 4 chars string you want, lowercase) : `rmll`

**First image** :

![I1](https://github.com/fallais/goangecryption/blob/master/example/png/google.png)

**Second image** :

![I2](https://github.com/fallais/goangecryption/blob/master/example/png/duckduckgo.png)

### Step 1 : determine the IV

In order to determine the first encrypted block :

- Open the `img1`
- Right padding of the `img1` (modulo *16*)
- Calculate the size of the `img1` and substract `16` (which is the **BlockSize**)
- Create the block : **PNG Header** +  **Size** + **Fake Type (rmll)**
- Decrypt the block with **AES-ECB**
- XOR this block with the first `16 bytes` of the `img1`

`IV` is the result : `78d002816ba7c3de88de568f6a591d06`. And the key is still : `AngeCryptionKey!`.

### Step 2 : prepare the storage chunk

In order to prepare the storage chunk :

- Encrypt `img1` with **AES-CBC** and the **key** and **IV** we calculated previously
- Calculate the CRC32 of the encrypted data of the `img1` and append it to the encrypted `img1`

### Step 3 : generate the result

In order to generate the result :

- Open the `img2`
- Right padding of the `img2` (modulo *16*)
- Append the `img2` to the to the encrypted `img1` (except the first 8 bytes not needed because they are the *PNG Header*)
- Right padding of the result (modulo *16*)
- Decrypt the result with **AES-CBC** and the **key** and **IV** we calculated previously
- Write the result into a file
- Provide the `IV` to allow the reversed operation

### Step 4 : reversed operation

The reversed operation can be achieved by encrypting the image with **AES-CBC** and the **key** and **IV** provided.

## How to use it ?

This library can be used as follow.

```go

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/fallais/goangecryption"
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

```