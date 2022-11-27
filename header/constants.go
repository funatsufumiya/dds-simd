package header

// private size flags
const (
	headerSize      = 124 // Size of DDS_HEADER structure
	pixelFormatSize = 32  // Size of DDS_PIXELFORMAT structure
)

// flags are related o the header itself
const (
	Caps        = 0x1
	Height      = 0x2
	Width       = 0x4
	Pitch       = 0x8
	PixelFormat = 0x1000
	MipMapCount = 0x20000
	LinearSize  = 0x80000
	Depth       = 0x800000
)

// combined header flags
const (
	HeaderFlagsTexture    = Caps | Height | Width | PixelFormat
	HeaderFlagsMipMap     = MipMapCount
	HeaderFlagsVolume     = Depth
	HeaderFlagsPitch      = Pitch
	HeaderFlagsLinearSize = LinearSize
)

// flags describe the pixel format a lot more
const (
	AlphaPixels = 0x1
	Alpha       = 0x2
	FourCC      = 0x4
	RGB         = 0x40
	YUV         = 0x200
	Luminance   = 0x20000
)
