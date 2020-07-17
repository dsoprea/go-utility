package rifs

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestSeekableBuffer_Write_FromEmpty(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")
	dataLen := len(data)

	n, err := sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written: (%d)", dataLen, n)
	}

	position, err := sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	if position != 0 {
		t.Fatalf("The seek did not return a position of (%d): (%d)", dataLen, position)
	}

	buffer := make([]byte, 10)

	n, err = sb.Read(buffer)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes were not read: (%d)", dataLen, n)
	}

	buffer = buffer[:dataLen]

	if bytes.Compare(buffer, data) != 0 {
		t.Fatalf("Data did not match.")
	}

	n, err = sb.Read(buffer)
	if n != 0 {
		t.Fatalf("Read after EOF did not return (0): (%d)", n)
	}

	if err != io.EOF {
		t.Fatalf("Expected EOF: [%v]", err)
	}
}

func TestSeekableBuffer_Write_Append(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")
	dataLen := len(data)

	n, err := sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written (1): (%d)", dataLen, n)
	}

	n, err = sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written (2): (%d)", dataLen, n)
	}

	position, err := sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	if position != 0 {
		t.Fatalf("The seek did not return a position of (%d): (%d)", dataLen, position)
	}

	buffer := make([]byte, 15)

	n, err = sb.Read(buffer)
	log.PanicIf(err)

	if n != dataLen*2 {
		t.Fatalf("Exactly (%d) bytes were not read: (%d)", dataLen, n)
	}

	buffer = buffer[:dataLen*2]

	if bytes.Compare(buffer[:dataLen], data) != 0 {
		t.Fatalf("Data did not match (1).")
	}

	if bytes.Compare(buffer[dataLen:], data) != 0 {
		t.Fatalf("Data did not match (2).")
	}

	n, err = sb.Read(buffer)
	if n != 0 {
		t.Fatalf("Read after EOF did not return (0): (%d)", n)
	}

	if err != io.EOF {
		t.Fatalf("Expected EOF: [%v]", err)
	}
}

func TestSeekableBuffer_Write_Replace(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("abcde-fghij-klmno")

	_, err := sb.Write(data)
	log.PanicIf(err)

	// Replace the characters in the middle.

	_, err = sb.Seek(6, os.SEEK_SET)
	log.PanicIf(err)

	data2 := []byte("12345")

	_, err = sb.Write(data2)
	log.PanicIf(err)

	// We're still six bytes from the end of the file. Read the rest.
	buffer := make([]byte, 10)

	n, err := sb.Read(buffer)
	log.PanicIf(err)

	if n != 6 {
		t.Fatalf("Could not read exactly the six remaining bytes: (%d)", n)
	}

	if bytes.Compare(buffer[:n], data[len(data)-n:]) != 0 {
		t.Fatalf("The last six bytes are not correct:\nACTUAL:\n%v\nEXPECTED:\n%v", buffer[:n], data[len(data)-n:])
	}

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	buffer = make([]byte, 20)

	expected := []byte("abcde-12345-klmno")

	n, err = sb.Read(buffer)
	log.PanicIf(err)

	if n != len(expected) {
		t.Fatalf("Exactly the right number of bytes were not read: (%d) != (%d)", n, len(expected))
	}

	if bytes.Compare(buffer[:n], expected) != 0 {
		t.Fatalf("The complete contents are not correct:\nACTUAL: [%s]\n%v\nEXPECTED: [%s]\n%v", string(buffer[:n]), buffer[:n], string(expected), expected)
	}

	_, err = sb.Read(buffer)
	if err != io.EOF {
		t.Fatalf("Expected EOF: [%v]", err)
	}
}

func TestSeekableBuffer_Write_Expand(t *testing.T) {
	sb := NewSeekableBuffer()

	// Write first string.

	data := []byte("word1-word2")

	_, err := sb.Write(data)
	log.PanicIf(err)

	// Seek and replace partway through, and replace insert more data than we
	// currently have.

	_, err = sb.Seek(6, os.SEEK_SET)
	log.PanicIf(err)

	data2 := []byte("word3-word4")

	_, err = sb.Write(data2)
	log.PanicIf(err)

	// Read contents.

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	buffer := make([]byte, 20)

	n, err := sb.Read(buffer)
	log.PanicIf(err)

	if n != 17 {
		t.Fatalf("Final size not correct: (%d)", n)
	}

	expected := []byte("word1-word3-word4")

	if bytes.Compare(buffer[:n], expected) != 0 {
		t.Fatalf("Data did not match.")
	}

	_, err = sb.Read(buffer)
	if err != io.EOF {
		t.Fatalf("Expected EOF: [%v]", err)
	}
}

