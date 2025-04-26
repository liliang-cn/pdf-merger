package mergemd

import (
	"fmt"
	"os"
	"path/filepath"
	"pdf-merger/pkg/merger"

	"github.com/spf13/cobra"
)

var (
	inputDir   string
	outputFile string
	addTitles  bool
	verbose    bool
	files      []string // 新增：直接指定文件列表
)

// NewMergeMdCommand 创建merge-md子命令
func NewMergeMdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-md",
		Short: "合并Markdown文件",
		Long:  `合并指定目录下的所有Markdown文件，或合并指定的Markdown文件列表，并按字符顺序排序`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMergeMd()
		},
	}

	// 添加命令行参数
	cmd.Flags().StringVarP(&inputDir, "input", "i", ".", "指定输入目录，包含要合并的Markdown文件")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "merged.md", "指定输出文件名")
	cmd.Flags().BoolVarP(&addTitles, "add-titles", "t", true, "为每个文件添加标题（使用文件名）")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "显示详细信息")
	cmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "指定要合并的Markdown文件列表，如果提供则忽略input参数") // 新增：文件列表参数

	return cmd
}

func runMergeMd() error {
	var result *merger.MergeResult
	var err error

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
		fmt.Printf("输出文件: %s\n", outputFile)
		fmt.Printf("添加标题: %v\n", addTitles)
	}

	// 根据参数选择处理模式：文件列表或目录
	if len(files) > 0 {
		// 使用指定的文件列表
		if verbose {
			fmt.Printf("将合并 %d 个指定的Markdown文件\n", len(files))
		}
		result, err = merger.MergeMarkdownFilesList(files, outputFile, addTitles, verbose)
	} else {
		// 使用目录模式
		// 确保输入目录路径存在且可访问
		inputInfo, err := os.Stat(inputDir)
		if err != nil {
			// 尝试检查是否是路径问题，而不是文件不存在
			if os.IsNotExist(err) {
				fmt.Printf("错误: 输入目录不存在: %s\n", inputDir)
				fmt.Println("注意: 如果是绝对路径，请确保路径完全正确")
			} else {
				fmt.Printf("错误: 无法访问输入目录: %v\n", err)
			}
			return err
		}

		if !inputInfo.IsDir() {
			return fmt.Errorf("%s 不是一个目录", inputDir)
		}

		if verbose {
			fmt.Printf("输入目录: %s\n", inputDir)
		}

		result, err = merger.MergeMarkdownFiles(inputDir, outputFile, addTitles, verbose)
	}

	if err != nil {
		return err
	}

	fmt.Printf("成功! 已将 %d 个Markdown文件合并为: %s\n", result.MergedFiles, result.OutputPath)
	return nil
}
