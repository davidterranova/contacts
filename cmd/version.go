package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	// Version is populated at compile time
	Version string = "na"
	// BuildTime is populated at compile time
	BuildTime string = "na"
)

// versionCmd prints the version of the program
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the contacts version",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("contacts build %s.%s on %s\n", Version, BuildTime, runtime.Version())
}
