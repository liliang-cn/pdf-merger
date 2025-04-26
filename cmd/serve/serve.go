package serve

import (
	"pdf-merger/api"

	"github.com/spf13/cobra"
)

var port int

// NewServeCommand 创建serve子命令
func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "启动API服务器",
		Long:  `启动HTTP API服务器，提供PDF合并功能的REST接口`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe()
		},
	}

	// 添加命令行参数
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "API服务器监听端口")

	return cmd
}

func runServe() error {
	return api.StartServer(port)
}
