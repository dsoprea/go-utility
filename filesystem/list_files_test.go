package rifs

import (
	"os"
	"path"
	"sort"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestListFiles_NoPredicate(t *testing.T) {
	filesC, errC := ListFiles(appPath, nil)

	visited := make([]VisitedFile, 0)

FilesRead:

	for {
		select {
		case err, ok := <-errC:
			if ok == true {
				close(filesC)
				close(errC)
			}

			log.PanicIf(err)

		case vf, ok := <-filesC:
			visited = append(visited, vf)

			// TODO(dustin): !! Can vf have a useful value when (ok == false)?
			if ok == false {
				// The goroutine finished.
				break FilesRead
			}
		}
	}

	checkedPathsSs := sort.StringSlice([]string{
		path.Join(appPath, ".git", "objects"),
		path.Join(appPath, "command"),
		path.Join(appPath, "command", "gi_extract_from_images"),
		path.Join(appPath, "utility.go"),
		path.Join(appPath, "utility_test.go"),
	})

	checkedPathsSs.Sort()

	found := 0
	for _, vf := range visited {
		i := checkedPathsSs.Search(vf.Filepath)
		if i < len(checkedPathsSs) && checkedPathsSs[i] == vf.Filepath {
			found++
		}
	}

	if found != len(checkedPathsSs) {
		t.Fatalf("Did not visit all expected paths.")
	}
}

func TestListFiles_WithPredicate(t *testing.T) {
	gitPath := path.Join(appPath, ".git")
	gitObjectsPath := path.Join(gitPath, "objects")

	filter := func(parent string, child os.FileInfo) (bool, error) {
		fullPath := path.Join(parent, child.Name())
		if fullPath == gitPath {
			return true, nil
		} else if fullPath == gitObjectsPath {
			return true, nil
		}

		return false, nil
	}

	filesC, errC := ListFiles(appPath, filter)

	visited := make([]VisitedFile, 0)

FilesRead:

	for {
		select {
		case err, ok := <-errC:
			if ok == true {
				// TODO(dustin): Can we close these on the other side after sending and still get our data?
				close(filesC)
				close(errC)
			}

			log.PanicIf(err)

		case vf, ok := <-filesC:
			// We have finished reading. `vf` has an empty value.
			if ok == false {
				// The goroutine finished.
				break FilesRead
			}

			visited = append(visited, vf)
		}
	}

	if len(visited) != 2 {
		t.Fatalf("We did not visit the count of path we expected: %v", visited)
	} else if visited[0].Filepath != gitPath || visited[1].Filepath != gitObjectsPath {
		t.Fatalf("We did not visit the paths we expected: %v", visited)
	}
}