func TestSeekableBuffer_Write_ReadParts(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("word1-word2-word3")

	_, err := sb.Write(data)
	log.PanicIf(err)

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	// Read first part.

	buffer := make([]byte, 12)

	n, err := io.ReadFull(sb, buffer)
	log.PanicIf(err)

	if n != 12 {
		t.Fatalf("Read size not correct: (%d)", n)
	}

	if bytes.Compare(buffer[:n], data[:12]) != 0 {
		t.Fatalf("Data did not match (1).")
	}

	// Read rest.

	n, err = sb.Read(buffer)
	log.PanicIf(err)

	if n != 5 {
		t.Fatalf("Read size not correct: (%d)", n)
	}

	if bytes.Compare(buffer[:n], data[12:]) != 0 {
		t.Fatalf("Data did not match (2).")
	}

	_, err = sb.Read(buffer)
	if err != io.EOF {
		t.Fatalf("Expected EOF: [%v]", err)
	}
}

func TestSeekableBuffer_Write_ReadFromMiddle(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("word1-word2-word3")

	_, err := sb.Write(data)
	log.PanicIf(err)

	_, err = sb.Seek(12, os.SEEK_SET)
	log.PanicIf(err)

	// Read rest.

	buffer := make([]byte, 10)

	n, err := sb.Read(buffer)
	log.PanicIf(err)

	if n != 5 {
		t.Fatalf("Read size not correct: (%d)", n)
	}

	if bytes.Compare(buffer[:n], data[12:]) != 0 {
		t.Fatalf("Data did not match (2).")
	}

	_, err = sb.Read(buffer)
	if err != io.EOF {
		t.Fatalf("Expected EOF: [%v]", err)
	}
}

func TestSeekableBuffer_Write_Append_SeekAwayAndBack(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")
	dataLen := len(data)

	n, err := sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written (1): (%d)", dataLen, n)
	}

	// Seek somewhere else.
	position, err := sb.Seek(100, os.SEEK_SET)
	log.PanicIf(err)

	if position != 100 {
		t.Fatalf("Far seek position not correct: (%d)", position)
	}

	// Seek back to where we were.
	position, err = sb.Seek(int64(dataLen), os.SEEK_SET)
	log.PanicIf(err)

	if position != int64(dataLen) {
		t.Fatalf("Far seek position not correct: (%d) != (%d)", position, dataLen)
	}

	// Make sure the size of the buffer wasn't affected by just seeking to a
	// position that didn't exist (nothing should change until an actual write
	// happens).
	if sb.Len() != dataLen {
		t.Fatalf("Buffer size not correct: (%d)", sb.Len())
	}

	n, err = sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written (2): (%d)", dataLen, n)
	}

	position, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	if position != 0 {
		t.Fatalf("The seek did not return a position of (%d): (%d)", dataLen, position)
	}

	buffer := make([]byte, 15)

	n, err = sb.Read(buffer)
	log.PanicIf(err)

	if n != dataLen*2 {
		t.Fatalf("Exactly (%d) bytes were not read: (%d)", dataLen, n)
	}

	buffer = buffer[:dataLen*2]

	if bytes.Compare(buffer[:dataLen], data) != 0 {
		t.Fatalf("Data did not match (1).")
	}

	if bytes.Compare(buffer[dataLen:], data) != 0 {
		t.Fatalf("Data did not match (2).")
	}

	n, err = sb.Read(buffer)
	if n != 0 {
		t.Fatalf("Read after EOF did not return (0): (%d)", n)
	}

	if err != io.EOF {
		t.Fatalf("Expected EOF: [%v]", err)
	}
}

func TestSeekableBuffer_Seek_End_Empty(t *testing.T) {
	sb := NewSeekableBuffer()

	position, err := sb.Seek(0, os.SEEK_END)
	log.PanicIf(err)

	if position != 0 {
		t.Fatalf("The seek did not move to the last position of a still-empty file: (%d)", position)
	}
}

