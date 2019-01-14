# imagehide

Based on the great work of `Ange Albertini` : https://www.youtube.com/watch?v=iIesDpv9F4s

## How it works ?

### Prerequisites

Key : `alpacaAndKoala!!`

I1 : first image

![I1](https://github.com/fallais/pnghide/blob/master/example/alpaca.png)

I2 : second image

![I2](https://github.com/fallais/pnghide/blob/master/example/koala.png)

## Step 1 : determine the IV

In order to determine the first encrypted block :

- Open the `img1`
- Right padding of the `img1`
- Calculate the size of the `img1` and substract `16` (which is the **BlockSize**)
- Create the block : **PNG Header** +  **Size** + **Fake Type (rmmll)**
- Decrypt the block with **AES-ECB**
- XOR this block with the first `16 bytes` of the `img1`

`IV` is the result : `56a26af016bfac33f529597c35ad977a`.  
And the key is still : `alpacaAndKoala!!`.
