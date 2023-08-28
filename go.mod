module github.com/server/txs-analysis

go 1.16

require (
	github.com/astaxie/beego v1.12.3
	github.com/ethereum/go-ethereum v1.10.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/sirupsen/logrus v1.9.0
	github.com/vesoft-inc/nebula-go/v3 v3.3.1 // indirect
	github.com/zhihu/norm v0.1.11
)

replace github.com/zhihu/norm => ../norm
