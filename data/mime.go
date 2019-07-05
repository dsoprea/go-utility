package ridata

import (
	"io"

	"net/http"

	"github.com/dsoprea/go-logging"
)

const (
	MimetypeLeadBytes = 512
)

// GetMimetypeFromContent uses net/http to map from magic-bytes to mime-type.
func GetMimetypeFromContent(r io.Reader, fileSize int) (mimetype string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): !! Add test.

	leadCount := MimetypeLeadBytes
	if fileSize > 0 && fileSize < leadCount {
		leadCount = fileSize
	}

	buffer := make([]byte, leadCount)

	n, err := io.ReadFull(r, buffer)
	log.PanicIf(err)

	buffer = buffer[:n]

	// Always returns a valid mime-type.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
