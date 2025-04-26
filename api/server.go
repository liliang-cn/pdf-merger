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

// StartServer 启动API服务器
func StartServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("API服务器启动在 http://localhost%s\n", addr)
	fmt.Printf("可用端点:\n")
	fmt.Printf("  POST /api/merge         - 合并PDF文件\n")
	fmt.Printf("  POST /api/merge-md      - 合并Markdown文件\n")
	fmt.Printf("  GET  /api/files?dir=... - 列出目录中的PDF文件\n")
	fmt.Printf("  GET  /api/md-files?dir=... - 列出目录中的Markdown文件\n")

	// 注册API路由处理程序
	http.HandleFunc("/api/merge", handleMerge)
	http.HandleFunc("/api/merge-md", handleMergeMd)
	http.HandleFunc("/api/files", handleListFiles)
	http.HandleFunc("/api/md-files", handleListMdFiles)
	http.HandleFunc("/api/download/", handleDownload)

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
