package cmd

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/database"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/task"
	"os"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start generator task",
	Long: `This command will start sentences generate task. For example:
	$ generator task
It will sync dataset and generate bundle.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkLockFile()
		err := database.Connect()
		if err != nil {
			logging.Logger.Fatal(err.Error())
		}
		task.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func checkLockFile() {
	if _, err := os.Stat("./data/init.lock"); err != nil {
		if os.IsNotExist(err) {
			// File is not exist
			// initSSHCmd.Run(nil, nil)
			cloneCmd.Run(nil, nil)
			if err = os.WriteFile("./data/init.lock", []byte(""), 0700); err != nil {
				logging.Logger.Fatal(err.Error())
			}
		} else {
			// unknown error
			logging.Logger.Fatal(err.Error())
		}
	}
}
