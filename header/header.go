package header

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// init checks the integrity of the dds header and
func init() {
	if binary.Size(deserializer{}) != 128 {
		panic("dds header definition of wrong size")
	}
}

// deserializer is used to parse all header bytes into a structure
type deserializer struct {
	MagicNumber       uint32     // magic number must be "DDS "
	HeaderSize        uint32     // header size. must be 124
	TextureHeader                // the texture header
	_                 [11]uint32 // reserved1
	PixelFormatSize   uint32     // pixel format size. must be 32
	PixelFormatHeader            // the pixel format header
	Caps              [4]uint32  //
	_                 [1]uint32  //reserved2
}

func New(r io.Reader) (*Header, error) {
	var (
		buf = make([]byte, 128)
		car = new(deserializer)
	)

	if n, err := r.Read(buf); err != nil {
		return nil, err
	} else if n != 128 {
		return nil, errors.New("not enough bytes for the file header")
	} else if err = car.Read(buf); err != nil {
		return nil, err
	}

	return &Header{TextureHeader: car.TextureHeader, PixelFormatHeader: car.PixelFormatHeader, Caps: car.Caps}, nil
}

func (h *deserializer) Read(buf []byte) (err error) {
	err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, h)
	if err == nil {
		err = h.Verify()
	}
	return
}

func (h *deserializer) Verify() error {
	if mn := toString(h.MagicNumber); mn != "DDS " {
		return fmt.Errorf("magic is incorrect, expected \"DDS \", got %v", mn)
	}
	if h.HeaderSize != headerSize {
		return fmt.Errorf("DDS_HEADER reports wrong size, expected %d, got %d", headerSize, h.HeaderSize)
	}
	if h.PixelFormatSize != pixelFormatSize {
		return fmt.Errorf("DDS_PIXEL_FORMAT reports wrong size, expected %d, got %d", pixelFormatSize, h.PixelFormatSize)
	}
	// check that flags is valid
	if h.TextureFlags&HeaderFlagsTexture != HeaderFlagsTexture {
		return fmt.Errorf("DDS_HEADER reports that one or more required fields are not set: flags was %x; should at least have %x set", h.TextureFlags, HeaderFlagsTexture)
	}
	return nil
}
