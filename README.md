# imagehide

Based on the great work of `Ange Albertini` : https://www.youtube.com/watch?v=iIesDpv9F4s

## Information

I1 : first image

![I1](https://github.com/fallais/go-siem/blob/master/alpaca.png)

I2 : second image

![I2](https://github.com/fallais/go-siem/blob/master/koala.png)

Key : `alpacaAndKoala!`

## Step 0 : checks

These conditions must be respected :

- `I1` and `I2` are PNGs ;
- `I1` fits into only on chunk of `I2`.

## Step 1 : determine the first encrypted block

`R` will have the same first block as `I1`.

Once `encrypted`, `R` will start with :

- PNG magic header of `8 bytes` : `\x89PNG\r\n\x1a\n`
- A storage chunk with a fake type `rmll` that will contains `I1` with a length of `16926 - 6 = 16920 = 00004218 (hex)`

> A PNG chunk is composed of the `length`, the `type`, the `data` and the `CRC`.

First block of `R`, `P1`, comes from `I1` : `\x89PNG\r\n\x1a\n \x00\x00\x00\x0D IHDR`.  
First encrypted block of `R`, `P2` :  `\x89PNG\r\n\x1a\n \x00\x00\x42\x18 rmll`.