func TestSeekableBuffer_Seek_End_Nonempty_Zero(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")
	dataLen := len(data)

	n, err := sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written: (%d)", dataLen, n)
	}

	position, err := sb.Seek(0, os.SEEK_END)
	log.PanicIf(err)

	if position != int64(dataLen) {
		t.Fatalf("The seek did not move to the last byte of a non-empty byte: (%d) != (%d)", position, dataLen)
	}
}

func TestSeekableBuffer_Seek_End_Empty_SeekToNegative(t *testing.T) {
	sb := NewSeekableBuffer()

	position, err := sb.Seek(-10, os.SEEK_END)
	log.PanicIf(err)

	if position != int64(0) {
		t.Fatalf("The seek to a negative position did not end up on zero: (%d)", position)
	}
}

func TestSeekableBuffer_Seek_End_Nonempty_Nonzero(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")
	dataLen := len(data)

	n, err := sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written: (%d)", dataLen, n)
	}

	position, err := sb.Seek(-2, os.SEEK_END)
	log.PanicIf(err)

	if position != 3 {
		t.Fatalf("The seek did not move to the last byte of a non-empty byte: (%d) != (%d)", position, dataLen)
	}
}

func TestSeekableBuffer_Seek_Set_Empty(t *testing.T) {
	sb := NewSeekableBuffer()

	position, err := sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	if position != 0 {
		t.Fatalf("The seek did not move to the first position of a still-empty file: (%d)", position)
	}
}

func TestSeekableBuffer_Seek_Set_Nonempty(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")
	dataLen := len(data)

	n, err := sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written: (%d)", dataLen, n)
	}

	position, err := sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	if position != 0 {
		t.Fatalf("The seek did not move to the first byte of a non-empty file: (%d) != (%d)", position, dataLen)
	}
}

func TestSeekableBuffer_Seek_Cur_Empty(t *testing.T) {
	sb := NewSeekableBuffer()

	position, err := sb.Seek(0, os.SEEK_CUR)
	log.PanicIf(err)

	if position != 0 {
		t.Fatalf("The seek did not stay on the first position of a still-empty file: (%d)", position)
	}
}

func TestSeekableBuffer_Seek_Cur_Nonempty(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")
	dataLen := len(data)

	n, err := sb.Write(data)
	log.PanicIf(err)

	if n != dataLen {
		t.Fatalf("Exactly (%d) bytes weren't written: (%d)", dataLen, n)
	}

	position, err := sb.Seek(0, os.SEEK_CUR)
	log.PanicIf(err)

	if position != 5 {
		t.Fatalf("The seek did not stay on the last byte of a non-empty file: (%d) != (%d)", position, dataLen)
	}
}

func TestSeekableBuffer_Write_SeekAndWritePastEnd(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")

	_, err := sb.Write(data)
	log.PanicIf(err)

	_, err = sb.Seek(10, os.SEEK_SET)
	log.PanicIf(err)

	data2 := []byte("xhellox")

	_, err = sb.Write(data2)
	log.PanicIf(err)

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	buffer := make([]byte, 20)

	n, err := sb.Read(buffer)
	log.PanicIf(err)

	if n != 17 {
		t.Fatalf("Read didn't return the correct number of bytes: (%d)", n)
	}

	expected := []byte("hello\000\000\000\000\000xhellox")

	if bytes.Compare(buffer[:n], expected) != 0 {
		t.Fatalf("Data did not match.")
	}
}

func TestSeekableBuffer_Bytes(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")

	_, err := sb.Write(data)
	log.PanicIf(err)

	_, err = sb.Seek(10, os.SEEK_SET)
	log.PanicIf(err)

	data2 := []byte("xhellox")

	_, err = sb.Write(data2)
	log.PanicIf(err)

	expected := []byte("hello\000\000\000\000\000xhellox")

	if bytes.Compare(sb.Bytes(), expected) != 0 {
		t.Fatalf("Data did not match.")
	}
}

func TestSeekableBuffer_Truncate(t *testing.T) {
	sb := NewSeekableBuffer()

	data := []byte("hello")

	_, err := sb.Write(data)
	log.PanicIf(err)

	err = sb.Truncate(3)
	log.PanicIf(err)

	result := sb.Bytes()
	if bytes.Compare(result, []byte("hel")) != 0 {
		t.Fatalf("Truncated result was not expected: %v", result)
	}
}
