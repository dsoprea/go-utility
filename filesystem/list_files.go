package rifs

import (
    "io"
    "os"
    "path"

    "github.com/dsoprea/go-logging"
)

type FileListFilterPredicate func(parent string, child os.FileInfo) (hit bool, err error)

type VisitedFile struct {
    Filepath string
    Info     os.FileInfo
}

// ListFiles feeds a continuous list of files from a recursive folder scan. An
// optional predicate can be provided in order to filter. When done, the
// `filesC` channel is closed. If there's an error, the `errC` channel will
// receive it.
func ListFiles(rootPath string, cb FileListFilterPredicate) (filesC chan VisitedFile, errC chan error) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.Panic(err)
        }
    }()

    // Make sure the path exists.

    f, err := os.Open(rootPath)
    log.PanicIf(err)

    f.Close()

    // Do our thing.

    filesC = make(chan VisitedFile, 100)
    errC = make(chan error, 1)

    go func() {
        defer func() {
            if state := recover(); state != nil {
                err := log.Wrap(state.(error))
                errC <- err
            }
        }()

        queue := []string{rootPath}
        for len(queue) > 0 {
            // Pop the next folder to process off the queue.
            var thisPath string
            thisPath, queue = queue[0], queue[1:]

            // Read information.

            folderF, err := os.Open(thisPath)
            if err != nil {
                errC <- log.Wrap(err)
                return
            }

            // Iterate through children.

            for {
                children, err := folderF.Readdir(1000)
                if err == io.EOF {
                    break
                } else if err != nil {
                    errC <- log.Wrap(err)
                    return
                }

                for _, child := range children {
                    // If a predicate was given, determine if this child will be
                    // left behind.
                    if cb != nil {
                        hit, err := cb(thisPath, child)

                        if err != nil {
                            errC <- log.Wrap(err)
                            return
                        }

                        if hit == false {
                            continue
                        }
                    }

                    // Push file to channel.

                    filepath := path.Join(thisPath, child.Name())

                    vf := VisitedFile{
                        Filepath: filepath,
                        Info:     child,
                    }

                    filesC <- vf

                    // If a folder, queue for later processing.

                    if child.IsDir() == true {
                        queue = append(queue, filepath)
                    }
                }
            }

            folderF.Close()
        }

        close(filesC)
        close(errC)
    }()

    return filesC, errC
}
