package main

import "testing"

func TestPrefixes(t *testing.T) {

	cases := []struct {
		word  string
		plist []string
	}{
		{"", []string{}},
		{"bam", []string{"b", "ba", "bam", "bam$"}},
	}

	for i, c := range cases {
		if plist := prefixes(c.word); len(plist) != len(c.plist) {
			t.Errorf("case #%d: got %s[%d] expected %s[%d]",
				i, plist, len(plist), c.plist, len(c.plist))
		}
	}

}
