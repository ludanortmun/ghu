package cmd

import (
	"log"

	"github.com/ludanortmun/ghu/internal"
	"github.com/ludanortmun/ghu/internal/gitignore"
	"github.com/spf13/cobra"
)

var gitignoreCmd = &cobra.Command{
	Use:   "gitignore [language] [output file]",
	Short: "Fetches a .gitignore template for the specified programming language from GitHub.",
	Long: `The gitignore command allows you to fetch a .gitignore template for a specified programming language from GitHub.

The first argument must be the name of the programming language for which you want to fetch the .gitignore template (e.g., "Go", "Python", "C++").
It can also be a common alias like "golang", "py" or "cpp".

The second argument is optional and specifies the output file or directory where the .gitignore file will be saved. 
If not provided, it defaults to a file named ".gitignore" in the current directory.
`,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir := ".gitignore"

		if len(args) == 2 {
			outputDir = args[1]
		}

		log.Printf("Fetching .gitignore for %s\n", args[0])

		client := internal.CreateGithubClient()

		command := gitignore.NewGetGitignoreCommand(args[0], outputDir, client)

		err := command.Execute()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(".gitignore for %s saved to %s\n", args[0], outputDir)
	},
	Args: cobra.RangeArgs(1, 2),
}

func init() {
	rootCmd.AddCommand(gitignoreCmd)
}
