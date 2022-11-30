package main

import (
	"fmt"
	"sort"

	"juice/array"
	"juice/converter"
	"juice/multistring"
	"juice/npm"
)

func main() {
	// Install packages
	npm.InstallPackageLock()

	// Get list of
	fmt.Println("Start gathering the list of production packages...")
	stringsClean := npm.GetListProductionPackagesFromPackageLock()

	fmt.Println("Start data cleaning process...")
	stringsClean = multistring.RemoveDedupedPackages(stringsClean) // remove deduped lines
	stringsClean = npm.RemoveAllNpmTreeCharacters(stringsClean)    // remove npm tree
	stringsClean = multistring.RemoveEmptyLines(stringsClean)      // clean remaining empty lines
	packages := multistring.MultilinestringToArray(stringsClean)   // to array
	cleanPackages := array.RemoveDuplicateStr(packages)            // remove duplicate lines
	sort.Strings(packages)                                         // sort lines

	fmt.Println("Fetching related license texts...")
	licenses := npm.FetchPackagesLicences(cleanPackages)

	// print results
	fmt.Println("Preparing final result...")
	converter.SaveDataToCSVFile(licenses)
}
