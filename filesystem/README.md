[![Build Status](https://travis-ci.org/RandomIngenuity/go-utility.svg?branch=master)](https://travis-ci.org/RandomIngenuity/go-utility)
[![Coverage Status](https://coveralls.io/repos/github/RandomIngenuity/go-utility/badge.svg?branch=master)](https://coveralls.io/github/RandomIngenuity/go-utility?branch=master)
[![GoDoc](https://godoc.org/github.com/RandomIngenuity/go-utility/filesystem?status.svg)](https://godoc.org/github.com/RandomIngenuity/go-utility/filesystem)

# bounceback_reader

An `io.Reader` that re-seeks to where it is supposed to be before reading. Useful when the file-position is being reused for reads or writes at different positions in the same file resource.

# list_files

A recursive path walker that supports filters.

# seekable_buffer

A memory structure that satisfies `io.ReadWriteSeeker`.
