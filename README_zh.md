# File Merger

一个强大的文件合并工具，可以将 PDF 和 Markdown 文件按字符顺序合并。提供命令行界面和 HTTP API 两种使用方式，支持目录扫描和直接指定文件。

## 功能特性

- 合并 PDF 文件：将多个 PDF 文件合并为一个 PDF 文件
- 合并 Markdown 文件：将多个 Markdown 文件合并为一个 Markdown 文件
- 文件排序：按字符顺序排序所有文件
- 直接指定文件：可以直接指定要合并的具体文件列表
- 支持绝对路径：完全支持绝对路径和相对路径
- 支持文件上传：可以上传文件到临时目录并进行合并
- 命令行界面：使用 Cobra 框架提供易用的命令行工具
- HTTP API：提供 REST API 接口，可远程调用
- 批量处理：一次处理目录下的所有相同类型文件
- 多平台支持：可用于 Linux、macOS（包括 Apple Silicon）和 Windows

## 安装

### 使用 Go Install 安装（推荐）

安装 File Merger 最简单的方法是使用 Go 的内置 `go install` 命令：

```bash
go install github.com/liliang-cn/pdf-merger@latest
```

此命令会下载、编译并将 File Merger 的最新版本安装到您的 `$GOPATH/bin` 目录中。请确保该目录已添加到您系统的 PATH 环境变量中。

### 从发布页面下载

访问项目的 [GitHub Releases](https://github.com/liliang-cn/pdf-merger/releases) 页面，下载适合您操作系统的预编译二进制文件。我们提供以下平台的二进制文件：

- Linux (amd64)
- macOS (Intel/amd64 和 Apple Silicon/arm64)
- Windows (amd64)

### 从源码编译

如果您喜欢，也可以从源码构建：

```bash
git clone https://github.com/liliang-cn/pdf-merger.git
cd pdf-merger
go build -o file-merger
```

## 使用方法

### 命令行模式

**查看帮助信息:**

```bash
file-merger --help
```

**合并 PDF 文件 (目录模式):**

```bash
file-merger merge -i <输入目录> -o <输出文件.pdf> -v
```

**合并 PDF 文件 (指定文件模式):**

```bash
file-merger merge -f <文件1.pdf> <文件2.pdf> <文件3.pdf> -o <输出文件.pdf> -v
```

参数说明:

- `-i, --input`: 指定输入目录 (默认为当前目录)
- `-o, --output`: 指定输出文件名 (默认为 merged.pdf)
- `-f, --files`: 指定要合并的 PDF 文件列表 (如果提供则忽略 input 参数)
- `-v, --verbose`: 显示详细信息

**合并 Markdown 文件 (目录模式):**

```bash
file-merger merge-md -i <输入目录> -o <输出文件.md>
```

**合并 Markdown 文件 (指定文件模式):**

```bash
file-merger merge-md -f <文件1.md> <文件2.md> <文件3.md> -o <输出文件.md>
```

参数说明:

- `-i, --input`: 指定输入目录 (默认为当前目录)
- `-o, --output`: 指定输出文件名 (默认为 merged.md)
- `-f, --files`: 指定要合并的 Markdown 文件列表 (如果提供则忽略 input 参数)
- `-t, --add-titles`: 是否为每个文件添加标题 (默认为 true)
- `-v, --verbose`: 显示详细信息

### API 服务器模式

**启动 API 服务器:**

```bash
file-merger serve -p 8080
```

参数说明:

- `-p, --port`: 指定 API 服务器监听端口 (默认为 8080)

**API 端点:**

1. **获取目录中的 PDF 文件列表:**

```bash
curl -X GET "http://localhost:8080/api/files?dir=<目录路径>"
```

2. **获取目录中的 Markdown 文件列表:**

```bash
curl -X GET "http://localhost:8080/api/md-files?dir=<目录路径>"
```

3. **合并 PDF 文件:**

```bash
curl -X POST "http://localhost:8080/api/merge" \
     -H "Content-Type: application/json" \
     -d '{"inputDir": "<目录路径>", "outputFile": "output.pdf"}'
```

4. **合并 Markdown 文件:**

```bash
curl -X POST "http://localhost:8080/api/merge-md" \
     -H "Content-Type: application/json" \
     -d '{"inputDir": "<目录路径>", "outputFile": "output.md", "addTitles": true}'
```

5. **下载合并后的文件:**

```bash
curl -X GET "http://localhost:8080/api/download/<文件路径>" --output downloaded_file
```

### 文件上传和临时目录 API

1. **创建临时目录:**

```bash
curl -X POST "http://localhost:8080/api/temp-dir"
```

2. **上传文件到临时目录:**

```bash
curl -X POST "http://localhost:8080/api/upload" \
     -F "tempDir=<临时目录路径>" \
     -F "file=@<要上传的文件路径>"
```

3. **列出临时目录中的文件:**

```bash
curl -X GET "http://localhost:8080/api/temp-files?dir=<临时目录路径>"
```

4. **合并临时目录中的文件:**

```bash
curl -X POST "http://localhost:8080/api/merge-files" \
     -H "Content-Type: application/json" \
     -d '{"tempDir": "<临时目录路径>", "outputFile": "merged.pdf", "addTitles": true}'
```

5. **删除临时目录:**

```bash
curl -X DELETE "http://localhost:8080/api/temp-dir" \
     -H "Content-Type: application/json" \
     -d '{"tempDir": "<临时目录路径>"}'
```

## 示例

### 合并所有 PDF 教程

```bash
file-merger merge -i ./PDF -o "Kubernetes教程合集.pdf" -v
```

### 合并指定的 PDF 文件（支持绝对路径）

```bash
file-merger merge -f "/Users/liliang/Things/backend/common/PDF/01｜初识容器：万事开头难.pdf" "/Users/liliang/Things/backend/common/PDF/02｜被隔离的进程：一起来看看容器的本质.pdf" -o "k8s入门.pdf" -v
```

### 合并所有 Markdown 笔记并添加标题

```bash
file-merger merge-md -i ./notes -o "笔记合集.md" -t -v
```

### 合并指定的 Markdown 文件

```bash
file-merger merge-md -f "intro.md" "chapter1.md" "chapter2.md" -o "文档.md" -v
```

## 依赖库

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - 命令行界面框架
- [github.com/pdfcpu/pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF 处理库

## 项目结构

```
file-merger/
├── api/                 # API服务器实现
│   └── server.go        # HTTP API处理逻辑
├── cmd/                 # 命令行界面实现
│   ├── root.go          # 根命令
│   ├── merge/           # PDF合并命令
│   ├── merge-md/        # Markdown合并命令
│   └── serve/           # API服务器命令
├── pkg/                 # 核心功能包
│   └── merger/          # 文件合并核心逻辑
│       ├── merger.go    # 合并功能实现
│       └── filemanager.go # 文件管理功能实现
├── main.go              # 主程序入口
├── README.md            # 项目文档（英文）
└── README_zh.md         # 项目文档（中文）
```

## 适用场景

- 合并教程 PDF 文件为一个大文档
- 整合多个 Markdown 笔记为一个文档
- 使用 API 通过 Web UI 上传和合并文件
- 处理任意路径（包括绝对路径）的文件
- 在不同设备上调用 API 进行文件合并
- 集成到其他系统作为文件处理服务

## 贡献

欢迎贡献！请随时提出问题或提交拉取请求。

## 许可证

[MIT](https://opensource.org/licenses/MIT)
