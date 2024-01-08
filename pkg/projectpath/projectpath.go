package projectpath

import (
	"path/filepath"
	"runtime"
)

/*
projectpath.go contains the Root variable which is used to get the root folder of the project.

This is used in various points of the application to provide a default folder, and is used
in the tests too
*/

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..")
)
