package train

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/emehrkay/cbt/internal/models/cdc"
	"github.com/emehrkay/cbt/internal/models/jwt/issuer"
	"github.com/emehrkay/cbt/internal/models/jwt/validator"
	"github.com/emehrkay/cbt/internal/models/service"
	"github.com/emehrkay/cbt/internal/storage"
	"github.com/emehrkay/cbt/internal/storage/demo"
)

var (
	err            error
	trainService   *service.Train
	store          storage.Storage
	rpcPort        string
	jwtIssuer      issuer.JWTIssuer
	jwtValidator   validator.JWTValidator
	capture        cdc.CDC
	privateKeyFile string
	publicKeyFile  string
	RootCmd        *cobra.Command
)

func init() {
	RootCmd = &cobra.Command{
		Use:   "train",
		Short: "commands associated with train service",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rpcPort = os.Getenv("RPC_PORT")
	if rpcPort == "" {
		rpcPort = ":7676"
	}

	store = demo.DemoStorage{}

	privateKeyFile = os.Getenv("PRIVATE_KEY_FILE")
	if privateKeyFile == "" {
		privateKeyFile = "./train/key"
	}

	privateKeyFile, err = filepath.Abs(privateKeyFile)
	if err != nil {
		panic(err)
	}

	publicKeyFile = os.Getenv("PUBLIC_KEY_FILE")
	if publicKeyFile == "" {
		publicKeyFile = "./train/key.pub"
	}

	publicKeyFile, err = filepath.Abs(publicKeyFile)
	if err != nil {
		panic(err)
	}

	jwtIssuer, err = issuer.New(privateKeyFile, rpcPort)
	if err != nil {
		panic(err)
	}

	jwtValidator, err = validator.New(publicKeyFile)
	if err != nil {
		panic(err)
	}

	capture = cdc.New()
	trainService, err = service.New(store, jwtIssuer, jwtValidator, capture)
	if err != nil {
		panic(err)
	}
}

func log(format string, args ...any) {
	fmt.Printf(format, args...)
	fmt.Println("")
	fmt.Println("")
}
