# garden
> go language backend playground

create package
```shell
go mod init github.com/stkali/garden
```

## cobra & cobra-cli
install
```shell
go get -u github.com/spf13/cobra@latest
go get -u github.com/spf13/cobra-cli@latest
```
create app structure use cobra and cobra-cli
```shell
cobra-cli init -a "<st·kali clarkmonkey@163.com>" -l MIT --viper
cobra-cli add config -a "<st·kali clarkmonkey@163.com>" -l MIT
cobra-cli add server -a "<st·kali clarkmonkey@163.com>" -l MIT
```

## migrate
import golang migration

[🏠 golang/migrate](https://github.com/golang-migrate/migrate)

[👉 download list](https://github.com/golang-migrate/migrate/releases)

```shell
migrate create -ext sql -dir db/migrations -seq desc_table
```

## sqlc

[🏠 kyleconroy/sqlc](https://github.com/kyleconroy/sqlc)

[👉 download list](https://github.com/kyleconroy/sqlc/releases)

```shell
sqlc generate
```
[📰 sqlc doc](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html)

sqlc config file sample
```yaml
version: 1
packages:
  - path: "db/sqlc"
    name: "db"
    engine: "postgresql"
    schema: "db/migration/"
    queries: "db/query"
    emit_json_tags: true
    emit_prepared_queries: false
    emit_interface: true
    emit_exact_table_names: false
    emit_empty_slices: true
```

genreate go code from sql query
```shell
sqlc generate
```

## protoc & gRPC
install proto buffer compiler

[🏠 protocolbuffers/protobuf](https://github.com/protocolbuffers/protobuf)

[👉 download list](https://github.com/protocolbuffers/protobuf/releases)

install Golang plugins for the protocol compiler:
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Update your PATH so that the protoc compiler can find the plugins:
```shell
export PATH="$PATH:$(go env GOPATH)/bin"
# or write it to profile or zshrc file
```

## evans

[🏠 ktr0731/evans](https://github.com/ktr0731/evans)

[👉 download list](https://github.com/ktr0731/evans/releases)

> evans is a gRPC client
```
# regsiter reflection

reflection.Register(grpcServer)
```
connect rpc server
```shell
evans -r repl --host <server address> --port <server port>
```

## grpc-gateway
generate http server by protobuf and grpc
[🏠 grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

[👉 download list](https://github.com/grpc-ecosystem/grpc-gateway/releases)

1 copy *.proto to project
```go
// +build tools

package tools

import (
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

```shell
# pull dependencies
go mod tidy

# install binaries
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
add file to $PROJECT/proto/google/api/*.proto

[🏠 googleapis/googleapis](https://github.com/googleapis/googleapis)
```shell
google/api/annotations.proto
google/api/field_behavior.proto
google/api/http.proto
google/api/httpbody.proto
```


