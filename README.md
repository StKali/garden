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
