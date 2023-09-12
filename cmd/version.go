package cmd

import (
	"fmt"
	"runtime"

	"github.com/hitokoto-osc/sentence-generator/config"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version and build information about hitokoto sentences generator",
	Long: `Show information of this tool about version and build status. For example:
	$ generator version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`A lightweight and swift sentences sync and generator tool. Powered by MoeTeam, a vivid Group.
Version: %s
Build Time: %s
Git Tag: %s
Commit Time: %s
Powered by: %s
`,
			config.Version,
			config.BuildTime,
			config.BuildTag,
			config.CommitTime,
			fmt.Sprintf("%s on %s(%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
