package utils

import "github.com/server/txs-analysis/constant"

func GetMethod(inputData string) string {
	if len(inputData) >= 10 {
		switch inputData[:10] {
		case constant.TransferCode:
			return constant.TransferMethod
		case constant.TransferEventCode:
			return constant.TransferMethod
		case constant.MintCode:
			return constant.MintMethod
		default:
			return ""
		}
	}
	return ""
}
