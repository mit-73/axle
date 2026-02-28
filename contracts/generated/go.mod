module github.com/ApeironFoundation/axle/contracts

go 1.26

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	google.golang.org/protobuf/cmd/protoc-gen-go
)

require (
	connectrpc.com/connect v1.19.1
	google.golang.org/protobuf v1.36.11
)
