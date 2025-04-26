package mergemd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/liliang-cn/pdf-merger/pkg/merger"

	"github.com/spf13/cobra"
)

var (
	inputDir   string
	outputFile string
	addTitles  bool
	verbose    bool
	files      []string // Added: directly specify file list
)

// NewMergeMdCommand creates merge-md subcommand
func NewMergeMdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge-md",
		Short: "Merge Markdown files",
		Long:  `Merge all Markdown files in the specified directory, or merge the specified list of Markdown files, sorted in alphanumeric order`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMergeMd()
		},
	}

	// Add command line parameters
	cmd.Flags().StringVarP(&inputDir, "input", "i", ".", "Specify input directory containing Markdown files to merge")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "merged.md", "Specify output filename")
	cmd.Flags().BoolVarP(&addTitles, "add-titles", "t", true, "Add title for each file (using filename)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display detailed information")
	cmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "Specify list of Markdown files to merge, ignores input parameter if provided") // Added: file list parameter

	return cmd
}

func runMergeMd() error {
	var result *merger.MergeResult
	var err error

	// Ensure output file path is absolute
	if !filepath.IsAbs(outputFile) {
		absPath, err := filepath.Abs(outputFile)
		if err != nil {
			fmt.Printf("Warning: Unable to get absolute path: %v, will use relative path\n", err)
		} else {
			outputFile = absPath
		}
	}

	if verbose {
		fmt.Printf("Output file: %s\n", outputFile)
		fmt.Printf("Add titles: %v\n", addTitles)
	}

	// Choose processing mode based on parameters: file list or directory
	if len(files) > 0 {
		// Use specified file list
		if verbose {
			fmt.Printf("Will merge %d specified Markdown files\n", len(files))
		}
		result, err = merger.MergeMarkdownFilesList(files, outputFile, addTitles, verbose)
	} else {
		// Use directory mode
		// Ensure input directory path exists and is accessible
		inputInfo, err := os.Stat(inputDir)
		if err != nil {
			// Try to check if it's a path issue, not file doesn't exist
			if os.IsNotExist(err) {
				fmt.Printf("Error: Input directory does not exist: %s\n", inputDir)
				fmt.Println("Note: If using absolute path, ensure path is completely correct")
			} else {
				fmt.Printf("Error: Cannot access input directory: %v\n", err)
			}
			return err
		}

		if !inputInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", inputDir)
		}

		if verbose {
			fmt.Printf("Input directory: %s\n", inputDir)
		}

		result, err = merger.MergeMarkdownFiles(inputDir, outputFile, addTitles, verbose)
	}

	if err != nil {
		return err
	}

	fmt.Printf("Success! %d Markdown files merged into: %s\n", result.MergedFiles, result.OutputPath)
	return nil
}
