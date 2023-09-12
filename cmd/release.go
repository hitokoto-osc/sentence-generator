package cmd

import (
	"log"

	"github.com/hitokoto-osc/sentence-generator/logging"
	"github.com/hitokoto-osc/sentence-generator/utils"

	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer logging.Logger.Sync()
		logging.Logger.Info("Commit changes...")
		if err := utils.CommitRepository(); err != nil {
			log.Fatal(err.Error())
		}
		logging.Logger.Info("Create tag changes...")
		if err := utils.ReleaseRepository(); err != nil {
			log.Fatal(err.Error())
		}
		logging.Logger.Info("Push...")
		if err := utils.Push(); err != nil {
			log.Fatal(err.Error())
		}
		logging.Logger.Info("Do release successfully.")
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
