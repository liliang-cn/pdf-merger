package mergemd

import (
	"fmt"
	"path/filepath"
	"pdf-merger/pkg/merger"

	"github.com/spf13/cobra"
)

var (
	inputDir   string
	outputFile string
	addTitles  bool
	verbose    bool
)

// NewMergeMdCommand 创建merge-md子命令
func NewMergeMdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-md",
		Short: "合并Markdown文件",
		Long:  `合并指定目录下的所有Markdown文件，并按字符顺序排序`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMergeMd()
		},
	}

	// 添加命令行参数
	cmd.Flags().StringVarP(&inputDir, "input", "i", ".", "指定输入目录，包含要合并的Markdown文件")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "merged.md", "指定输出文件名")
	cmd.Flags().BoolVarP(&addTitles, "add-titles", "t", true, "为每个文件添加标题（使用文件名）")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "显示详细信息")

	return cmd
}

func runMergeMd() error {
	// 确保输出文件路径是绝对路径
	if !filepath.IsAbs(outputFile) {
		absPath, err := filepath.Abs(outputFile)
		if err != nil {
			fmt.Printf("警告: 无法获取绝对路径: %v, 将使用相对路径\n", err)
		} else {
			outputFile = absPath
		}
	}

	if verbose {
		fmt.Printf("输入目录: %s\n", inputDir)
		fmt.Printf("输出文件: %s\n", outputFile)
		fmt.Printf("添加标题: %v\n", addTitles)
	}

	// 调用核心逻辑合并Markdown
	result, err := merger.MergeMarkdownFiles(inputDir, outputFile, addTitles, verbose)
	if err != nil {
		return err
	}

	fmt.Printf("成功! 已将 %d 个Markdown文件合并为: %s\n", result.MergedFiles, result.OutputPath)
	return nil
}
