package multistring

import (
	"testing"
)

func TestRemoveEmptyLines_fromSingleString(t *testing.T) {
	input := `			  some test line	    `
	want := `some test line`
	got := RemoveEmptyLines(input)

	if got != want {
		t.Errorf(`got %q, wantend %q`, got, want)
	}
}

func TestRemoveDedupedPackages(t *testing.T) {
	inputMultiline := `
default-browser-id@1.0.4
bplist-parser@0.1.1 deduped
camelcase@2.1.1
map-obj@1.0.1 deduped
`
	want := `
default-browser-id@1.0.4
camelcase@2.1.1
`
	got := RemoveDedupedPackages(inputMultiline)

	if got != want {
		t.Errorf(`got %q, wantend %q`, got, want)
	}
}

// Convert multiline text to go array
func TestMultilinestringToArray(t *testing.T) {
	inputMultiline := `
test1
test2
`
	want := []string{"", "test1", "test2"}
	got := MultilinestringToArray(inputMultiline)

	if len(got) != len(want) {
		t.Errorf(`got %q, wantend %q`, got, want)
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf(`got %q, wantend %q`, got, want)
		}
	}
}
