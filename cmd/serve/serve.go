package serve

import (
	"github.com/liliang-cn/pdf-merger/api"

	"github.com/spf13/cobra"
)

var port int

// NewServeCommand creates a serve subcommand
func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start API server",
		Long:  `Start HTTP API server, providing REST interface for PDF merging functionality`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe()
		},
	}

	// Add command line parameters
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "API server listening port")

	return cmd
}

func runServe() error {
	return api.StartServer(port)
}
