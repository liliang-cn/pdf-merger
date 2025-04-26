# File Merger

一个强大的文件合并工具，可以将 PDF 和 Markdown 文件按字符顺序合并。提供命令行界面和 HTTP API 两种使用方式。

## 功能特性

- 合并 PDF 文件：将多个 PDF 文件合并为一个 PDF 文件
- 合并 Markdown 文件：将多个 Markdown 文件合并为一个 Markdown 文件
- 文件排序：按字符顺序排序所有文件
- 命令行界面：使用 Cobra 框架提供易用的命令行工具
- HTTP API：提供 REST API 接口，可远程调用
- 批量处理：一次处理目录下的所有相同类型文件

## 安装

### 从源码编译

```bash
git clone <repository-url>
cd pdf-merger
go build -o file-merger
```

## 使用方法

### 命令行模式

**查看帮助信息:**

```bash
./pdf-merger --help
```

**合并 PDF 文件:**

```bash
./pdf-merger merge -i <输入目录> -o <输出文件.pdf> -v
```

参数说明:

- `-i, --input`: 指定输入目录 (默认为当前目录)
- `-o, --output`: 指定输出文件名 (默认为 merged.pdf)
- `-v, --verbose`: 显示详细信息

**合并 Markdown 文件:**

```bash
./pdf-merger merge-md -i <输入目录> -o <输出文件.md>
```

参数说明:

- `-i, --input`: 指定输入目录 (默认为当前目录)
- `-o, --output`: 指定输出文件名 (默认为 merged.md)
- `-t, --add-titles`: 是否为每个文件添加标题 (默认为 true)
- `-v, --verbose`: 显示详细信息

### API 服务器模式

**启动 API 服务器:**

```bash
./pdf-merger serve -p 8080
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

## 示例

### 合并所有 PDF 教程

```bash
./pdf-merger merge -i .. -o "Rust教程合集.pdf" -v
```

### 合并所有 Markdown 笔记并添加标题

```bash
./pdf-merger merge-md -i ./notes -o "笔记合集.md" -t -v
```

## 依赖库

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - 命令行界面框架
- [github.com/pdfcpu/pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF 处理库

## 项目结构

```
pdf-merger/
├── api/           # API服务器实现
│   └── server.go  # HTTP API处理逻辑
├── cmd/           # 命令行界面实现
│   ├── root.go    # 根命令
│   ├── merge/     # PDF合并命令
│   ├── merge-md/  # Markdown合并命令
│   └── serve/     # API服务器命令
├── pkg/           # 核心功能包
│   └── merger/    # 文件合并核心逻辑
├── main.go        # 主程序入口
└── README.md      # 项目文档
```

## 适用场景

- 合并教程 PDF 文件为一个大文档
- 整合多个 Markdown 笔记为一个文档
- 在不同设备上调用 API 进行文件合并
- 集成到其他系统作为文件处理服务

## 许可证

[MIT](https://opensource.org/licenses/MIT)
