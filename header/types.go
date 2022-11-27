package header

import "encoding/binary"

type (
	Header struct {
		TextureHeader
		PixelFormatHeader
		Caps [4]uint32
	}

	TextureHeader struct {
		TextureFlags      uint32 // flag-set
		Height            uint32 // width of the texture in pixels
		Width             uint32 // height of the texture in pixels
		PitchOrLinearSize uint32 // total size of the image
		Depth             uint32
		MipMapCount       uint32
	}

	PixelFormatHeader struct {
		PixelFlags  uint32
		FourCC      uint32 // code for the used texture compression
		RgbBitCount uint32 // bit-size for every color
		RBitMask    uint32
		GBitMask    uint32
		BBitMask    uint32
		ABitMask    uint32
	}
)

func (h *Header) FourCCString() string {
	return toString(h.FourCC)
}

func toString(i uint32) string {
	return string(binary.LittleEndian.AppendUint32(nil, i))
}
