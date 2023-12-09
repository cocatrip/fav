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
	Args:  cobra.MinimumNArgs(1),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".fav.yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Can't read config: %v", err)
		os.Exit(1)
	}
}
