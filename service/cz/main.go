package main

import (
	"fmt"
	"github.com/server/txs-analysis/models/nebulaModels"
)

func main() {
	trace()


}

func trace()  {
	contract_address := "0xd801e8cf801ecd9b29c3ca0aec3ddcbe041af381"
	address := "0x260f2f029b5b985da0d0b7984dd2c75b269610b4"
	token_id := "10"
	count := "2"
	direction := "out"

	nebulaDB := nebulaModels.Init()
	txsRoute, err := nebulaModels.GetNFTTxsPath(nebulaDB, contract_address, address, token_id, count, direction)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("txsRoute:", txsRoute)
}
