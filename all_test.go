package gosass

/*
#cgo linux LDFLAGS: -L. -lsass -lstdc++
#cgo windows LDFLAGS: libsass_windows.a -lstdc++ -lm
*/

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"
)

func runParallel(testFunc func(chan bool), concurrency int) {
	runtime.GOMAXPROCS(4)
	done := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		go testFunc(done)
	}
	for i := 0; i < concurrency; i++ {
		<-done
		<-done
	}
	runtime.GOMAXPROCS(1)
}

const numConcurrentRuns = 200

// const testFileName1     = "test1.scss"
// const testFileName2     = "test2.scss"
// const desiredOutput     = "div {\n  color: black; }\n  div span {\n    color: blue; }\n"

func compileTest(t *testing.T, fileName string) (result string) {

	ctx := FileContext{
		Options: Options{
			OutputStyle:  NESTED_STYLE,
			IncludePaths: make([]string, 0),
		},
		InputPath:    fileName,
		OutputString: "",
		ErrorStatus:  0,
		ErrorMessage: "",
	}

	CompileFile(&ctx)

	if ctx.ErrorStatus != 0 {
		if ctx.ErrorMessage != "" {
			t.Error("ERROR: ", ctx.ErrorMessage)
		} else {
			t.Error("UNKNOWN ERROR")
		}
	} else {
		result = ctx.OutputString
	}

	return result
}

const numTests = 3 // TO DO: read the test dir and set this dynamically

func TestConcurrent(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false
		for i := 1; i <= numTests; i++ {
			inputFile := fmt.Sprintf("test/test%d.scss", i)
			result := compileTest(t, inputFile)
			result = strings.Replace(result, "\r", "", -1)
			desiredOutput, err := ioutil.ReadFile(fmt.Sprintf("test/test%d.css", i))
			if err != nil {
				t.Error(fmt.Sprintf("ERROR: couldn't read test/test%d.css", i))
			}
			desiredOutput = bytes.Replace(desiredOutput, []byte("\r"), []byte(""), -1)
			if result != string(desiredOutput) {
				t.Errorf("ERROR: incorrect output")
			}
		}
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)
}

var testSassFuncs = []struct {
	name     string
	context  Context
	expected string
}{
	{name: "variable with bgcolor",
		context: Context{
			Options: Options{
				OutputStyle:    COMPACT_STYLE,
				SourceComments: false,
				IncludePaths:   make([]string, 0),
			},
			SourceString: `$primaryColor: #eeffcc; body { background: $primaryColor; }`,
		},
		expected: "body { background: #eeffcc; }",
	},
}

func TestCompileFile(t *testing.T) {
	compileTest(t, "test/test1.scss")
}

func TestSassFunctions(t *testing.T) {
	for _, tobj := range testSassFuncs {
		Compile(&tobj.context)

		if tobj.context.ErrorStatus != 0 {
			if tobj.context.ErrorMessage != "" {
				t.Error("ERROR: ", tobj.context.ErrorMessage)
			} else {
				t.Error("UNKNOWN ERROR")
			}
		} else if strings.TrimSpace(tobj.context.OutputString) != strings.TrimSpace(tobj.expected) {
			t.Errorf("Test case %s failed.  Expected \"%s\" but received \"%s\".", tobj.name, tobj.expected, tobj.context.OutputString)
		}
	}
}
