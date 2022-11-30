package main

import (
	"fmt"
	"juice/multistring"
	"juice/npm"
	"log"
	"os/exec"
	"reflect"
	"sort"
)

func main() {
	cmd := exec.Command("npm", "ls", "--production", "--all", "--package-lock-only")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	stringsClean := string(stdout[:])

	// remove deduped lines
	stringsClean = multistring.RemoveDedupedPackages(stringsClean)

	// remove npm tree
	stringsClean = npm.RemoveAllNpmTreeCharacters(stringsClean)

	// clean remaining empty lines
	stringsClean = multistring.RemoveEmptyLines(stringsClean)

	// to array
	packages := multistring.MultilinestringToArray(stringsClean)

	// remove duplicate lines
	cleanPackages := multistring.RemoveDuplicateStr(packages)

	// sort lines
	sort.Strings(packages)

	// count := 5
	var licenses []npm.License

	// create array of channels
	var myChannels []chan npm.License
	for i := 0; i < len(licenses); i++ {
		ch := make(chan npm.License)
		myChannels = append(myChannels, ch)
	}

	cases := make([]reflect.SelectCase, len(myChannels))
	for k, ch := range myChannels {
		cases[k] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		go npm.GetPackageDetails(ch, cleanPackages[k])
	}
	fmt.Println("Waiting for results...")

	// remaining
	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining -= 1
			continue
		}
		fmt.Printf("Read from channel %#v and received %#v\n", myChannels[chosen], value)
	}
}
