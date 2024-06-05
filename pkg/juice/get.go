package juice

import (
	"fmt"
	"sort"

	"github.com/egel/juice/internal/converter"
	"github.com/egel/juice/internal/multistring"
	"github.com/egel/juice/internal/npm"
	"github.com/egel/juice/pkg/array"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Command to extract licensing information of Node's production packages and pull information from npm registry.",
	Run: func(cmd *cobra.Command, args []string) {
		FetchAndExport()
	},
}

func FetchAndExport() {
	// Check files or die
	pjson, err := npm.IsPathExists(npm.PackageJsonFilePath)
	if err != nil {
		log.Fatal().Err(err)
	}
	if !pjson {
		log.Fatal().
			Str("PackageJsonFilePath", npm.PackageJsonFilePath).
			Msg("package.json file path does not exist")
	}
	pljson, err := npm.IsPathExists(npm.PackageLockJsonFilePath)
	if err != nil {
		log.Fatal().Err(err)
	}
	if !pljson {
		log.Fatal().
			Str("PackageLockJsonFilePath", npm.PackageJsonFilePath).
			Msg("package-lock.json file path does not exist")
	}

	npm.InstallPackageLock()
	nodeModules := npm.IsNodeModuleExist()
	if !nodeModules {
		log.Fatal().Msg("node_module directory not exist. Install packages first.")
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
	fmt.Printf("\nDone\n")
}
