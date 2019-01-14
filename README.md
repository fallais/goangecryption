# go-angecryption

This library is the **Golang** version of the great work of [Ange Albertini](https://corkami.github.io/).  

## How it works ?

### Prerequisites

**Key** (16 bytes long) : `alpacaAndKoala!!`

**FakeChunkType** : `ilym`

**First image** :

![I1](https://github.com/fallais/go-angecryption/blob/master/example/alpaca.png)

**Second image** :

![I2](https://github.com/fallais/go-angecryption/blob/master/example/koala.png)

> A PNG chunk is composed of : **Size (4-byte)** | **Name (4-byte)** | **Data (n-byte)** | **CRC32 (4-byte)**

### Step 1 : determine the IV

In order to determine the first encrypted block :

- Open the `img1`
- Right padding of the `img1`
- Calculate the size of the `img1` and substract `16` (which is the **BlockSize**)
- Create the block : **PNG Header** +  **Size** + **Fake Type (rmmll)**
- Decrypt the block with **AES-ECB**
- XOR this block with the first `16 bytes` of the `img1`

`IV` is the result : `56a26af016bfac33f529597c35ad977a`.  
And the key is still : `alpacaAndKoala!!`.

### Step 2 : prepare the storage chunk

In order to prepare the storage chunk :

- Encrypt `img1` with **AES-CBC** and the **key** and **IV** we calculated previously
- Calculate the CRC32 of the encrypted data of the `img1` and append it to the encrypted `img1`

### Step 3 : generate the result

In order to generate the result :

- Open the `img2`
- Right padding of the `img2`
- Append the `img2` to the to the encrypted `img1` (except the first 8 bytes)
- Right padding of the result
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
	"time"

	"github.com/fallais/goangecryption"

	"github.com/sirupsen/logrus"
)

var (
	img1    = flag.String("img1", "koala.png", "First image path")
	img2    = flag.String("img2", "alpaca.png", "Second image path")
	key = flag.String("key", "alpacaAndKoala!!", "Key")
)

func main() {
	ph := goangecryption.NewPNGHide(*key)
	_, err := ph.Hide(*img1, *img2)
	if err != nil {
		logrus.Fatalln("Error while hidding :", err)
	}
}
```