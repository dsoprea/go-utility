package ritesting

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestRedirectTty(t *testing.T) {
	if rtty != nil {
		t.Fatalf("rtty global is not nil")
	}

	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	RedirectTty()

	if rtty == nil {
		t.Fatalf("rtty is nil after redirect.")
	}

	if rtty.originalStdin != originalStdin {
		t.Fatalf("Captured STDIN is not correct.")
	} else if rtty.originalStdout != originalStdout {
		t.Fatalf("Captured STDOUT is not correct.")
	} else if rtty.originalStderr != originalStderr {
		t.Fatalf("Captured STDERR is not correct.")
	}

	// Load test data.

	testStdinData := []byte("test input")

	_, err := rtty.newStdinWriter.Write(testStdinData)
	log.PanicIf(err)

	testStdoutData := []byte("test stdout output")

	_, err = os.Stdout.Write(testStdoutData)
	log.PanicIf(err)

	testStderrData := []byte("test stderr output")

	_, err = os.Stderr.Write(testStderrData)
	log.PanicIf(err)

	// Validate.

	stdinBuffer := make([]byte, len(testStdinData))
	_, err = io.ReadFull(os.Stdin, stdinBuffer)
	log.PanicIf(err)

	if bytes.Equal(stdinBuffer, testStdinData) != true {
		t.Fatalf("STDIN data not correct.")
	}

	stdoutBuffer := make([]byte, len(testStdoutData))
	_, err = io.ReadFull(rtty.newStdoutReader, stdoutBuffer)
	log.PanicIf(err)

	if bytes.Equal(stdoutBuffer, testStdoutData) != true {
		t.Fatalf("STDOUT data not correct.")
	}

	stderrBuffer := make([]byte, len(testStderrData))
	_, err = io.ReadFull(rtty.newStderrReader, stderrBuffer)
	log.PanicIf(err)

	if bytes.Equal(stderrBuffer, testStderrData) != true {
		t.Fatalf("STDERR data not correct.")
	}

	RestoreTty()

	if os.Stdin != originalStdin {
		t.Fatalf("Restored STDIN is not correct.")
	} else if os.Stdout != originalStdout {
		t.Fatalf("Restored STDOUT is not correct.")
	} else if os.Stderr != originalStderr {
		t.Fatalf("Restored STDERR is not correct.")
	}
}

func TestRestoreTty(t *testing.T) {
	if rtty != nil {
		t.Fatalf("rtty global is not nil")
	}

	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	RedirectTty()

	if rtty == nil {
		t.Fatalf("rtty is nil after redirect.")
	}

	if rtty.originalStdin != originalStdin {
		t.Fatalf("Captured STDIN is not correct.")
	} else if rtty.originalStdout != originalStdout {
		t.Fatalf("Captured STDOUT is not correct.")
	} else if rtty.originalStderr != originalStderr {
		t.Fatalf("Captured STDERR is not correct.")
	}

	RestoreTty()

	if os.Stdin != originalStdin {
		t.Fatalf("Restored STDIN is not correct.")
	} else if os.Stdout != originalStdout {
		t.Fatalf("Restored STDOUT is not correct.")
	} else if os.Stderr != originalStderr {
		t.Fatalf("Restored STDERR is not correct.")
	}
}

func TestRestoreAndDumpTty(t *testing.T) {
	if rtty != nil {
		t.Fatalf("rtty global is not nil")
	}

	originalStdin := os.Stdin
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	RedirectTty()

	if rtty == nil {
		t.Fatalf("rtty is nil after redirect.")
	}

	if rtty.originalStdin != originalStdin {
		t.Fatalf("Captured STDIN is not correct.")
	} else if rtty.originalStdout != originalStdout {
		t.Fatalf("Captured STDOUT is not correct.")
	} else if rtty.originalStderr != originalStderr {
		t.Fatalf("Captured STDERR is not correct.")
	}

	RestoreAndDumpTty()

	if os.Stdin != originalStdin {
		t.Fatalf("Restored STDIN is not correct.")
	} else if os.Stdout != originalStdout {
		t.Fatalf("Restored STDOUT is not correct.")
	} else if os.Stderr != originalStderr {
		t.Fatalf("Restored STDERR is not correct.")
	}
}

func TestStdinWriter(t *testing.T) {
	RedirectTty()

	defer RestoreTty()

	if StdinWriter() != rtty.newStdinWriter {
		t.Fatalf("Not correct resource.")
	}
}

func TestStdoutReader(t *testing.T) {
	RedirectTty()

	defer RestoreTty()

	if StdoutReader() != rtty.newStdoutReader {
		t.Fatalf("Not correct resource.")
	}
}

func TestStderrReader(t *testing.T) {
	RedirectTty()

	defer RestoreTty()

	if StderrReader() != rtty.newStderrReader {
		t.Fatalf("Not correct resource.")
	}
}

func TestIsTtyRedirected(t *testing.T) {
	if IsTtyRedirected() != false {
		t.Fatalf("Incorrect initial state.")
	}

	RedirectTty()

	if IsTtyRedirected() != true {
		t.Fatalf("Incorrect redirected state.")
	}

	RestoreTty()

	if IsTtyRedirected() != false {
		t.Fatalf("Incorrect final state.")
	}
}
