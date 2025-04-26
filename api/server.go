package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/liliang-cn/pdf-merger/pkg/merger"
)

// MergeRequest represents the JSON structure for a PDF merge request
type MergeRequest struct {
	InputDir   string `json:"inputDir"`
	OutputFile string `json:"outputFile"`
}

// MergeMdRequest represents the JSON structure for a Markdown merge request
type MergeMdRequest struct {
	InputDir   string `json:"inputDir"`
	OutputFile string `json:"outputFile"`
	AddTitles  bool   `json:"addTitles"`
}

// TempDirRequest represents the JSON structure for a new temporary directory request
type TempDirRequest struct {
	Purpose string `json:"purpose,omitempty"`
}

// FileUploadInfo represents information about uploaded files
type FileUploadInfo struct {
	TempDir  string   `json:"tempDir"`
	Files    []string `json:"files"`
	FileType string   `json:"fileType"` // 'pdf' or 'markdown'
}

// MergeFilesRequest represents the request structure for merging uploaded files
type MergeFilesRequest struct {
	TempDir    string   `json:"tempDir"`
	FileNames  []string `json:"fileNames,omitempty"` // Optional list of filenames, if empty use all files in directory
	OutputFile string   `json:"outputFile"`
	AddTitles  bool     `json:"addTitles,omitempty"` // Only for Markdown files
}

// StartServer starts the API server
func StartServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("API server started at http://localhost%s\n", addr)
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("  POST /api/merge         - Merge PDF files\n")
	fmt.Printf("  POST /api/merge-md      - Merge Markdown files\n")
	fmt.Printf("  GET  /api/files?dir=... - List PDF files in directory\n")
	fmt.Printf("  GET  /api/md-files?dir=... - List Markdown files in directory\n")
	fmt.Printf("  POST /api/temp-dir      - Create new temporary directory\n")
	fmt.Printf("  POST /api/upload        - Upload files to temporary directory\n")
	fmt.Printf("  GET  /api/temp-files?dir=... - List files in temporary directory\n")
	fmt.Printf("  POST /api/merge-files   - Merge files in temporary directory\n")
	fmt.Printf("  DELETE /api/temp-dir    - Delete temporary directory\n")

	// Register API route handlers
	http.HandleFunc("/api/merge", handleMerge)
	http.HandleFunc("/api/merge-md", handleMergeMd)
	http.HandleFunc("/api/files", handleListFiles)
	http.HandleFunc("/api/md-files", handleListMdFiles)
	http.HandleFunc("/api/download/", handleDownload)
	http.HandleFunc("/api/temp-dir", handleTempDir)
	http.HandleFunc("/api/upload", handleFileUpload)
	http.HandleFunc("/api/temp-files", handleListTempFiles)
	http.HandleFunc("/api/merge-files", handleMergeUploadedFiles)

	return http.ListenAndServe(addr, nil)
}

