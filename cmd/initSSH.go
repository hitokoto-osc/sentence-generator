package cmd

/*
// initSSHCmd represents the initSsh command
var initSSHCmd = &cobra.Command{
	Use:   "initSSH",
	Short: "init ssh public/private key",
	Long: `This command is intended to set SSH public and private key. For example:
	generator initSSH
It will set key according to your config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.Logger.Info("Set private key...")
		if err := os.WriteFile("/root/.ssh/id_rsa", []byte(config.SSH.PrivateKey), 0700); err != nil {
			logging.Logger.Fatal(errors.WithMessage(err, "can't set ssh private key").Error())
		}
		logging.Logger.Info("Set private key successfully.")
		logging.Logger.Info("Set public key...")
		if err := os.WriteFile("/root/.ssh/id_rsa.pub", []byte(config.SSH.PublicKey), 0700); err != nil {
			logging.Logger.Fatal(errors.WithMessage(err, "can't set ssh public key.").Error())
		}
		logging.Logger.Info("Set public key successfully.")
		logging.Logger.Info("Set Github ssh policy...")
		if f, err := os.OpenFile("/root/.ssh/config", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
			logging.Logger.Fatal(err.Error())
		} else {
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					logging.Logger.Fatal(err.Error())
				}
			}(f)
			if _, err = f.WriteString(`Host github.com
StrictHostKeyChecking no`); err != nil {
				logging.Logger.Fatal(errors.WithMessage(err, "can't set  Github ssh policy").Error())
			}
			logging.Logger.Info("Set Github ssh policy successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initSSHCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initSSHCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initSSHCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
*/
