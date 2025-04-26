package merger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Temporary file directory prefix
const TempDirPrefix = "file-merger-tmp-"

// File upload result information structure
type FileUploadResult struct {
	Success      bool   `json:"success"`
	FilePath     string `json:"filePath,omitempty"`
	TempDir      string `json:"tempDir,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	FileSize     int64  `json:"fileSize,omitempty"`
	FileType     string `json:"fileType,omitempty"` // File type: pdf or markdown
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// CreateTempDirectory creates a temporary directory
func CreateTempDirectory() (string, error) {
	// Create temporary directory name with timestamp
	timestamp := time.Now().Format("20060102-150405")
	tempDirName := TempDirPrefix + timestamp

	// Create a subdirectory in the system temporary directory
	tempDir := filepath.Join(os.TempDir(), tempDirName)

	// Create the directory
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("Failed to create temporary directory: %v", err)
	}

	return tempDir, nil
}

// SaveUploadedFile saves uploaded file to a temporary directory
func SaveUploadedFile(fileReader io.Reader, fileName string, tempDir string) (*FileUploadResult, error) {
	// If no temporary directory is provided, create a new one
	var err error
	if tempDir == "" {
		tempDir, err = CreateTempDirectory()
		if err != nil {
			return &FileUploadResult{
				Success:      false,
				ErrorMessage: fmt.Sprintf("Failed to create temporary directory: %v", err),
			}, err
		}
	}

	// Create new file
	filePath := filepath.Join(tempDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return &FileUploadResult{
			Success:      false,
			TempDir:      tempDir,
			ErrorMessage: fmt.Sprintf("Failed to create file: %v", err),
		}, err
	}
	defer file.Close()

	// Write uploaded content to file
	written, err := io.Copy(file, fileReader)
	if err != nil {
		return &FileUploadResult{
			Success:      false,
			TempDir:      tempDir,
			ErrorMessage: fmt.Sprintf("Failed to write to file: %v", err),
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

// CleanTempDirectory cleans up temporary directory
func CleanTempDirectory(tempDir string) error {
	// Check if directory name starts with our temporary directory prefix to avoid deleting other directories
	dirName := filepath.Base(tempDir)
	if len(dirName) < len(TempDirPrefix) || dirName[:len(TempDirPrefix)] != TempDirPrefix {
		return fmt.Errorf("Not a valid temporary directory: %s", tempDir)
	}

	// Remove directory and all its contents
	return os.RemoveAll(tempDir)
}

// ListFilesInTempDir lists all files in the temporary directory
func ListFilesInTempDir(tempDir string) ([]string, error) {
	// Check if directory exists
	info, err := os.Stat(tempDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to access temporary directory: %v", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", tempDir)
	}

	// Read directory contents
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to read directory contents: %v", err)
	}

	// Collect file paths
	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(tempDir, file.Name()))
		}
	}

	return filePaths, nil
}
