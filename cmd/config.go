package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dongfg/dogecli/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

// configCmd writes access and secret keys to config file
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Interactively set AccessKey and SecretKey",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)
		curAccess := viper.GetString("access_key")
		curSecret := viper.GetString("secret_key")

		fmt.Printf("AccessKey [%s]: ", curAccess)
		access, _ := reader.ReadString('\n')
		access = strings.TrimSpace(access)
		if access == "" {
			access = curAccess
		}

		fmt.Printf("SecretKey [%s]: ", curSecret)
		secret, _ := reader.ReadString('\n')
		secret = strings.TrimSpace(secret)
		if secret == "" {
			secret = curSecret
		}

		viper.Set("access_key", access)
		viper.Set("secret_key", secret)

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir := path.Join(home, ".", constants.CLIName)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
		cfgPath := path.Join(dir, "config.yaml")
		if err := viper.WriteConfigAs(cfgPath); err != nil {
			return err
		}
		fmt.Printf("Credentials saved to %s\n", cfgPath)
		return nil
	},
}
