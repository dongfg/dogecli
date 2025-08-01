package config

import (
	"fmt"
	"github.com/dongfg/dogecli/internal/constants"
	"github.com/spf13/viper"
	"os"
	"path"
)

// Init loads configuration from file and environment
func Init(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configDir := path.Join(home, ".", constants.CLIName)
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("DOGECLI")
	_ = viper.BindEnv(constants.EnvAccessKey)
	_ = viper.BindEnv(constants.EnvSecretKey)

	if err := viper.ReadInConfig(); err != nil {
		// Ignoring if config not found
	}
}
