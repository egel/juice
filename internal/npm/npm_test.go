package npm

import (
	"testing"
)

func TestRemoveAllNpmTreeCharacters(t *testing.T) {
	input := "│ ├─ somePackageName "
	want := "somePackageName"
	got := RemoveAllNpmTreeCharacters(input)

	if got != want {
		t.Errorf(`got %q, wantend %q`, got, want)
	}
}

func TestIsPathExists(t *testing.T) {
	input := "./npm_test.go"
	want := true
	got, _ := IsPathExists(input)

	if got != want {
		t.Errorf(`got %t, wantend %t`, got, want)
	}
}
