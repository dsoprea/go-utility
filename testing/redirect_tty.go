package ritesting

import (
	"fmt"
	"os"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

type redirectedTty struct {
	originalStdin  *os.File
	newStdinWriter *os.File

	originalStdout  *os.File
	newStdoutReader *os.File

	originalStderr  *os.File
	newStderrReader *os.File
}

var (
	rtty *redirectedTty
)

// RedirectTty redirects STDIN, STDOUT, and STDERR to alternative resources.
// This allows us to create unit-tests for executables by directly calling
// their main() entrypoints.
func RedirectTty() {
	if rtty != nil {
		log.Panicf("TTY is already redirected")
	}

	rtty = &redirectedTty{
		originalStdin:  os.Stdin,
		originalStdout: os.Stdout,
		originalStderr: os.Stderr,
	}

	var err error

	// The caller can write to the new writer and it will be available on STDIN.

	os.Stdin, rtty.newStdinWriter, err = os.Pipe()
	log.PanicIf(err)

	// Any output to STDOUT or STDERR will be available to the caller on the
	// new readers.

	rtty.newStdoutReader, os.Stdout, err = os.Pipe()
	log.PanicIf(err)

	rtty.newStderrReader, os.Stderr, err = os.Pipe()
	log.PanicIf(err)
}

// RestoreTty restores original TTY resources.
func RestoreTty() {
	if rtty == nil {
		return
	}

	os.Stdin = rtty.originalStdin
	os.Stdout = rtty.originalStdout
	os.Stderr = rtty.originalStderr

	rtty = nil
}

// RestoreAndDumpTty restores original TTY resources but not before grabbing
// their screen output and then printing it before returning.
func RestoreAndDumpTty() {
	if rtty == nil {
		return
	}

	// TODO(dustin): !! Finish. Close os.Stdout and os.Stderr, read each, and print each between anchors.

	os.Stdout.Close()

	stdoutOutput, err := ioutil.ReadAll(rtty.newStdoutReader)
	log.PanicIf(err)

	os.Stderr.Close()

	stderrOutput, err := ioutil.ReadAll(rtty.newStderrReader)
	log.PanicIf(err)

	RestoreTty()

	fmt.Printf(">>>>>>>>>>>>>\n")
	fmt.Printf("STDOUT OUTPUT\n")
	fmt.Printf(">>>>>>>>>>>>>\n")
	fmt.Println(string(stdoutOutput))
	fmt.Printf("<<<<<<<<<<<<<\n")
	fmt.Printf("\n")

	fmt.Printf(">>>>>>>>>>>>>\n")
	fmt.Printf("STDERR OUTPUT\n")
	fmt.Printf(">>>>>>>>>>>>>\n")
	fmt.Println(string(stderrOutput))
	fmt.Printf("<<<<<<<<<<<<<\n")
	fmt.Printf("\n")
}

// StdinWriter returns a writer that can be used to feed STDIN data (if
// redirected).
func StdinWriter() *os.File {
	if rtty == nil {
		log.Panicf("TTY not redirected; STDIN writer not available")
	}

	return rtty.newStdinWriter
}

// StdoutReader returns a reader that can be used to read STDOUT output (if
// redirected).
func StdoutReader() *os.File {
	if rtty == nil {
		log.Panicf("TTY not redirected; STDOUT reader not available")
	}

	return rtty.newStdoutReader
}

// StderrReader returns a reader that can be used to read STDERR output (if
// redirected).
func StderrReader() *os.File {
	if rtty == nil {
		log.Panicf("TTY not redirected; STDERR reader not available")
	}

	return rtty.newStderrReader
}

// IsTtyRedirected returns whether the TTY is currently redirected.
func IsTtyRedirected() bool {
	return rtty != nil
}
