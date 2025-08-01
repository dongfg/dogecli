package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(bucketListCmd)
}

var bucketListCmd = &cobra.Command{
	Use:   "bucket",
	Short: "List all buckets",
	RunE: func(cmd *cobra.Command, args []string) error {
		buckets, err := cli.BucketList()
		if err != nil {
			return err
		}
		// fmt.Printf("%-15s %-20s %-10s %-8s\n", "名称", "所在区", "大小", "创建时间")
		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"名称", "所在区", "大小", "创建时间"})
		for _, b := range buckets {
			_ = table.Append([]string{b.Name, b.RegionName, b.SpaceHuman, b.CTimeStr})
		}
		_ = table.Render()
		return nil
	},
}
