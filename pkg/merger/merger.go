package merger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// PDFFileInfo 存储PDF文件信息
type PDFFileInfo struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

// MergeResult 存储合并结果信息
type MergeResult struct {
	Success      bool     `json:"success"`
	OutputPath   string   `json:"outputPath,omitempty"`
	MergedFiles  int      `json:"mergedFiles,omitempty"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
	FilesList    []string `json:"filesList,omitempty"`
}

// MarkdownFileInfo 存储Markdown文件信息
type MarkdownFileInfo struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

// MergePDFs 合并指定目录下的所有PDF文件
func MergePDFs(inputDir, outputFile string, verbose bool) (*MergeResult, error) {
	// 检查输入目录是否存在
	info, err := os.Stat(inputDir)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("无法访问输入目录 %s: %v", inputDir, err),
		}, err
	}
	if !info.IsDir() {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("%s 不是一个目录", inputDir),
		}, fmt.Errorf("%s 不是一个目录", inputDir)
	}

	// 获取目录中所有的PDF文件
	var pdfFiles []string
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".pdf" {
			pdfFiles = append(pdfFiles, path)
		}
		return nil
	})
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("扫描目录时发生错误: %v", err),
		}, err
	}

	if len(pdfFiles) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("未找到任何PDF文件在目录 %s", inputDir),
		}, fmt.Errorf("未找到任何PDF文件在目录 %s", inputDir)
	}

	// 按字符顺序排序文件
	sort.Strings(pdfFiles)

	if verbose {
		fmt.Printf("找到 %d 个PDF文件，准备合并...\n", len(pdfFiles))
		for i, file := range pdfFiles {
			fmt.Printf("%d: %s\n", i+1, file)
		}
	}

	// 创建配置
	conf := model.NewDefaultConfiguration()

	// 执行合并
	// 设置 dividerPage 为 false，表示不在合并的PDF之间添加分隔页
	err = api.MergeCreateFile(pdfFiles, outputFile, false, conf)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("合并PDF文件失败: %v", err),
		}, err
	}

	return &MergeResult{
		Success:     true,
		OutputPath:  outputFile,
		MergedFiles: len(pdfFiles),
		FilesList:   pdfFiles,
	}, nil
}

// GetPDFFiles 获取指定目录下的所有PDF文件
func GetPDFFiles(inputDir string) ([]PDFFileInfo, error) {
	// 检查输入目录是否存在
	info, err := os.Stat(inputDir)
	if err != nil {
		return nil, fmt.Errorf("无法访问输入目录 %s: %v", inputDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s 不是一个目录", inputDir)
	}

	// 获取目录中所有的PDF文件
	var pdfInfos []PDFFileInfo
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".pdf" {
			title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			pdfInfos = append(pdfInfos, PDFFileInfo{
				Path:  path,
				Title: title,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("扫描目录时发生错误: %v", err)
	}

	// 按字符顺序排序文件
	sort.Slice(pdfInfos, func(i, j int) bool {
		return pdfInfos[i].Path < pdfInfos[j].Path
	})

	return pdfInfos, nil
}

// MergeMarkdownFiles 合并指定目录下的所有Markdown文件
func MergeMarkdownFiles(inputDir, outputFile string, addTitles bool, verbose bool) (*MergeResult, error) {
	// 检查输入目录是否存在
	info, err := os.Stat(inputDir)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("无法访问输入目录 %s: %v", inputDir, err),
		}, err
	}
	if !info.IsDir() {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("%s 不是一个目录", inputDir),
		}, fmt.Errorf("%s 不是一个目录", inputDir)
	}

	// 获取目录中所有的Markdown文件
	var mdFiles []string
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !info.IsDir() && (ext == ".md" || ext == ".markdown") {
			mdFiles = append(mdFiles, path)
		}
		return nil
	})
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("扫描目录时发生错误: %v", err),
		}, err
	}

	if len(mdFiles) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("未找到任何Markdown文件在目录 %s", inputDir),
		}, fmt.Errorf("未找到任何Markdown文件在目录 %s", inputDir)
	}

	// 按字符顺序排序文件
	sort.Strings(mdFiles)

	if verbose {
		fmt.Printf("找到 %d 个Markdown文件，准备合并...\n", len(mdFiles))
		for i, file := range mdFiles {
			fmt.Printf("%d: %s\n", i+1, file)
		}
	}

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("无法创建输出文件: %v", err),
		}, err
	}
	defer outFile.Close()

	// 合并所有Markdown文件
	for i, mdFile := range mdFiles {
		// 读取Markdown文件内容
		content, err := os.ReadFile(mdFile)
		if err != nil {
			return &MergeResult{
				Success:      false,
				ErrorMessage: fmt.Sprintf("读取文件 %s 失败: %v", mdFile, err),
			}, err
		}

		// 如果需要添加标题，则添加文件名作为标题
		if addTitles {
			title := strings.TrimSuffix(filepath.Base(mdFile), filepath.Ext(mdFile))

			// 如果不是第一个文件，先添加分隔符
			if i > 0 {
				outFile.WriteString("\n\n---\n\n")
			}

			// 写入标题
			outFile.WriteString(fmt.Sprintf("# %s\n\n", title))
		} else if i > 0 {
			// 如果不添加标题但不是第一个文件，添加两个换行符作为分隔
			outFile.WriteString("\n\n")
		}

		// 写入文件内容
		outFile.Write(content)
	}

	return &MergeResult{
		Success:     true,
		OutputPath:  outputFile,
		MergedFiles: len(mdFiles),
		FilesList:   mdFiles,
	}, nil
}

// GetMarkdownFiles 获取指定目录下的所有Markdown文件
func GetMarkdownFiles(inputDir string) ([]MarkdownFileInfo, error) {
	// 检查输入目录是否存在
	info, err := os.Stat(inputDir)
	if err != nil {
		return nil, fmt.Errorf("无法访问输入目录 %s: %v", inputDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s 不是一个目录", inputDir)
	}

	// 获取目录中所有的Markdown文件
	var mdInfos []MarkdownFileInfo
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))
		if !info.IsDir() && (ext == ".md" || ext == ".markdown") {
			title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			mdInfos = append(mdInfos, MarkdownFileInfo{
				Path:  path,
				Title: title,
			})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("扫描目录时发生错误: %v", err)
	}

	// 按字符顺序排序文件
	sort.Slice(mdInfos, func(i, j int) bool {
		return mdInfos[i].Path < mdInfos[j].Path
	})

	return mdInfos, nil
}
