package juice

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	CLI_NAME  = "juice"
	LOGO_TEXT = `
    _      _        
   (_)_  _(_)__ ___ 
   | | || | / _/ -_)
  _/ |\_,_|_\__\___|
 |__/               
`
)

var (
	verboseVar bool
	inputFile  string
	outputFile string

	rootCmdShortDesc = "Quick and easy tool to help extract licensing information of Node's production packages."
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   CLI_NAME,
	Short: rootCmdShortDesc,
	Long:  fmt.Sprintf("%s\n\n%s", LOGO_TEXT, rootCmdShortDesc),
	Run: func(cmd *cobra.Command, args []string) {
		printHelpInfo()
	},
}

func init() {
	// rootCmd.PersistentFlags().BoolVar(&verboseVar, "verbose", false, "print verbose output")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(versionCmd)
}

// Print information about program and it's usage
func printHelpInfo() {
	fmt.Printf("%s\n", LOGO_TEXT)
	fmt.Printf("Welcome to %s! For more info type '%s help' or use '-h', '-help' flag.\n", CLI_NAME, CLI_NAME)
	os.Exit(0)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
