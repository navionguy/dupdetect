package processfile

import (
	"io/ioutil"
	"strings"
	"testing"
)

func Test_loadFile(t *testing.T) {
	inp := "OneCode"

	count, err := loadFile(ioutil.NopCloser(strings.NewReader(inp)))

	if err != nil {
		t.Fatal("loadfile had error on simple load")
	}

	if count != 1 {
		t.Fatal("loadfile failed to count codes")
	}
}

func Test_detectDupe(t *testing.T) {
	inp := "OneCode, AnotherCode, OneCode"

	_, err := loadFile(ioutil.NopCloser(strings.NewReader(inp)))

	if err == nil {
		t.Fatal("duplicate code went unnoticed")
	}
}
