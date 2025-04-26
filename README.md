# File Merger

A powerful file merging tool that can combine PDF and Markdown files in alphanumeric order. It provides both a command-line interface and HTTP API, supporting directory scanning and direct file specification.

[中文文档](README_zh.md)

## Features

- Merge PDF Files: Combine multiple PDF files into a single PDF file
- Merge Markdown Files: Combine multiple Markdown files into a single Markdown document
- File Sorting: Sort all files in alphanumeric order
- Direct File Specification: Specify exact files to merge
- Absolute Path Support: Full support for both absolute and relative paths
- File Upload Support: Upload files to a temporary directory for merging
- Command-line Interface: Easy-to-use CLI tool built with the Cobra framework
- HTTP API: RESTful API interface for remote operation
- Batch Processing: Process all files of the same type in a directory at once
- Multi-platform Support: Available for Linux, macOS (including Apple Silicon), and Windows

## Installation

### Using Go Install (Recommended)

The easiest way to install File Merger is with Go's built-in `go install` command:

```bash
go install github.com/liliang-cn/pdf-merger@latest
```

This command will download, compile, and install the latest version of File Merger to your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH.

### Download Pre-built Binaries

Visit the [GitHub Releases](https://github.com/liliang-cn/pdf-merger/releases) page to download pre-compiled binaries for your operating system. We provide binaries for:

- Linux (amd64)
- macOS (Intel/amd64 and Apple Silicon/arm64)
- Windows (amd64)

### Build from Source

You can also build from source if you prefer:

```bash
git clone https://github.com/liliang-cn/pdf-merger.git
cd pdf-merger
go build -o file-merger
```

## Usage

### Command Line Mode

**View help information:**

```bash
file-merger --help
```

**Merge PDF files (directory mode):**

```bash
file-merger merge -i <input_directory> -o <output_file.pdf> -v
```

**Merge PDF files (specific files mode):**

```bash
file-merger merge -f <file1.pdf> <file2.pdf> <file3.pdf> -o <output_file.pdf> -v
```

Parameter description:

- `-i, --input`: Specify the input directory (default is the current directory)
- `-o, --output`: Specify the output filename (default is merged.pdf)
- `-f, --files`: Specify the list of PDF files to merge (ignores the input parameter if provided)
- `-v, --verbose`: Display detailed information

**Merge Markdown files (directory mode):**

```bash
file-merger merge-md -i <input_directory> -o <output_file.md>
```

**Merge Markdown files (specific files mode):**

```bash
file-merger merge-md -f <file1.md> <file2.md> <file3.md> -o <output_file.md>
```

Parameter description:

- `-i, --input`: Specify the input directory (default is the current directory)
- `-o, --output`: Specify the output filename (default is merged.md)
- `-f, --files`: Specify the list of Markdown files to merge (ignores the input parameter if provided)
- `-t, --add-titles`: Whether to add titles for each file (default is true)
- `-v, --verbose`: Display detailed information

### API Server Mode

**Start the API server:**

```bash
file-merger serve -p 8080
```

Parameter description:

- `-p, --port`: Specify the API server listening port (default is 8080)

**API Endpoints:**

1. **Get a list of PDF files in a directory:**

```bash
curl -X GET "http://localhost:8080/api/files?dir=<directory_path>"
```

2. **Get a list of Markdown files in a directory:**

```bash
curl -X GET "http://localhost:8080/api/md-files?dir=<directory_path>"
```

3. **Merge PDF files:**

```bash
curl -X POST "http://localhost:8080/api/merge" \
     -H "Content-Type: application/json" \
     -d '{"inputDir": "<directory_path>", "outputFile": "output.pdf"}'
```

4. **Merge Markdown files:**

```bash
curl -X POST "http://localhost:8080/api/merge-md" \
     -H "Content-Type: application/json" \
     -d '{"inputDir": "<directory_path>", "outputFile": "output.md", "addTitles": true}'
```

5. **Download the merged file:**

```bash
curl -X GET "http://localhost:8080/api/download/<file_path>" --output downloaded_file
```

### File Upload and Temporary Directory API

1. **Create a temporary directory:**

```bash
curl -X POST "http://localhost:8080/api/temp-dir"
```

2. **Upload file to a temporary directory:**

```bash
curl -X POST "http://localhost:8080/api/upload" \
     -F "tempDir=<temp_dir_path>" \
     -F "file=@<file_to_upload_path>"
```

3. **List files in a temporary directory:**

```bash
curl -X GET "http://localhost:8080/api/temp-files?dir=<temp_dir_path>"
```

4. **Merge files in a temporary directory:**

```bash
curl -X POST "http://localhost:8080/api/merge-files" \
     -H "Content-Type: application/json" \
     -d '{"tempDir": "<temp_dir_path>", "outputFile": "merged.pdf", "addTitles": true}'
```

5. **Delete a temporary directory:**

```bash
curl -X DELETE "http://localhost:8080/api/temp-dir" \
     -H "Content-Type: application/json" \
     -d '{"tempDir": "<temp_dir_path>"}'
```

## Examples

### Merge all PDF tutorials

```bash
file-merger merge -i ./PDF -o "Kubernetes_Tutorials.pdf" -v
```

### Merge specific PDF files (supporting absolute paths)

```bash
file-merger merge -f "/Users/liliang/Things/backend/common/PDF/01_intro_to_containers.pdf" "/Users/liliang/Things/backend/common/PDF/02_isolated_processes.pdf" -o "k8s_intro.pdf" -v
```

### Merge all Markdown notes with titles

```bash
file-merger merge-md -i ./notes -o "combined_notes.md" -t -v
```

### Merge specific Markdown files

```bash
file-merger merge-md -f "intro.md" "chapter1.md" "chapter2.md" -o "document.md" -v
```

## Dependencies

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - Command-line interface framework
- [github.com/pdfcpu/pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF processing library

## Project Structure

```
file-merger/
├── api/                 # API server implementation
│   └── server.go        # HTTP API handling logic
├── cmd/                 # Command-line interface implementation
│   ├── root.go          # Root command
│   ├── merge/           # PDF merge command
│   ├── merge-md/        # Markdown merge command
│   └── serve/           # API server command
├── pkg/                 # Core functionality packages
│   └── merger/          # File merging core logic
│       ├── merger.go    # Merge functionality implementation
│       └── filemanager.go # File management implementation
├── main.go              # Main program entry
├── README.md            # Project documentation (English)
└── README_zh.md         # Project documentation (Chinese)
```

## Use Cases

- Merge tutorial PDF files into a single comprehensive document
- Combine multiple Markdown notes into one document
- Use API through Web UI to upload and merge files
- Process files with arbitrary paths (including absolute paths)
- Call the API from different devices for file merging
- Integrate into other systems as a file processing service

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

[MIT](https://opensource.org/licenses/MIT)
