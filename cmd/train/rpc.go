package train

import (
	"github.com/emehrkay/cbt/internal/api/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverRunning bool
)

func startRPCServer() {
	if serverRunning {
		return
	}

	rpcServer, err := rpc.New(trainService)
	if err != nil {
		panic(err)
	}

	err = rpcServer.Run(rpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	serverRunning = true
}

func init() {
	rpcCmd := &cobra.Command{
		Use:   "rpc",
		Short: "starts the rpc server",
		Run: func(cmd *cobra.Command, args []string) {
			startRPCServer()
		},
	}

	RootCmd.AddCommand(rpcCmd)
}
