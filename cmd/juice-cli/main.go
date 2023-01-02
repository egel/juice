package main

import (
	"fmt"
	"log"
	"sort"

	"juice/internal/array"
	"juice/internal/converter"
	"juice/internal/multistring"
	"juice/internal/npm"
)

func main() {
	// Install packages or die
	if _, err := npm.IsPathExists(npm.PackageJsonFilePath); err != nil {
		log.Fatalf("%s does not exist", npm.PackageJsonFilePath)
	}
	if _, err := npm.IsPathExists(npm.PackageLockJsonFilePath); err != nil {
		log.Fatalf("%s does not exist", npm.PackageLockJsonFilePath)
	}

	npm.InstallPackageLock()
	nodeModules := npm.IsNodeModuleExist()
	if nodeModules != true {
		log.Fatalf("node_module directory not exist. Install packages first.")
	}

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

	npm.PrintNumberOfPackages(packages)

	fmt.Println("Start fetching licenses...")
	licenses := npm.FetchPackagesLicences(cleanPackages)

	// print results
	fmt.Printf("\nPreparing final result...")
	converter.SaveDataToCSVFile(licenses)
	fmt.Printf("\nDone")
}
