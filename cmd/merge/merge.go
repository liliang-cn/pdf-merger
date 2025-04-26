package merge

import (
	"fmt"
	"path/filepath"
	"pdf-merger/pkg/merger"

	"github.com/spf13/cobra"
)

var (
	inputDir   string
	outputFile string
	verbose    bool
)

// NewMergeCommand 创建merge子命令
func NewMergeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "合并PDF文件",
		Long:  `合并指定目录下的所有PDF文件，并按字符顺序排序`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge()
		},
	}

	// 添加命令行参数
	cmd.Flags().StringVarP(&inputDir, "input", "i", ".", "指定输入目录，包含要合并的PDF文件")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "merged.pdf", "指定输出文件名")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "显示详细信息")

	return cmd
}

func runMerge() error {
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
	}

	// 调用核心逻辑合并PDF
	result, err := merger.MergePDFs(inputDir, outputFile, verbose)
	if err != nil {
		return err
	}

	fmt.Printf("成功! 已将 %d 个PDF文件合并为: %s\n", result.MergedFiles, result.OutputPath)
	return nil
}
