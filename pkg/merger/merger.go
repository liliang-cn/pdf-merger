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

// PDFFileInfo stores PDF file information
type PDFFileInfo struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

// MergeResult stores merge operation result information
type MergeResult struct {
	Success      bool     `json:"success"`
	OutputPath   string   `json:"outputPath,omitempty"`
	MergedFiles  int      `json:"mergedFiles,omitempty"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
	FilesList    []string `json:"filesList,omitempty"`
}

// MarkdownFileInfo stores Markdown file information
type MarkdownFileInfo struct {
	Path  string `json:"path"`
	Title string `json:"title"`
}

// MergePDFs merges all PDF files in the specified directory
func MergePDFs(inputDir, outputFile string, verbose bool) (*MergeResult, error) {
	// Check if input directory exists
	info, err := os.Stat(inputDir)
	if err != nil {
		// Add error handling for absolute paths
		errMsg := fmt.Sprintf("Cannot access input directory %s: %v", inputDir, err)
		if os.IsNotExist(err) {
			if filepath.IsAbs(inputDir) {
				errMsg = fmt.Sprintf("Specified absolute path directory does not exist: %s", inputDir)
			} else {
				errMsg = fmt.Sprintf("Specified directory does not exist: %s", inputDir)
			}
		}
		return &MergeResult{
			Success:      false,
			ErrorMessage: errMsg,
		}, err
	}

	if !info.IsDir() {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("%s is not a directory", inputDir),
		}, fmt.Errorf("%s is not a directory", inputDir)
	}

	// Get all PDF files in the directory
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
			ErrorMessage: fmt.Sprintf("Error scanning directory: %v", err),
		}, err
	}

	if len(pdfFiles) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("No PDF files found in directory %s", inputDir),
		}, fmt.Errorf("No PDF files found in directory %s", inputDir)
	}

	// Sort files in alphanumeric order
	sort.Strings(pdfFiles)

	if verbose {
		fmt.Printf("Found %d PDF files, preparing to merge...\n", len(pdfFiles))
		for i, file := range pdfFiles {
			fmt.Printf("%d: %s\n", i+1, file)
		}
	}

	// Create configuration
	conf := model.NewDefaultConfiguration()

	// Execute merge
	// Set dividerPage to false, meaning don't add separator pages between merged PDFs
	err = api.MergeCreateFile(pdfFiles, outputFile, false, conf)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to merge PDF files: %v", err),
		}, err
	}

	return &MergeResult{
		Success:     true,
		OutputPath:  outputFile,
		MergedFiles: len(pdfFiles),
		FilesList:   pdfFiles,
	}, nil
}

// GetPDFFiles gets all PDF files in the specified directory
func GetPDFFiles(inputDir string) ([]PDFFileInfo, error) {
	// Check if input directory exists
	info, err := os.Stat(inputDir)
	if err != nil {
		return nil, fmt.Errorf("Cannot access input directory %s: %v", inputDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", inputDir)
	}

	// Get all PDF files in the directory
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
		return nil, fmt.Errorf("Error scanning directory: %v", err)
	}

	// Sort files in alphanumeric order
	sort.Slice(pdfInfos, func(i, j int) bool {
		return pdfInfos[i].Path < pdfInfos[j].Path
	})

	return pdfInfos, nil
}

// MergeMarkdownFiles merges all Markdown files in the specified directory
func MergeMarkdownFiles(inputDir, outputFile string, addTitles bool, verbose bool) (*MergeResult, error) {
	// Check if input directory exists
	info, err := os.Stat(inputDir)
	if err != nil {
		// Add error handling for absolute paths
		errMsg := fmt.Sprintf("Cannot access input directory %s: %v", inputDir, err)
		if os.IsNotExist(err) {
			if filepath.IsAbs(inputDir) {
				errMsg = fmt.Sprintf("Specified absolute path directory does not exist: %s", inputDir)
			} else {
				errMsg = fmt.Sprintf("Specified directory does not exist: %s", inputDir)
			}
		}
		return &MergeResult{
			Success:      false,
			ErrorMessage: errMsg,
		}, err
	}

	if !info.IsDir() {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("%s is not a directory", inputDir),
		}, fmt.Errorf("%s is not a directory", inputDir)
	}

	// Get all Markdown files in the directory
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
			ErrorMessage: fmt.Sprintf("Error scanning directory: %v", err),
		}, err
	}

	if len(mdFiles) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("No Markdown files found in directory %s", inputDir),
		}, fmt.Errorf("No Markdown files found in directory %s", inputDir)
	}

	// Sort files in alphanumeric order
	sort.Strings(mdFiles)

	if verbose {
		fmt.Printf("Found %d Markdown files, preparing to merge...\n", len(mdFiles))
		for i, file := range mdFiles {
			fmt.Printf("%d: %s\n", i+1, file)
		}
	}

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Cannot create output file: %v", err),
		}, err
	}
	defer outFile.Close()

	// Merge all Markdown files
	for i, mdFile := range mdFiles {
		// Read Markdown file content
		content, err := os.ReadFile(mdFile)
		if err != nil {
			return &MergeResult{
				Success:      false,
				ErrorMessage: fmt.Sprintf("Failed to read file %s: %v", mdFile, err),
			}, err
		}

		// If titles should be added, add filename as title
		if addTitles {
			title := strings.TrimSuffix(filepath.Base(mdFile), filepath.Ext(mdFile))

			// If not the first file, add separator first
			if i > 0 {
				outFile.WriteString("\n\n---\n\n")
			}

			// Write title
			outFile.WriteString(fmt.Sprintf("# %s\n\n", title))
		} else if i > 0 {
			// If not adding titles but not the first file, add two newlines as separator
			outFile.WriteString("\n\n")
		}

		// Write file content
		outFile.Write(content)
	}

	return &MergeResult{
		Success:     true,
		OutputPath:  outputFile,
		MergedFiles: len(mdFiles),
		FilesList:   mdFiles,
	}, nil
}

// GetMarkdownFiles gets all Markdown files in the specified directory
func GetMarkdownFiles(inputDir string) ([]MarkdownFileInfo, error) {
	// Check if input directory exists
	info, err := os.Stat(inputDir)
	if err != nil {
		return nil, fmt.Errorf("Cannot access input directory %s: %v", inputDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", inputDir)
	}

	// Get all Markdown files in the directory
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
		return nil, fmt.Errorf("Error scanning directory: %v", err)
	}

	// Sort files in alphanumeric order
	sort.Slice(mdInfos, func(i, j int) bool {
		return mdInfos[i].Path < mdInfos[j].Path
	})

	return mdInfos, nil
}

// MergePDFFiles merges the specified list of PDF files
func MergePDFFiles(files []string, outputFile string, verbose bool) (*MergeResult, error) {
	if len(files) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: "No files provided",
		}, fmt.Errorf("No files provided")
	}

	// Validate that each file exists and is a PDF
	validFiles := make([]string, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			if verbose {
				fmt.Printf("Warning: Cannot access file %s: %v, skipped\n", file, err)
			}
			continue
		}

		if info.IsDir() {
			if verbose {
				fmt.Printf("Warning: %s is a directory, not a file, skipped\n", file)
			}
			continue
		}

		if strings.ToLower(filepath.Ext(file)) != ".pdf" {
			if verbose {
				fmt.Printf("Warning: %s is not a PDF file, skipped\n", file)
			}
			continue
		}

		validFiles = append(validFiles, file)
	}

	if len(validFiles) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: "No valid PDF files to merge",
		}, fmt.Errorf("No valid PDF files to merge")
	}

	if verbose {
		fmt.Printf("Found %d valid PDF files, preparing to merge...\n", len(validFiles))
		for i, file := range validFiles {
			fmt.Printf("%d: %s\n", i+1, file)
		}
	}

	// Create configuration
	conf := model.NewDefaultConfiguration()

	// Execute merge
	err := api.MergeCreateFile(validFiles, outputFile, false, conf)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Failed to merge PDF files: %v", err),
		}, err
	}

	return &MergeResult{
		Success:     true,
		OutputPath:  outputFile,
		MergedFiles: len(validFiles),
		FilesList:   validFiles,
	}, nil
}

// MergeMarkdownFilesList merges the specified list of Markdown files
func MergeMarkdownFilesList(files []string, outputFile string, addTitles bool, verbose bool) (*MergeResult, error) {
	if len(files) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: "No files provided",
		}, fmt.Errorf("No files provided")
	}

	// Validate that each file exists and is a Markdown file
	validFiles := make([]string, 0, len(files))
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			if verbose {
				fmt.Printf("Warning: Cannot access file %s: %v, skipped\n", file, err)
			}
			continue
		}

		if info.IsDir() {
			if verbose {
				fmt.Printf("Warning: %s is a directory, not a file, skipped\n", file)
			}
			continue
		}

		ext := strings.ToLower(filepath.Ext(file))
		if ext != ".md" && ext != ".markdown" {
			if verbose {
				fmt.Printf("Warning: %s is not a Markdown file, skipped\n", file)
			}
			continue
		}

		validFiles = append(validFiles, file)
	}

	if len(validFiles) == 0 {
		return &MergeResult{
			Success:      false,
			ErrorMessage: "No valid Markdown files to merge",
		}, fmt.Errorf("No valid Markdown files to merge")
	}

	if verbose {
		fmt.Printf("Found %d valid Markdown files, preparing to merge...\n", len(validFiles))
		for i, file := range validFiles {
			fmt.Printf("%d: %s\n", i+1, file)
		}
	}

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return &MergeResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("Cannot create output file: %v", err),
		}, err
	}
	defer outFile.Close()

	// Merge all Markdown files
	for i, mdFile := range validFiles {
		// Read Markdown file content
		content, err := os.ReadFile(mdFile)
		if err != nil {
			return &MergeResult{
				Success:      false,
				ErrorMessage: fmt.Sprintf("Failed to read file %s: %v", mdFile, err),
			}, err
		}

		// If titles should be added, add filename as title
		if addTitles {
			title := strings.TrimSuffix(filepath.Base(mdFile), filepath.Ext(mdFile))

			// If not the first file, add separator first
			if i > 0 {
				outFile.WriteString("\n\n---\n\n")
			}

			// Write title
			outFile.WriteString(fmt.Sprintf("# %s\n\n", title))
		} else if i > 0 {
			// If not adding titles but not the first file, add two newlines as separator
			outFile.WriteString("\n\n")
		}

		// Write file content
		outFile.Write(content)
	}

	return &MergeResult{
		Success:     true,
		OutputPath:  outputFile,
		MergedFiles: len(validFiles),
		FilesList:   validFiles,
	}, nil
}
