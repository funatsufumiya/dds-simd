package header

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	sizeDDTF = 128 // Size of the whole texture file header. is 128
	sizeDDSD = 124 // Size of the serialized DDSHeader. is 124
	sizeDDPF = 32  // Size of the serialized DDPFHeader. is 32
)

// Deserializer is used to parse all header bytes into a structure
type Deserializer struct {
	MagicNumber     uint32     // magic number must be "DDS "
	HeaderSize      uint32     // header size. must be 124
	DDSHeader                  // the texture header
	_               [11]uint32 // reserved1
	PixelFormatSize uint32     // pixel format size. must be 32
	DDPFHeader                 // the pixel format header
	CapsHeader                 // header for more complex textures with mipmaps or cube-maps
	_               [1]uint32  // reserved2
}

// New creates a new empty Deserializer waiting for a call to Deserializer.Read or Deserializer.Parse
func New() *Deserializer {
	return new(Deserializer)
}

// Read tries to take size_DDTF Bytes from the reader and then calls Deserializer.Parse with it.
func (d *Deserializer) Read(r io.Reader) (*Header, error) {
	buf := make([]byte, sizeDDTF, sizeDDTF)
	if n, err := r.Read(buf); err != nil {
		return nil, err
	} else if n != sizeDDTF {
		return nil, fmt.Errorf("bytes in header. expected %d, only got : %d", sizeDDTF, n)
	}
	return d.Parse(*(*[sizeDDTF]byte)(buf))
}

// Parse takes size_DDTF Bytes and tries to create a Header from it. Calls verification on a successful
// parsed Header, which might return an error in the case of a wrongly configured header.
func (d *Deserializer) Parse(in [sizeDDTF]byte) (h *Header, err error) {
	if err = binary.Read(bytes.NewReader(in[:]), binary.LittleEndian, d); err == nil {
		err = d.verify()
		h = &Header{
			DDSHeader:    d.DDSHeader,
			DDPFHeader:   d.DDPFHeader,
			CapsHeader:   d.CapsHeader,
			FourCCString: d.toString(d.FourCC),
		}
	}
	return
}

// verify makes some semantic checks for validity
func (d *Deserializer) verify() error {
	if mn := d.toString(d.MagicNumber); mn != "DDS " {
		return fmt.Errorf("magic is incorrect, expected \"DDS \", got %v", mn)
	}
	if d.HeaderSize != sizeDDSD {
		return fmt.Errorf("DDS_HEADER reports wrong size, expected %d, got %d", sizeDDSD, d.HeaderSize)
	}
	if d.PixelFormatSize != sizeDDPF {
		return fmt.Errorf("DDS_PIXEL_FORMAT reports wrong size, expected %d, got %d", sizeDDPF, d.PixelFormatSize)
	}

	// check that it's actually a texture per requirements
	if !d.TextureFlags.Has(DDSDHeaderFlagsTexture) {
		return fmt.Errorf("DDS_HEADER reports that one or more required fields are not set: flags was %x; should at least have %x set", d.TextureFlags, DDSDHeaderFlagsTexture)
	}
	return nil
}

func (*Deserializer) toString(i uint32) string {
	return string(binary.LittleEndian.AppendUint32(nil, i))
}
