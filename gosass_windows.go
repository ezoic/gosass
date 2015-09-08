package gosass

// +build windows

/*
#cgo windows LDFLAGS: ${SRCDIR}/libsass_windows.a -lstdc++ -lm
#cgo CFLAGS: -Ilibsass
#include "sass.h"
*/
import "C"

import (
	"bytes"
	"os/exec"
	"strings"
)

// return the stubbed Libsass version
func GetLibsassVersion() string {
	return "LIBSASS_WINDOWS"
}

// create the command line for compiling a file
func generateCompileFileCommand(goCtx *FileContext) *exec.Cmd {

	args := make([]string, 0)

	// handle the output style
	switch goCtx.OutputStyle {
	case C.SASS_STYLE_NESTED:
		args = append(args, []string{"--style", "nested"}...)
	case C.SASS_STYLE_EXPANDED:
		args = append(args, []string{"--style", "expanded"}...)
	case C.SASS_STYLE_COMPACT:
		args = append(args, []string{"--style", "compact"}...)
	case C.SASS_STYLE_COMPRESSED:
		args = append(args, []string{"--style", "compressed"}...)
	}

	// include line comments or not
	if goCtx.Options.SourceComments {
		args = append(args, " --line-comments")
	}

	// set the import path
	if len(goCtx.Options.IncludePaths) > 0 {
		args = append(args, []string{"-I", strings.Join(goCtx.Options.IncludePaths, ":")}...)
	}

	args = append(args, goCtx.InputPath)

	cmd := exec.Command("sassc.exe", args...)

	return cmd
}

// create the command line for compiling a string
func generateCompileStringCommand(goCtx *Context) *exec.Cmd {

	args := make([]string, 0)

	// tell command to read from stdin
	args = append(args, "--stdin")

	// handle the output style
	switch goCtx.OutputStyle {
	case C.SASS_STYLE_NESTED:
		args = append(args, []string{"--style", "nested"}...)
	case C.SASS_STYLE_EXPANDED:
		args = append(args, []string{"--style", "expanded"}...)
	case C.SASS_STYLE_COMPACT:
		args = append(args, []string{"--style", "compact"}...)
	case C.SASS_STYLE_COMPRESSED:
		args = append(args, []string{"--style", "compressed"}...)
	}

	// include line comments or not
	if goCtx.Options.SourceComments {
		args = append(args, " --line-comments")
	}

	// set the import path
	if len(goCtx.Options.IncludePaths) > 0 {
		args = append(args, []string{"-I", strings.Join(goCtx.Options.IncludePaths, ":")}...)
	}

	cmd := exec.Command("sassc.exe", args...)

	return cmd
}

// run the sassc command and read the results into an object
func runCommand(cmd *exec.Cmd) (errorStatus int, errorMessage, outputString string) {
	cmdErr := &bytes.Buffer{}
	cmdOut := &bytes.Buffer{}

	cmd.Stderr = cmdErr
	cmd.Stdout = cmdOut

	err := cmd.Run()
	if err != nil {
		errorStatus = 1
		errorMessage = cmdErr.String()
	} else {
		errorStatus = 0
		errorMessage = ""
	}

	outputString = cmdOut.String()

	return errorStatus, errorMessage, outputString
}

// Compile sass string
func Compile(goCtx *Context) {

	cmd := generateCompileStringCommand(goCtx)
	cmd.Stdin = strings.NewReader(goCtx.SourceString)

	goCtx.ErrorStatus, goCtx.ErrorMessage, goCtx.OutputString = runCommand(cmd)

}

// Compile sass from file
func CompileFile(goCtx *FileContext) {

	cmd := generateCompileFileCommand(goCtx)
	cmd.Stdin = nil

	goCtx.ErrorStatus, goCtx.ErrorMessage, goCtx.OutputString = runCommand(cmd)

}
