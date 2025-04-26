package cmd

import (
	"os"

	"github.com/liliang-cn/pdf-merger/cmd/merge"
	mergemd "github.com/liliang-cn/pdf-merger/cmd/merge-md"
	"github.com/liliang-cn/pdf-merger/cmd/serve"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "file-merger",
		Short: "File Merger Tool",
		Long:  `A command-line tool and API server for merging PDF and Markdown files, capable of combining all files in a specified directory into one file, sorted in alphanumeric order.`,
	}

	// Add subcommands
	rootCmd.AddCommand(merge.NewMergeCommand())
	rootCmd.AddCommand(mergemd.NewMergeMdCommand())
	rootCmd.AddCommand(serve.NewServeCommand())
}
