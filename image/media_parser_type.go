package riimage

import (
	"io"

	"github.com/dsoprea/go-exif/v2"
)

// MediaContext is an accessor that knows how to extract specific metadata from
// the media.
type MediaContext interface {
	// Exif returns the EXIF's root IFD.
	Exif() (rootIfd *exif.Ifd, data []byte, err error)
}

// MediaParser prescribes a specific structure for the parser types that are
// imported from other projects. We don't use it directly, but we use this to
// impose structure.
type MediaParser interface {
	// Parse parses a stream using an `io.Reader`. `ec` should *actually* be a
	// `ExifContext`.
	Parse(r io.ReadSeeker, size int) (ec MediaContext, err error)

	// ParseFile parses a stream using a file. `ec` should *actually* be a
	// `ExifContext`.
	ParseFile(filepath string) (ec MediaContext, err error)

	// ParseBytes parses a stream direct from bytes. `ec` should *actually* be
	// a `ExifContext`.
	ParseBytes(data []byte) (ec MediaContext, err error)

	// Parses the data to determine if it's a compatible format.
	LooksLikeFormat(data []byte) bool
}
