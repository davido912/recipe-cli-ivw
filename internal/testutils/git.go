package testutils

import (
	"path/filepath"
	"runtime"
)

// helper function to get git project root path

var (
	_, b, _, _ = runtime.Caller(0)

	GitRoot = filepath.Join(filepath.Dir(b), "../..")
)
