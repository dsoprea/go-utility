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
func Exit(returnCode int) {
	if exitsAreMarshaled == false {
		os.Exit(returnCode)
	}

	// NOTE(dustin): If we need non-zero success-codes later, we can add a new function for it. Since this is the common case, we'd like this function to have the simplest signature.

	if returnCode == 0 {
		return
	}

	log.Panicf("marshaled exit of (%d)", returnCode)
}
