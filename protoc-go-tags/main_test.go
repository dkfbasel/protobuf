package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testDirectory = "./testdata"

func TestProtocGoTags(t *testing.T) {

	// go through all files in the testdata directory
	testfiles, err := ioutil.ReadDir(testDirectory)
	if err != nil {
		t.Error("could not open test directory")
		t.Fail()
	}

	// go through all test files and check example and expected output
	// a parsing error is expected if no expected output file is present
	for _, info := range testfiles {
		if info.IsDir() {
			continue
		}

		// get the name of the file
		exampleFilepath := info.Name()

		// skip all non example files
		if strings.HasSuffix(exampleFilepath, "_example.go") == false {
			continue
		}

		// format the test name
		testName := formatTestName(exampleFilepath)

		// read the example file
		example, err := ioutil.ReadFile(filepath.Join(testDirectory, exampleFilepath))
		if err != nil {
			t.Errorf("could not open sample file: %s, %v\n", exampleFilepath, err)
			continue
		}

		// expect a result when running go tags
		expectParseError := false

		// read the expected golden file
		goldenFilepath := strings.Replace(exampleFilepath, "_example.go", "_expected.go", -1)
		expected, err := ioutil.ReadFile(filepath.Join(testDirectory, goldenFilepath))

		// if no result is given, that we would not expect a result
		if err != nil {
			expectParseError = true
		}

		// generate a temporary file for parsing
		tmpfile, err := ioutil.TempFile(testDirectory, "test")
		if err != nil {
			t.Errorf("could not create temporary file: %s, %v\n", exampleFilepath, err)
			continue
		}

		defer func() {
			err2 := os.Remove(tmpfile.Name())
			if err2 != nil {
				t.Errorf("could not remove temporary file: %s, %v\n", tmpfile.Name(), err)
			}
		}()

		_, err = tmpfile.Write(example)
		if err != nil {
			t.Errorf("could not write to temporary file: %v\n", err)
			continue
		}

		err = tmpfile.Close()
		if err != nil {
			t.Errorf("could not close temporary file: %v\n", err)
			continue
		}

		// run go tags on the temp file
		parseError := applyStructTags(tmpfile.Name())

		if expectParseError == true && parseError != nil {
			t.Logf("Test %s successful\n", testName)
			continue
		}

		if parseError != nil {
			t.Errorf("Test %s faileds: parsing error: %v\n", testName, parseError)
		}

		result, err := ioutil.ReadFile(tmpfile.Name())
		if err != nil {
			t.Errorf("could not read output file: %v\n", err)
			continue
		}

		if bytes.Equal(expected, result) == false {
			t.Errorf("Test %s failed: expected result not matched\n", testName)
			continue
		}

		t.Logf("Test %s successful\n", testName)

	}
}

// formatTestName will format the file name of the example file for printing
func formatTestName(name string) string {
	name = strings.Replace(name, "_example.go", "", -1)
	name = strings.Replace(name, "-", " ", -1)
	return name
}
