[![GoDoc](https://godoc.org/github.com/dsoprea/go-utility/testing?status.svg)](https://godoc.org/github.com/dsoprea/go-utility/testing)
[![Build Status](https://travis-ci.org/dsoprea/go-utility.svg?branch=master)](https://travis-ci.org/dsoprea/go-utility)
[![Coverage Status](https://coveralls.io/repos/github/dsoprea/go-utility/badge.svg?branch=master)](https://coveralls.io/github/dsoprea/go-utility?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dsoprea/go-utility)](https://goreportcard.com/report/github.com/dsoprea/go-utility)

# redirect_tty

This will temporarily redirect the process TTY resources to support writing
unit-tests directly against `main()` functions.

# handled_exit

Can switch between `os.Exit()` and panicing a return-code. Supports testing
`main()`. Requires calls to `os.Exit()` to call `Exit()` here instead.
