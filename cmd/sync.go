package cmd

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/utils"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync local repository from remote repos",
	Long: `You can use this command to sync repos. For example:
	$ generator sync
It will do git fetch and git reset --hard origin/master command`,
	Run: func(cmd *cobra.Command, args []string) {
		defer logging.Logger.Sync()
		if err := utils.SyncRepository(); err != nil {
			logging.Logger.Fatal(err.Error())
		}
		logging.Logger.Info("Sync repository successfully.")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
