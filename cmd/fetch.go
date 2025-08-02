package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fileFetchCmd)
	rootCmd.AddCommand(listFetchCmd)
}

var fileFetchCmd = &cobra.Command{
	Use:   "fetch <remoteUrl> <bucket:filename>",
	Short: "Download file to bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("缺少参数, copy <remoteUrl> <bucket:filename>")
		}
		remoteUrl := args[0]
		bucket, name, err := parse(args[1])
		if err != nil {
			return err
		}
		if name == "" {
			return fmt.Errorf("缺少参数, filename 不能为空")
		}
		err = cli.FileFetch(remoteUrl, bucket, name)
		if err != nil {
			return err
		}
		fmt.Println("已提交")
		return nil
	},
}

var listFetchCmd = &cobra.Command{
	Use:     "list-fetch",
	Short:   "Get fetch status",
	Aliases: []string{"ls-fetch"},
	RunE: func(cmd *cobra.Command, args []string) error {
		fetches, err := cli.FileFetchList()
		if err != nil {
			return err
		}
		fmt.Printf("%-20s  %-16s  %s \n", "Time", "State", "Name")
		for _, fetch := range fetches {
			fmt.Printf("%-20s  %-16s  %s \n", fetch.CTime, fetch.State, fetch.Name)
		}
		return nil
	},
}
