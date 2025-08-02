dogecli
> 多吉云基础型云储存管理工具

## 安装

go install

```shell
go install github.com/dongfg/dogecli@latest
```

shell install on unix
```shell
curl -sSL https://raw.githubusercontent.com/dongfg/dogecli/refs/heads/master/scripts/install.sh | bash
```

or [download binary](https://github.com/dongfg/dogecli/releases)

## 使用

> 先使用 dogecli config 进行配置

```shell
dogecli -h

多吉云基础型云储存管理工具

Usage:
  dogecli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Interactively set AccessKey and SecretKey
  copy        Upload file to bucket
  fetch       Download file to bucket
  help        Help about any command
  list        List files in bucket
  list-bucket List all buckets
  list-fetch  Get fetch status
  version     Print the version number of dogecli

Flags:
      --config string   config file (default is $HOME/.dogecli/config.yaml)
  -h, --help            help for dogecli
  -v, --verbose         enable verbose output

Use "dogecli [command] --help" for more information about a command.
```