package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghu",
	Short: "CLI tools for interacting with GitHub repositories.",
	Long: `GHU is a CLI library that can be used to interact with GitHub repositories.

It supports features such as:
- Downloading a specific directory from a repository
- Serving websites from GitHub repos
- Supports all public repos and private repositories (requires a PAT)
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
