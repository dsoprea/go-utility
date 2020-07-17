package rifs

import (
	"os"
	"testing"
	"time"
)

func TestNewSimpleFileInfoWithFile(t *testing.T) {
	modTime := time.Now()

	sfi := NewSimpleFileInfoWithFile("aa", 11, 22, modTime)
	if sfi.Name() != "aa" {
		t.Fatalf("Filename not correct.")
	} else if sfi.Size() != 11 {
		t.Fatalf("Size not correct.")
	} else if sfi.Mode() != 22 {
		t.Fatalf("Mode not correct.")
	} else if sfi.IsDir() != false {
		t.Fatalf("IsDir() not correct.")
	} else if sfi.ModTime() != modTime {
		t.Fatalf("ModTime() not correct.")
	}
}

func TestNewSimpleFileInfoWithDirectory(t *testing.T) {
	modTime := time.Now()

	sfi := NewSimpleFileInfoWithDirectory("aa", modTime)
	if sfi.Name() != "aa" {
		t.Fatalf("Filename not correct.")
	} else if sfi.Size() != 0 {
		t.Fatalf("Size not correct.")
	} else if sfi.Mode() != os.ModeDir {
		t.Fatalf("Mode not correct.")
	} else if sfi.IsDir() != true {
		t.Fatalf("IsDir() not correct.")
	} else if sfi.ModTime() != modTime {
		t.Fatalf("ModTime() not correct.")
	}
}

func TestSimpleFileInfo_Name(t *testing.T) {
	sfi := &SimpleFileInfo{
		filename: "aa",
	}

	if sfi.Name() != "aa" {
		t.Fatalf("Name not correct.")
	}
}

func TestSimpleFileInfo_Size(t *testing.T) {
	sfi := &SimpleFileInfo{
		size: 123,
	}

	if sfi.Size() != 123 {
		t.Fatalf("Size not correct.")
	}
}

func TestSimpleFileInfo_Mode(t *testing.T) {
	sfi := &SimpleFileInfo{
		mode: 11,
	}

	if sfi.Mode() != 11 {
		t.Fatalf("Mode not correct.")
	}
}

func TestSimpleFileInfo_ModTime(t *testing.T) {
	modTime := time.Now()

	sfi := &SimpleFileInfo{
		modTime: modTime,
	}

	if sfi.ModTime() != modTime {
		t.Fatalf("ModTime not correct.")
	}
}

func TestSimpleFileInfo_IsDir(t *testing.T) {
	sfi := &SimpleFileInfo{
		isDir: true,
	}

	if sfi.IsDir() != true {
		t.Fatalf("IsDir not correct.")
	}
}

func TestSimpleFileInfo_Sys(t *testing.T) {
	sfi := &SimpleFileInfo{}

	if sfi.Sys() != nil {
		t.Fatalf("Sys not correct.")
	}
}
