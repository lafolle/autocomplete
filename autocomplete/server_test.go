package server

import (
	"testing"
)

func TestAutocomplete(t *testing.T) {

	cases := []struct {
		prefix   string
		expected []string
	}{
		{"zzzzzzskjf", []string{}},

		// In case of exact match the exact word must be returned.
		{"karabagh", []string{"karabagh"}},
		{"bu", []string{"bu", "bual", "buat", "buaze", "bub", "buba", "bubal", "bubale", "bubales", "bubaline"}},
	}

	for i, c := range cases {
		if got, err := Autocomplete(c.prefix); err == nil {
			if len(got) != len(c.expected) {
				t.Errorf("case #%d: got %s expected %s", i, got, c.expected)
			}
		} else {
			t.Errorf("case #%d: got %s expected %s", i, got, c.expected)
		}
	}

}
