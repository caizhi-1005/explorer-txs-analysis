package utils

import "github.com/server/txs-analysis/constant"

func GetMethod(inputData string) string {
	if len(inputData) >= 10 {
		switch inputData[:10] {
		case constant.TransferCode:
			return constant.Transfer
		case constant.TransferEventCode:
			return constant.Transfer
		case constant.MintCode:
			return constant.Mint
		default:
			return ""
		}
	}
	return ""
}