// handleMerge handles PDF merge requests
func handleMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var req MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.InputDir == "" {
		http.Error(w, "Input directory must be specified", http.StatusBadRequest)
		return
	}

	if req.OutputFile == "" {
		req.OutputFile = "merged.pdf"
	}

	// Ensure output file path is absolute
	if !filepath.IsAbs(req.OutputFile) {
		absPath, err := filepath.Abs(req.OutputFile)
		if err == nil {
			req.OutputFile = absPath
		}
	}

	// Call core logic to merge PDFs
	result, err := merger.MergePDFs(req.InputDir, req.OutputFile, false)
	if err != nil {
		http.Error(w, "Failed to merge PDFs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleListFiles handles requests to list PDF files
func handleListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Get directory parameter
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		http.Error(w, "Directory parameter '?dir=...' must be specified", http.StatusBadRequest)
		return
	}

	// Get PDF files in the directory
	files, err := merger.GetPDFFiles(dir)
	if err != nil {
		http.Error(w, "Failed to get PDF files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// handleDownload provides download for merged PDF files
func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Extract file path from URL path
	filePath := r.URL.Path[len("/api/download/"):]
	if filePath == "" {
		http.Error(w, "File path must be specified", http.StatusBadRequest)
		return
	}

	// Check if file exists
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Cannot access file: "+err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Cannot get file information: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Write file content to response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error transferring file: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleMergeMd handles Markdown merge requests
func handleMergeMd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var req MergeMdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.InputDir == "" {
		http.Error(w, "Input directory must be specified", http.StatusBadRequest)
		return
	}

	if req.OutputFile == "" {
		req.OutputFile = "merged.md"
	}

	// Ensure output file path is absolute
	if !filepath.IsAbs(req.OutputFile) {
		absPath, err := filepath.Abs(req.OutputFile)
		if err == nil {
			req.OutputFile = absPath
		}
	}

	// Call core logic to merge Markdown
	result, err := merger.MergeMarkdownFiles(req.InputDir, req.OutputFile, req.AddTitles, false)
	if err != nil {
		http.Error(w, "Failed to merge Markdown: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleListMdFiles handles requests to list Markdown files
func handleListMdFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Get directory parameter
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		http.Error(w, "Directory parameter '?dir=...' must be specified", http.StatusBadRequest)
		return
	}

	// Get Markdown files in the directory
	files, err := merger.GetMarkdownFiles(dir)
	if err != nil {
		http.Error(w, "Failed to get Markdown files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// handleTempDir handles creation and deletion of temporary directories
func handleTempDir(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Create new temporary directory
		tempDir, err := merger.CreateTempDirectory()
		if err != nil {
			http.Error(w, "Failed to create temporary directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return temporary directory path
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"tempDir": tempDir,
			"message": "Temporary directory created successfully",
		})

	case http.MethodDelete:
		// Delete existing temporary directory
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON request: "+err.Error(), http.StatusBadRequest)
			return
		}

		tempDir := req["tempDir"]
		if tempDir == "" {
			http.Error(w, "Temporary directory path must be specified", http.StatusBadRequest)
			return
		}

		// Delete directory
		if err := merger.CleanTempDirectory(tempDir); err != nil {
			http.Error(w, "Failed to delete temporary directory: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return success information
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Temporary directory successfully deleted",
		})

	default:
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}

// handleFileUpload handles file upload to temporary directory
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Set maximum file size (100MB)
	r.ParseMultipartForm(100 << 20)

	// Get temporary directory
	tempDir := r.FormValue("tempDir")
	if tempDir == "" {
		// If not provided, create a new temporary directory
		var err error
		tempDir, err = merger.CreateTempDirectory()
		if err != nil {
			http.Error(w, "Failed to create temporary directory: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Get uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Verify file type (only PDF and Markdown supported)
	fileName := header.Filename
	fileExt := strings.ToLower(filepath.Ext(fileName))

	var fileType string
	if fileExt == ".pdf" {
		fileType = "pdf"
	} else if fileExt == ".md" || fileExt == ".markdown" {
		fileType = "markdown"
	} else {
		http.Error(w, "Unsupported file type, only PDF or Markdown files are allowed", http.StatusBadRequest)
		return
	}

	// Save file
	result, err := merger.SaveUploadedFile(file, fileName, tempDir)
	if err != nil {
		http.Error(w, "Failed to save uploaded file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add file type information
	result.FileType = fileType

	// Return result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleListTempFiles gets list of files in temporary directory
func handleListTempFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Get temporary directory parameter
	tempDir := r.URL.Query().Get("dir")
	if tempDir == "" {
		http.Error(w, "Directory parameter '?dir=...' must be specified", http.StatusBadRequest)
		return
	}

	// Get files in directory
	files, err := merger.ListFilesInTempDir(tempDir)
	if err != nil {
		http.Error(w, "Failed to get files in temporary directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Distinguish PDF and Markdown files
	var pdfFiles []string
	var mdFiles []string

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		if ext == ".pdf" {
			pdfFiles = append(pdfFiles, file)
		} else if ext == ".md" || ext == ".markdown" {
			mdFiles = append(mdFiles, file)
		}
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tempDir":    tempDir,
		"allFiles":   files,
		"pdfFiles":   pdfFiles,
		"mdFiles":    mdFiles,
		"totalFiles": len(files),
	})
}

// handleMergeUploadedFiles handles merging of uploaded files
func handleMergeUploadedFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var req MergeFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.TempDir == "" {
		http.Error(w, "Temporary directory must be specified", http.StatusBadRequest)
		return
	}

	if req.OutputFile == "" {
		http.Error(w, "Output file name must be specified", http.StatusBadRequest)
		return
	}

	// Ensure output file path is absolute
	if !filepath.IsAbs(req.OutputFile) {
		absPath, err := filepath.Abs(req.OutputFile)
		if err == nil {
			req.OutputFile = absPath
		}
	}

	var result *merger.MergeResult
	var err error

	// Get all files in the temporary directory
	var filesToMerge []string

	if len(req.FileNames) > 0 {
		// Use specified filenames
		for _, fileName := range req.FileNames {
			filesToMerge = append(filesToMerge, filepath.Join(req.TempDir, fileName))
		}
	} else {
		// Use all files in the directory
		filesToMerge, err = merger.ListFilesInTempDir(req.TempDir)
		if err != nil {
			http.Error(w, "Failed to get files in temporary directory: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(filesToMerge) == 0 {
		http.Error(w, "No files found to merge", http.StatusBadRequest)
		return
	}

	// Determine file type based on extension of first file
	fileExt := strings.ToLower(filepath.Ext(filesToMerge[0]))

	if fileExt == ".pdf" {
		// Merge PDF files
		result, err = merger.MergePDFFiles(filesToMerge, req.OutputFile, false)
	} else if fileExt == ".md" || fileExt == ".markdown" {
		// Merge Markdown files
		result, err = merger.MergeMarkdownFilesList(filesToMerge, req.OutputFile, req.AddTitles, false)
	} else {
		http.Error(w, "Unsupported file type, can only merge PDF or Markdown files", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to merge files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
