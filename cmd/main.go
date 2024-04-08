package main

import (
	"github.com/emehrkay/cbt/cmd/train"
)

func main() {
	if err := train.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
