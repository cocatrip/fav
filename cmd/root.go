package cmd

import (
	"os"

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
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("encrypt", "e", false, "encrypt a file")
	rootCmd.PersistentFlags().BoolP("decrypt", "d", false, "decrypt a file")

	rootCmd.PersistentFlags().StringP("secret-file", "f", "", "path to the secret file")

	rootCmd.PersistentFlags().StringP("secret-key", "k", "", "secret key to use")
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".fav.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Can't read config: %v", err)
		os.Exit(1)
	}
}
