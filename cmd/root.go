package cmd

import (
	"os"
	"strings"

	"github.com/cocatrip/fav/pkg/crypto"
	"github.com/cocatrip/fav/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logger.GetLogger()

var rootCmd = &cobra.Command{
	Use:   "fav",
	Short: "",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		flags := cmd.Flags()

		encrypt, err := flags.GetBool("encrypt")
		if err != nil {
			return err
		}

		decrypt, err := flags.GetBool("decrypt")
		if err != nil {
			return err
		}

		if encrypt {
			log.Info("encrypt mode")
			if err := crypto.Encrypt(file); err != nil {
				return err
			}
		} else if decrypt {
			log.Info("decrypt mode")
			if err := crypto.Decrypt(file); err != nil {
				return err
			}
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("encrypt", "e", false, "encrypt a file")
	rootCmd.PersistentFlags().BoolP("decrypt", "d", false, "decrypt a file")
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("fav")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.AddConfigPath(".")
	viper.SetConfigName(".fav")

	if err := viper.ReadInConfig(); err != nil {
		os.Exit(1)
	}
}
