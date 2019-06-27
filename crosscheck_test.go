package main

import (
	"testing"
)

type testrun struct {
	set1   []string // first set of codes
	set2   []string // second set of codes
	result int      // what I expect in the
	title  string   // name of the test run
}

var checkTest = []testrun{
	{[]string{"code1", "code2", "code3"}, []string{"code4", "code5", "code6"}, 3, "no duplicates"},
	{[]string{"code1", "code2", "code3"}, []string{"code4", "code5", "code6", "code1"}, -1, "duplicate keys"},
}

func Test_checkFriend(t *testing.T) {

	for _, td := range checkTest {
		s1 := make(codeMap)
		s2 := make(codeMap)

		for _, cd := range td.set1 {
			s1[cd] = true
		}
		for _, cd := range td.set2 {
			s2[cd] = true
		}

		w := worker{
			fname: "Test_checkFriend",
			outp:  make(chan results),
		}

		chk := request{
			fname: "testRequest",
			codes: s2,
		}
		go checkFriend(&w, chk, s1)

		res := <-w.outp

		if res.codesSeen != td.result {
			t.Errorf("Test_checkFriend returned %d, was expecting %d\n", res.codesSeen, td.result)
		}
	}
}
