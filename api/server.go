package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pdf-merger/pkg/merger"
	"strconv"
	"strings"
)

// MergeRequest 表示合并PDF请求的JSON结构
type MergeRequest struct {
	InputDir   string `json:"inputDir"`
	OutputFile string `json:"outputFile"`
}

// MergeMdRequest 表示合并Markdown请求的JSON结构
type MergeMdRequest struct {
	InputDir   string `json:"inputDir"`
	OutputFile string `json:"outputFile"`
	AddTitles  bool   `json:"addTitles"`
}

// TempDirRequest 表示请求新临时目录的JSON结构
type TempDirRequest struct {
	Purpose string `json:"purpose,omitempty"`
}

// FileUploadInfo 表示上传的文件信息
type FileUploadInfo struct {
	TempDir  string   `json:"tempDir"`
	Files    []string `json:"files"`
	FileType string   `json:"fileType"` // 'pdf' 或 'markdown'
}

// MergeFilesRequest 表示合并上传文件的请求结构
type MergeFilesRequest struct {
	TempDir    string   `json:"tempDir"`
	FileNames  []string `json:"fileNames,omitempty"` // 可选的文件名列表，若为空则使用目录中的所有文件
	OutputFile string   `json:"outputFile"`
	AddTitles  bool     `json:"addTitles,omitempty"` // 仅用于Markdown文件
}

// StartServer 启动API服务器
func StartServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("API服务器启动在 http://localhost%s\n", addr)
	fmt.Printf("可用端点:\n")
	fmt.Printf("  POST /api/merge         - 合并PDF文件\n")
	fmt.Printf("  POST /api/merge-md      - 合并Markdown文件\n")
	fmt.Printf("  GET  /api/files?dir=... - 列出目录中的PDF文件\n")
	fmt.Printf("  GET  /api/md-files?dir=... - 列出目录中的Markdown文件\n")
	fmt.Printf("  POST /api/temp-dir      - 创建新的临时目录\n")
	fmt.Printf("  POST /api/upload        - 上传文件到临时目录\n")
	fmt.Printf("  GET  /api/temp-files?dir=... - 列出临时目录中的文件\n")
	fmt.Printf("  POST /api/merge-files   - 合并临时目录中的文件\n")
	fmt.Printf("  DELETE /api/temp-dir    - 删除临时目录\n")

	// 注册API路由处理程序
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

