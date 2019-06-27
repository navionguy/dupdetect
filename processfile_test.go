package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

var loadTest = []struct {
	codes   string // codes to load
	result  int
	succeed bool
	title   string
}{
	{"OneCode", 3, true, "onecode"},
	{"OneCode, AnotherCode, OneCode", 3, false, "duplicate keys"},
	{"OneCode, AnotherCode, OneMoreCode", 3, true, "multiple keys"},
}

func Test_loadFile(t *testing.T) {
	for _, td := range loadTest {
		cm, err := loadFile(ioutil.NopCloser(strings.NewReader(td.codes)), td.title)

		if len(cm) != td.result {
			t.Errorf("loadFile(%s) returned %d, expected %d\n", td.title, len(cm), td.result)
		}

		if (err != nil) && td.succeed {
			t.Errorf("loadFile(%s) succeeded but was expected to fail", td.title)
		}

		if (err == nil) && !td.succeed {
			t.Errorf("loadFile(%s) failed but was expected to succeed", td.title)
		}
	}
}
