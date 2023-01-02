package npm

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

func TestIsNodeModuleExist_returnFalseWhenNodeModulesDirNotExist(t *testing.T) {
	want := false
	got := IsNodeModuleExist() // current path should not exist

	if got != want {
		t.Errorf(`got %t, wantend %t`, got, want)
	}
}

func TestIsPathExists_returnTrueWhenPathExists(t *testing.T) {
	input := "./npm_test.go"
	want := true
	got, _ := IsPathExists(input)

	if got != want {
		t.Errorf(`got %t, wantend %t`, got, want)
	}
}

// fakeExecCommandSuccess is a function that initialises a new exec.Cmd,
// one which will simply call TestExecCommandHelperSuccess rather than the
// command it is provided. It will also pass through the command and its
// arguments as an argument to TestExecCommandHelperSuccess
func fakeExecCommandSuccess(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelperSuccess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)

	mockedStdout := `fake "npm list" command run successfully`
	cmd.Env = []string{
		"GO_TEST_COMMAND_HELPER_PROCESS=1",
		"STDOUT=" + mockedStdout,
		"EXIT_STATUS=0",
	}
	return cmd
}

// TestExecCommandHelperSuccess is a method that is called as a substitute
// for a shell command, the GO_TEST_PROCESS flag ensures that if it is called
// as part of the test suite, it is skipped.
func TestExecCommandHelperSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_COMMAND_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(i)
}

func TestGetListProductionPackagesFromPackageLock(t *testing.T) {
	execCommand_NpmList = fakeExecCommandSuccess

	want := `fake "npm list" command run successfully`
	got := GetListProductionPackagesFromPackageLock()
	if got != want {
		t.Errorf(`got %q, wantend %q`, got, want)
	}
}
