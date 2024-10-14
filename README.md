

protobuf compiler

Install the `protoc` tool using the instructions available at [https://grpc.io/docs/protoc-installation/]


 Go protobuf codegen tools

`go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Generate Go code from .proto files (if using window)
```
protoc --go_out=. --go_opt=paths=source_relative ^
  --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
  proto/mail.proto
```
