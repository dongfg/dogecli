package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bucketListCmd)
}

var bucketListCmd = &cobra.Command{
	Use:     "list-bucket",
	Short:   "List all buckets",
	Aliases: []string{"ls-bucket", "lsb"},
	RunE: func(cmd *cobra.Command, args []string) error {
		buckets, err := cli.BucketList()
		if err != nil {
			return err
		}
		// fmt.Printf("%-15s %-20s %-10s %-8s\n", "名称", "所在区", "大小", "创建时间")
		table := tablewriter.NewTable(os.Stdout)
		table.Options(tablewriter.WithRendition(tw.Rendition{Borders: tw.BorderNone}))
		table.Header([]string{"名称", "所在区", "大小", "创建时间"})
		for _, b := range buckets {
			_ = table.Append([]string{b.Name, b.RegionName, b.SpaceHuman, b.CTimeStr})
		}
		_ = table.Render()
		return nil
	},
}
