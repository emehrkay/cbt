.EXPORT_ALL_VARIABLES:

DIR = $(dir $(abspath $(firstword $(MAKEFILE_LIST))))
RPC_PORT = :7676
PRIVATE_KEY_FILE = $(DIR)/cmd/train/key
PUBLIC_KEY_FILE = $(DIR)/cmd/train/key.pub

protos:
	@protoc --go_out=pkg/rpc --go_opt=paths=source_relative --go-grpc_out=pkg/rpc --go-grpc_opt=paths=source_relative ./pkg/rpc/protos/train.proto --proto_path=./pkg/rpc/protos

help:
	@go run cmd/main.go

ex1:
	@go run cmd/main.go ex1

ex2:
	@go run cmd/main.go ex2

ex3:
	@go run cmd/main.go ex3

ex4:
	@go run cmd/main.go ex4

ex5:
	@go run cmd/main.go ex5

tests:
	@ginkgo -v -r
