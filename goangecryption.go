package goangecryption

// PNGHeader is the PNG magic header.
const PNGHeader = "\x89PNG\r\n\x1a\n"

// JPGHeader is the JPG magic header.
const JPGHeader = "\xFF\xD8"

// FakeChunkType is the fake PNG chunk type.
const FakeChunkType = "ilym"

// BlockSize is the size of the block.
const BlockSize = 16

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// GoAngecryption is the Go version of Angecryption.
type GoAngecryption struct {
	// Key is the secret key.
	Key string
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// NewGoAngecryption returns a new GoAngecryption.
func NewGoAngecryption(key string) *GoAngecryption {
	p := &GoAngecryption{
		Key: key,
	}

	return p
}
