package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAddGoTags(t *testing.T) {
	for _, sample := range testCases {
		var path string
		tmpfile, err := ioutil.TempFile("", "test")
		if err != nil {
			t.Error(err)
		} else {
			path = tmpfile.Name()
			defer os.Remove(path)
		}
		input, err := os.Open(sample.Input)
		if err != nil {
			t.Error(err)
		}
		bs, err := ioutil.ReadAll(input)
		if err != nil {
			t.Error(err)
		}
		if err := input.Close(); err != nil {
			t.Error(err)
		}
		if _, err := tmpfile.Write(bs); err != nil {
			t.Error(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Error(err)
		}

		err = addGoTags(path)
		if sample.Error {
			if err == nil {
				t.Errorf("TestCase %s execpt some error, but got nil\n",
						 sample.Name)
			}
		} else if err != nil {
			t.Errorf("TestCase %s execpt no error, but got %v\n",
					  sample.Name, err)
		} else {
			file, err := os.Open(path)
			if err != nil {
				t.Error(err)
			}
			output, err := ioutil.ReadAll(file)
			if err != nil {
				t.Error(err)
			}
			if err := file.Close(); err != nil {
				t.Error(err)
			}

			file, err = os.Open(sample.Output)
			if err != nil {
				t.Error(err)
			}
			stand, err := ioutil.ReadAll(file)
			if err != nil {
				t.Error(err)
			}
			if err := file.Close(); err != nil {
				t.Error(err)
			}

			if string(output) != string(stand) {
				t.Errorf("TestCase %s except\n%s\n, but got\n%s\n",
					     sample.Name, string(stand), string(output))
			}
		}
	}
}

var testCases = []struct {
	Name   string
	Input  string
	Output string
	Error  bool
}{
	{"AddNothing", "./testdata/nothing.go", "./testdata/nothing.go", false},
	{"ParseError", "./testdata/invalid.go", "./testdata/invalid.go", true},
	{"AddGoTag", "./testdata/addtag.go", "./testdata/addedtag.go", false},
	{"AddGoTags", "./testdata/addtags.go", "./testdata/addedtags.go", false},
	{"AddGoTagsFromMultiComments", "./testdata/addtagsm.go", "./testdata/addedtagsm.go", false},
	{"AddGoTagsOverride", "./testdata/addtagso.go", "./testdata/addedtagso.go", false},
	{"AddGoTagMatchSyntax", "./testdata/addtagMatch.go", "./testdata/addedtagMatch.go", false},
}
