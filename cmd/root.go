package cmd

import (
	"os"

	"github.com/liliang-cn/pdf-merger/cmd/merge"
	mergemd "github.com/liliang-cn/pdf-merger/cmd/merge-md"
	"github.com/liliang-cn/pdf-merger/cmd/serve"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "file-merger",
		Short: "文件合并工具",
		Long:  `一个合并PDF文件和Markdown文件的命令行工具和API服务器，能够将指定目录下的所有文件合并为一个文件，并按照字符顺序排序。`,
	}

	// 添加子命令
	rootCmd.AddCommand(merge.NewMergeCommand())
	rootCmd.AddCommand(mergemd.NewMergeMdCommand())
	rootCmd.AddCommand(serve.NewServeCommand())
}
