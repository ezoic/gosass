package gosass

/*
#cgo LDFLAGS: -L. -lsass -lstdc++
#cgo CFLAGS: -Ilibsass

#include <stdlib.h>
#include <sass_context.h>
#include <sass_interface.h>
*/
import "C"
import (
	"strings"
	"unsafe"
)

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

func GetLibsassVersion() string {
	cStr := C.libsass_version()
	return C.GoString(cStr)
}

func Compile(goCtx *Context) {
	// set up the underlying C context struct
	source_string := C.CString(goCtx.SourceString)
	// libsass is deleting this during Block* Context::parse_string() in context.cpp
	//defer C.free(unsafe.Pointer(source_string))

	dataContext := C.sass_make_data_context(source_string)
	defer C.sass_delete_data_context(dataContext)

	options := C.sass_data_context_get_options(dataContext)

	C.sass_option_set_output_style(options, uint32(goCtx.Options.OutputStyle))
	if goCtx.Options.SourceComments {
		C.sass_option_set_source_comments(options, true)
	} else {
		C.sass_option_set_source_comments(options, false)
	}

	include_path := C.CString(strings.Join(goCtx.Options.IncludePaths, ":"))
	// i think that libsass is freeing this when we destroy the context
	defer C.free(unsafe.Pointer(include_path))
	C.sass_option_set_include_path(options, include_path)

	C.sass_data_context_set_options(dataContext, options)
	context := C.sass_data_context_get_context(dataContext)
	compiler := C.sass_make_data_compiler(dataContext)
	defer C.sass_delete_compiler(compiler)

	C.sass_compiler_parse(compiler)
	C.sass_compiler_execute(compiler)

	goCtx.OutputString = C.GoString(C.sass_context_get_output_string(context))
	goCtx.ErrorStatus = int(C.sass_context_get_error_status(context))
	goCtx.ErrorMessage = C.GoString(C.sass_context_get_error_message(context))

}

func CompileFile(goCtx *FileContext) {
	// set up the underlying C context struct
	input_path := C.CString(goCtx.InputPath)
	// libsass is deleting this during Block* Context::parse_string() in context.cpp
	//defer C.free(unsafe.Pointer(input_path))

	fileContext := C.sass_make_file_context(input_path)
	defer C.sass_delete_file_context(fileContext)

	options := C.sass_file_context_get_options(fileContext)

	C.sass_option_set_output_style(options, uint32(goCtx.Options.OutputStyle))
	if goCtx.Options.SourceComments {
		C.sass_option_set_source_comments(options, true)
	} else {
		C.sass_option_set_source_comments(options, false)
	}

	include_path := C.CString(strings.Join(goCtx.Options.IncludePaths, ":"))
	defer C.free(unsafe.Pointer(include_path))
	C.sass_option_set_include_path(options, include_path)

	C.sass_file_context_set_options(fileContext, options)
	context := C.sass_file_context_get_context(fileContext)
	compiler := C.sass_make_file_compiler(fileContext)
	defer C.sass_delete_compiler(compiler)

	C.sass_compiler_parse(compiler)
	C.sass_compiler_execute(compiler)

	goCtx.OutputString = C.GoString(C.sass_context_get_output_string(context))
	goCtx.ErrorStatus = int(C.sass_context_get_error_status(context))
	goCtx.ErrorMessage = C.GoString(C.sass_context_get_error_message(context))

}
