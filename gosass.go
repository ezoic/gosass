package gosass

/*
#cgo CFLAGS: -Ilibsass

#include <stdlib.h>
#include <sass_context.h>
#include <sass_interface.h>
*/
import "C"

type Options struct {
	OutputStyle    int
	SourceComments bool
	IncludePaths   []string
	// eventually gonna' have things like callbacks and whatnot
}

type Context struct {
	Options
	SourceString string
	OutputString string
	ErrorStatus  int
	ErrorMessage string
}

type FileContext struct {
	Options
	InputPath    string
	OutputString string
	ErrorStatus  int
	ErrorMessage string
}

// Constants/enums for the output style.
const (
	NESTED_STYLE = iota
	EXPANDED_STYLE
	COMPACT_STYLE
	COMPRESSED_STYLE
)
