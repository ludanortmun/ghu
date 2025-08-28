package cmd

import (
	"log"

	"github.com/google/go-github/v74/github"
	"github.com/ludanortmun/ghu/internal"
	"github.com/ludanortmun/ghu/internal/clonedir"
	"github.com/spf13/cobra"
)

var cloneDirCmd = &cobra.Command{
	Use:   "clonedir [github url] [output dir]",
	Short: "Downloads a specified directory from a GitHub repository",
	Long: `The clonedir command allows you to download a specific directory from a GitHub repository.

The first argument must be a valid GitHub URL pointing to a repository, and optionally to a specific branch, tag, or commit.
- If the URL does not include a directory path, the entire repository will be downloaded.
- If this URL includes a path to a directory, only that directory and its contents will be downloaded. Note that in this case, the downloaded directory will not be tracked by git.
- If the URL is for a specific file, then only that file will be downloaded.

The second argument is optional and specifies the output directory where the contents will be downloaded. If not provided, it defaults to the current directory.
`,
	Run: func(cmd *cobra.Command, args []string) {
		githubURL := args[0]
		outputDir := "."

		if len(args) == 2 {
			outputDir = args[1]
		}

		target, err := internal.InferTargetFromUrl(githubURL)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Downloading (%s) %s/%s/%s\n", target.Ref, target.Owner, target.Repository, target.Directory)

		client := github.NewClient(nil)
		ghToken, ok := internal.GetAuthToken()

		if ok {
			log.Println("Using authenticated GitHub API client.")
			client = client.WithAuthToken(ghToken)
		} else {
			log.Println("No PAT was found, using unauthenticated GitHub API client. " +
				"If you want to access private repositories, please set a PAT using the `ghws auth set-token` command.")
		}
		downloadCommand := clonedir.NewDownloadCommand(target, outputDir, client)

		err = downloadCommand.Execute()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Successfully downloaded directory.")

	},
	Args: cobra.RangeArgs(1, 2),
}

func init() {
	rootCmd.AddCommand(cloneDirCmd)
}
