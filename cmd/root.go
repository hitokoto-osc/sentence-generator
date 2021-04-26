package cmd

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/spf13/cobra"
)

var cfgFile string
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "A lightweight and swift hitokoto sentences generator tool.",
	Long: `A lightweight and swift hitokoto sentences generator tool. We provide these commands:
	$ generator clone   # Clone remote repository
	$ generator sync    # Sync repository
	$ generator release # Release sentences bundle
	$ generator start   # Start sync task`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.toml)")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "program debug mode")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.InitConfigDriver(cfgFile, logging.Logger)
}

func initLogger() {
	config.SetDebug(debug)
	logging.InitLogger()
}
