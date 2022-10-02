package main

// all imports required for the test
import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_getFileData(t *testing.T) {

	// defining test slice
	tests := []struct {
		name	string // nam of the test
		want 	inputFile // the input instance to return
		wantErr	bool // whether or not want an error
		osArgs	[]string // command arguments for the test
	}{

		// declaring each unit test input and output data
		{"Default parameters", inputFile{"test.csv", "comma", false}, false, []string{"cmd", "test.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"semicolon enabled", inputFile{"test.csv", "semicolon", false}, false, []string{"cmd", "--separator=semicolon", "test.csv"}},
		{"Pretty enabled", inputFile{"test.csv", "comma", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and semicolon enabled", inputFile{"test.csv", "semicolon", true}, false, []string{"cmd", "--pretty", "--separator=semicolon", "test.csv"}},
		{"Separator not identified", inputFile{}, true, []string{"cmd", "--separator=pipe", "test.csv"}},
	}

	// iterating over the slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T)) {
			// saving original os.Args reference
			actualOsArgs := os.Args

			// defer function to run after the test is done
			defer func() {
				os.Args = actualOsArgs // Restoring the original os.Args
				flag.commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // resetting the flag command line to parse the flag again
			}()

			os.Args = tt.osArgs // setting specific command args for this test
		}
	}
}