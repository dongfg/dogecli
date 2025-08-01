package cmd

import (
	"fmt"
	"github.com/dongfg/dogecli/internal/client"
	"github.com/dongfg/dogecli/internal/config"
	"github.com/dongfg/dogecli/internal/constants"
	"github.com/dongfg/dogecli/internal/logger"
	"os"

	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
)

var (
	cfgFile string
	cli     *client.Client
	rootCmd = &cobra.Command{
		Use:   constants.CLIName,
		Short: `多吉云基础型云储存管理工具`,
		Long:  `多吉云基础型云储存管理工具`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.Init(cfgFile)
			logger.Init()
			cli = client.New()
		},
	}
)

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		fmt.Sprintf("config file (default is $HOME/.%s/config.yaml)", constants.CLIName))
}
