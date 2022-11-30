package main

import (
	"fmt"
	"sort"

	"juice/multistring"
	"juice/npm"
)

func main() {
	// Install packages
	npm.InstallPackageLock()

	// Get list of
	stringsClean := npm.GetListProductionPackagesFromPackageLock()

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

	licenses := npm.FetchPackagesLicences(cleanPackages)

	// print results
	fmt.Println("Printing results")
	for _, p := range licenses {
		fmt.Printf("%+v\n", p)
	}
}
