module github.com/server/txs-analysis

go 1.16

require (
	github.com/astaxie/beego v1.12.3
	github.com/ethereum/go-ethereum v1.10.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/shopspring/decimal v1.3.1
	github.com/vesoft-inc/nebula-go/v3 v3.3.1
	github.com/zhihu/norm v0.1.11
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/zhihu/norm => ../norm
