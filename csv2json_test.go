package main

import (
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_getFileData(t *testing.T) {

	// defining test slice
	tests := []struct {
		name    string    // nam of the test
		want    inputFile // the input instance to return
		wantErr bool      // whether or not want an error
		osArgs  []string  // command arguments for the test
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

		t.Run(tt.name, func(t *testing.T) {
			// saving original os.Args reference
			actualOsArgs := os.Args

			// defer function to run after the test is done
			defer func() {
				os.Args = actualOsArgs                                           // Restoring the original os.Args
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // resetting the flag command line to parse the flag again
			}()

			os.Args = tt.osArgs       // setting specific command args for this test
			got, err := getFileData() // running the function we want to test

			if (err != nil) != tt.wantErr { // asserting whether or not the correct value id given
				t.Errorf("getFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) { // asserting whether the correct value wanted is given
				t.Errorf("getFileData() = %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_checkIfValidFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfValidFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfValidFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkIfValidFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processCsvFile(t *testing.T) {
	// Defining the maps we're expecting to get from our function
	wantMapSlice := []map[string]string{
		{"COL1": "1", "COL2": "2", "COL3": "3"},
		{"COL1": "4", "COL2": "5", "COL3": "6"},
	}
	// Defining our test cases
	tests := []struct {
		name      string // The name of the test
		csvString string // The content of our tested CSV file
		separator string // The separator used for each test case
	}{
		{"Comma separator", "COL1,COL2,COL3\n1,2,3\n4,5,6\n", "comma"},
		{"Semicolon separator", "COL1;COL2;COL3\n1;2;3\n4;5;6\n", "semicolon"},
	}
	// Iterating our test cases as usual
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating a CSV temp file for testing
			tmpfile, err := ioutil.TempFile("", "test*.csv")
			check(err)
			
			defer os.Remove(tmpfile.Name()) // Removing the CSV test file before living
			_, err = tmpfile.WriteString(tt.csvString) // Writing the content of the CSV test file
			tmpfile.Sync() // Persisting data on disk
			// Defining the inputFile struct that we're going to use as one parameter of our function
			testFileData := inputFile{
				filepath:  tmpfile.Name(),
				pretty:    false,
				separator: tt.separator,
			}
			// Defining the writerChanel
			writerChannel := make(chan map[string]string)
			// Calling the targeted function as a go routine
			go processCsvFile(testFileData, writerChannel)
			// Iterating over the slice containing the expected map values
			for _, wantMap := range wantMapSlice {
				record := <-writerChannel // Waiting for the record that we want to compare
				if !reflect.DeepEqual(record, wantMap) { // Making the corresponding test assertion
					t.Errorf("processCsvFile() = %v, want %v", record, wantMap)
				}
			}
		})
	}
}
