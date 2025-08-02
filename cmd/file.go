package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fileListCmd)
	rootCmd.AddCommand(fileUploadCmd)
}

const clearLine = "\033[2K"

var fileListCmd = &cobra.Command{
	Use:     "list <bucket[:prefix]>",
	Short:   "List files in bucket",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("缺少参数: bucket")
		}
		bucket, name, err := parse(args[0])
		if err != nil {
			return err
		}
		files, cursor, err := cli.FileList(bucket, name, "")
		if err != nil {
			return err
		}
		// 格式化输出
		fmt.Printf("%-32s  %8s  %s \n", "Hash", "Size", "Name")
		err = keyboard.Open()
		if err != nil {
			return err
		}
		defer func() {
			_ = keyboard.Close()
		}()
		// 第一页标记
		firstPage := true
		// 捕获 Ctrl+C
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		for {
			for _, file := range files {
				fmt.Printf("%-32s  %8s  %s \n", file.Hash, file.FSizeHuman, file.Key)
			}
			if cursor == "" { // 没有更多文件了
				if firstPage {
					fmt.Println("-- 已显示全部，按任意键退出 --")
					_, _, _ = keyboard.GetKey()
				}
				break
			}
			fmt.Print("-- 按空格加载更多，按 q 或 Ctrl+C 退出 --")

			char, key, err := keyboard.GetKey()
			if err != nil {
				log.Fatal(err)
			}

			// 清除提示行（并回到上方）
			fmt.Print("\r" + clearLine)

			if key == keyboard.KeySpace {
				// 清除提示 + 空行
				fmt.Print("\r" + clearLine)
				files, cursor, err = cli.FileList(bucket, name, cursor)
				if err != nil {
					return err
				}
				firstPage = false
			} else if char == 'q' || key == keyboard.KeyCtrlC {
				break
			}

		}
		return nil
	},
}

var fileUploadCmd = &cobra.Command{
	Use:     "copy <localFile> <bucket[:name]>",
	Short:   "Upload file to bucket",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("缺少参数, copy <localFile> <bucket:filename]>")
		}
		localFile := args[0]
		bucket, name, err := parse(args[1])
		if err != nil {
			return err
		}
		if name == "" {
			return fmt.Errorf("缺少参数, filename 不能为空")
		}
		return cli.FileUpload(localFile, bucket, name)
	},
}

func parse(arg string) (bucket, name string, err error) {
	if arg == "" {
		return "", "", fmt.Errorf("缺少参数, bucket 不能为空")
	}
	parts := strings.SplitN(arg, ":", 2)
	bucket = parts[0]

	if bucket == "" {
		return "", "", fmt.Errorf("缺少参数, bucket 不能为空")
	}
	if len(parts) == 2 {
		name = parts[1]
	} else {
		name = ""
	}
	return bucket, name, nil
}
