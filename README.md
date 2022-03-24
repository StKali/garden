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

[golang/migrate](https://github.com/golang-migrate/migrate)

download: https://github.com/golang-migrate/migrate/releases


```shell
migrate create -ext sql -dir db/migrations -seq desc_table
```

