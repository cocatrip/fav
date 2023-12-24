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
		log.Infof("file: %v\n", args[0])
		file := args[0]

		log.Infof("secretFile: %v\n", viper.Get("secret-file"))
		secretFile := viper.GetString("secret-file")

		encrypt, err := cmd.Flags().GetBool("encrypt")
		if err != nil {
			return err
		}

		if encrypt {
			log.Infoln("encrypting file using age")
			if err := crypto.Encrypt(secretFile, file); err != nil {
				return err
			}
		}

		decrypt, err := cmd.Flags().GetBool("decrypt")
		if err != nil {
			return err
		}

		if decrypt {
			log.Infoln("decrypting file using age")
			if err := crypto.Decrypt(secretFile, file); err != nil {
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

	rootCmd.PersistentFlags().StringP("secret-file", "f", "", "path to the secret file")
	viper.BindPFlag("secret-file", rootCmd.PersistentFlags().Lookup("secret-file"))
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
