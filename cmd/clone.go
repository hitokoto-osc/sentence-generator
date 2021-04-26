package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
	"github.com/pkg/errors"
	"os"

	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone sentences repository",
	Long: `This command is intended to clone sentences remote repositoty. For example:
    $ generator clone
It will clone remote repository to specific path.`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.Logger.Info(fmt.Sprintf("Start Clone repository(%s)...", config.Core.RemoteRepository))
		auth, err := utils.GetGitAuth()
		if err != nil {
			logging.Logger.Fatal(err.Error())
		}
		if _, err := git.PlainClone(config.Core.Workdir, false, &git.CloneOptions{
			URL:      config.Core.RemoteRepository,
			Progress: os.Stdout,
			Auth:     auth,
		}); err != nil {
			logging.Logger.Fatal(errors.WithMessage(err, "can't clone repository successfully").Error())
		}
		logging.Logger.Info("Clone repository successfully.")
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
