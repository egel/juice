package juice

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display program version details.",
	Run: func(cmd *cobra.Command, args []string) {
		version := getProgramVersion()
		fmt.Printf("%s CLI version: %v\n", CLI_NAME, version)
	},
}

// Extract version of the program
//
// TODO: enhance it with specify version while building program
func getProgramVersion() string {
	return version
}
