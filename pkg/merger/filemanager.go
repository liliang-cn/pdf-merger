package merger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// 临时文件目录
const TempDirPrefix = "file-merger-tmp-"

// 文件上传的结果信息结构
type FileUploadResult struct {
	Success      bool   `json:"success"`
	FilePath     string `json:"filePath,omitempty"`
	TempDir      string `json:"tempDir,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	FileSize     int64  `json:"fileSize,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// CreateTempDirectory 创建临时目录
func CreateTempDirectory() (string, error) {
	// 创建带有时间戳的临时目录名称
	timestamp := time.Now().Format("20060102-150405")
	tempDirName := TempDirPrefix + timestamp

	// 在系统临时目录下创建子目录
	tempDir := filepath.Join(os.TempDir(), tempDirName)

	// 创建目录
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	return tempDir, nil
}

// SaveUploadedFile 保存上传的文件到临时目录
func SaveUploadedFile(fileReader io.Reader, fileName string, tempDir string) (*FileUploadResult, error) {
	// 如果未提供临时目录，则创建一个新的
	var err error
	if tempDir == "" {
		tempDir, err = CreateTempDirectory()
		if err != nil {
			return &FileUploadResult{
				Success:      false,
				ErrorMessage: fmt.Sprintf("创建临时目录失败: %v", err),
			}, err
		}
	}

	// 创建新文件
	filePath := filepath.Join(tempDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return &FileUploadResult{
			Success:      false,
			TempDir:      tempDir,
			ErrorMessage: fmt.Sprintf("创建文件失败: %v", err),
		}, err
	}
	defer file.Close()

	// 将上传的内容写入文件
	written, err := io.Copy(file, fileReader)
	if err != nil {
		return &FileUploadResult{
			Success:      false,
			TempDir:      tempDir,
			ErrorMessage: fmt.Sprintf("写入文件失败: %v", err),
		}, err
	}

	return &FileUploadResult{
		Success:  true,
		FilePath: filePath,
		TempDir:  tempDir,
		FileName: fileName,
		FileSize: written,
	}, nil
}

// CleanTempDirectory 清理临时目录
func CleanTempDirectory(tempDir string) error {
	// 检查目录名是否以我们的临时目录前缀开头，以避免删除其他目录
	dirName := filepath.Base(tempDir)
	if len(dirName) < len(TempDirPrefix) || dirName[:len(TempDirPrefix)] != TempDirPrefix {
		return fmt.Errorf("不是有效的临时目录: %s", tempDir)
	}

	// 删除目录及其全部内容
	return os.RemoveAll(tempDir)
}

// ListFilesInTempDir 列出临时目录中的所有文件
func ListFilesInTempDir(tempDir string) ([]string, error) {
	// 检查目录是否存在
	info, err := os.Stat(tempDir)
	if err != nil {
		return nil, fmt.Errorf("访问临时目录失败: %v", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s 不是一个目录", tempDir)
	}

	// 读取目录内容
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, fmt.Errorf("读取目录内容失败: %v", err)
	}

	// 收集文件路径
	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(tempDir, file.Name()))
		}
	}

	return filePaths, nil
}
