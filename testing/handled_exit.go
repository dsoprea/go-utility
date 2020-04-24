package ritesting

import (
	"os"

	"github.com/dsoprea/go-logging"
)

var (
	exitsAreMarshaled = false
)

// EnableMarshaledExits enables exit marshaling.
func EnableMarshaledExits() {
	exitsAreMarshaled = true
}

// DisableMarshaledExits disables exit marshaling.
func DisableMarshaledExits() {
	exitsAreMarshaled = false
}

// Exit will marshal an exit into a panic if marshaling is turned on. If not
// turned on, forward through to `os.Exit()`.
func Exit(returnCode, successCode int) {
	if exitsAreMarshaled == false {
		os.Exit(returnCode)
	}

	if returnCode == successCode {
		return
	}

	log.Panicf("marshaled exit of (%d)", returnCode)
}