// handleMerge 处理PDF合并请求
func handleMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的JSON请求: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.InputDir == "" {
		http.Error(w, "必须指定输入目录", http.StatusBadRequest)
		return
	}

	if req.OutputFile == "" {
		req.OutputFile = "merged.pdf"
	}

	// 确保输出文件路径是绝对路径
	if !filepath.IsAbs(req.OutputFile) {
		absPath, err := filepath.Abs(req.OutputFile)
		if err == nil {
			req.OutputFile = absPath
		}
	}

	// 调用核心逻辑合并PDF
	result, err := merger.MergePDFs(req.InputDir, req.OutputFile, false)
	if err != nil {
		http.Error(w, "合并PDF失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleListFiles 处理列出PDF文件请求
func handleListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取目录参数
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		http.Error(w, "必须指定目录参数 '?dir=...'", http.StatusBadRequest)
		return
	}

	// 获取目录中的PDF文件
	files, err := merger.GetPDFFiles(dir)
	if err != nil {
		http.Error(w, "获取PDF文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// handleDownload 提供已合并的PDF文件下载
func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中提取文件路径
	filePath := r.URL.Path[len("/api/download/"):]
	if filePath == "" {
		http.Error(w, "必须指定文件路径", http.StatusBadRequest)
		return
	}

	// 检查文件是否存在
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "无法访问文件: "+err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "无法获取文件信息: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// 将文件内容写入响应
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "传输文件时出错: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleMergeMd 处理Markdown合并请求
func handleMergeMd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req MergeMdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的JSON请求: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.InputDir == "" {
		http.Error(w, "必须指定输入目录", http.StatusBadRequest)
		return
	}

	if req.OutputFile == "" {
		req.OutputFile = "merged.md"
	}

	// 确保输出文件路径是绝对路径
	if !filepath.IsAbs(req.OutputFile) {
		absPath, err := filepath.Abs(req.OutputFile)
		if err == nil {
			req.OutputFile = absPath
		}
	}

	// 调用核心逻辑合并Markdown
	result, err := merger.MergeMarkdownFiles(req.InputDir, req.OutputFile, req.AddTitles, false)
	if err != nil {
		http.Error(w, "合并Markdown失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleListMdFiles 处理列出Markdown文件请求
func handleListMdFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取目录参数
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		http.Error(w, "必须指定目录参数 '?dir=...'", http.StatusBadRequest)
		return
	}

	// 获取目录中的Markdown文件
	files, err := merger.GetMarkdownFiles(dir)
	if err != nil {
		http.Error(w, "获取Markdown文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// handleTempDir 处理临时目录的创建和删除
func handleTempDir(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// 创建新的临时目录
		tempDir, err := merger.CreateTempDirectory()
		if err != nil {
			http.Error(w, "创建临时目录失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回临时目录路径
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"tempDir": tempDir,
			"message": "临时目录创建成功",
		})

	case http.MethodDelete:
		// 删除现有临时目录
		var req map[string]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "无效的JSON请求: "+err.Error(), http.StatusBadRequest)
			return
		}

		tempDir := req["tempDir"]
		if tempDir == "" {
			http.Error(w, "必须指定临时目录路径", http.StatusBadRequest)
			return
		}

		// 删除目录
		if err := merger.CleanTempDirectory(tempDir); err != nil {
			http.Error(w, "删除临时目录失败: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回成功信息
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "临时目录已成功删除",
		})

	default:
		http.Error(w, "不支持的HTTP方法", http.StatusMethodNotAllowed)
	}
}

// handleFileUpload 处理文件上传到临时目录
func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	// 设置最大文件大小（100MB）
	r.ParseMultipartForm(100 << 20)

	// 获取临时目录
	tempDir := r.FormValue("tempDir")
	if tempDir == "" {
		// 如果未提供，则创建一个新的临时目录
		var err error
		tempDir, err = merger.CreateTempDirectory()
		if err != nil {
			http.Error(w, "创建临时目录失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 获取上传的文件
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "获取上传的文件失败: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 验证文件类型（仅支持PDF和Markdown）
	fileName := header.Filename
	fileExt := strings.ToLower(filepath.Ext(fileName))

	var fileType string
	if fileExt == ".pdf" {
		fileType = "pdf"
	} else if fileExt == ".md" || fileExt == ".markdown" {
		fileType = "markdown"
	} else {
		http.Error(w, "不支持的文件类型，只允许上传PDF或Markdown文件", http.StatusBadRequest)
		return
	}

	// 保存文件
	result, err := merger.SaveUploadedFile(file, fileName, tempDir)
	if err != nil {
		http.Error(w, "保存上传的文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 添加文件类型信息
	result.FileType = fileType

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// handleListTempFiles 获取临时目录中的文件列表
func handleListTempFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取临时目录参数
	tempDir := r.URL.Query().Get("dir")
	if tempDir == "" {
		http.Error(w, "必须指定临时目录参数 '?dir=...'", http.StatusBadRequest)
		return
	}

	// 获取目录中的文件
	files, err := merger.ListFilesInTempDir(tempDir)
	if err != nil {
		http.Error(w, "获取临时目录中的文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 区分PDF和Markdown文件
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

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tempDir":    tempDir,
		"allFiles":   files,
		"pdfFiles":   pdfFiles,
		"mdFiles":    mdFiles,
		"totalFiles": len(files),
	})
}

// handleMergeUploadedFiles 处理合并上传的文件
func handleMergeUploadedFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req MergeFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的JSON请求: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.TempDir == "" {
		http.Error(w, "必须指定临时目录", http.StatusBadRequest)
		return
	}

	if req.OutputFile == "" {
		http.Error(w, "必须指定输出文件名", http.StatusBadRequest)
		return
	}

	// 确保输出文件路径是绝对路径
	if !filepath.IsAbs(req.OutputFile) {
		absPath, err := filepath.Abs(req.OutputFile)
		if err == nil {
			req.OutputFile = absPath
		}
	}

	var result *merger.MergeResult
	var err error

	// 获取临时目录中的所有文件
	var filesToMerge []string

	if len(req.FileNames) > 0 {
		// 使用指定的文件名
		for _, fileName := range req.FileNames {
			filesToMerge = append(filesToMerge, filepath.Join(req.TempDir, fileName))
		}
	} else {
		// 使用目录中的所有文件
		filesToMerge, err = merger.ListFilesInTempDir(req.TempDir)
		if err != nil {
			http.Error(w, "获取临时目录中的文件失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(filesToMerge) == 0 {
		http.Error(w, "没有找到要合并的文件", http.StatusBadRequest)
		return
	}

	// 根据第一个文件的扩展名判断文件类型
	fileExt := strings.ToLower(filepath.Ext(filesToMerge[0]))

	if fileExt == ".pdf" {
		// 合并PDF文件
		result, err = merger.MergePDFFiles(filesToMerge, req.OutputFile, false)
	} else if fileExt == ".md" || fileExt == ".markdown" {
		// 合并Markdown文件
		result, err = merger.MergeMarkdownFilesList(filesToMerge, req.OutputFile, req.AddTitles, false)
	} else {
		http.Error(w, "不支持的文件类型，只能合并PDF或Markdown文件", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "合并文件失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
