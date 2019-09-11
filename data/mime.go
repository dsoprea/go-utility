package ridata

import (
	"io"

	"net/http"

	"github.com/dsoprea/go-logging"
)

const (
	MimetypeLeadBytesCount = 512
)

// GetMimetypeFromContent uses net/http to map from magic-bytes to mime-type.
func GetMimetypeFromContent(r io.Reader, fileSize int64) (mimetype string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): !! Add test.

	leadCount := int64(MimetypeLeadBytesCount)
	if fileSize > 0 && fileSize < leadCount {
		leadCount = fileSize
	}

	buffer := make([]byte, leadCount)

	n, err := io.ReadFull(r, buffer)
	if err != nil {
		// We can return EOF if a) we weren't given a filesize and the file did
		// not haveat least as many bytes as we check by default, or b) the file-
		// size is actually (0).
		if err == io.EOF {
			return "", err
		}

		log.Panic(err)
	}

	buffer = buffer[:n]

	// Always returns a valid mime-type.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
