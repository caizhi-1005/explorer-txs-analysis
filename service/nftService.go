package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
	"github.com/server/txs-analysis/utils"
)

type NftService struct {
}

// NFTList NFT溯源-全部NFT列表下拉列表
func (this *NftService) NFTList() ([]*apiModels.RespNftList, error) {
	list, err := dbModels.NFTList()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// TokenIdList NFT溯源-Token ID下拉列表
func (this *NftService) TokenIDList(req apiModels.ReqTokenIDList) ([]string, error) {
	list, err := dbModels.TokenIDList(req)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// NFTAddressDetail NFT溯源-地址详情
// todo 改我nebula查询
//func (this *NftService) NFTAddressDetail(req apiModels.ReqNFTAddressDetail) (*apiModels.RespNFTAddressDetail, error) {
	//	地址类型
	//	持有的 Token ID
	//	历史持有过的 Token ID
	//	最长持有 Token ID
	//	最长持有时间
	//todo 该地址没有所选的NFT交易记录

	// 如果当前地址在当前合约里持有的NFT数量超过100个，不支持显示
	//tokenIdCount, err := dbModels.HoldTokenIdCount(req.ContractAddress, req.Address)
	//if err != nil {
	//	beego.Error("dbModels.HoldTokenIdCount error.", err)
	//}
	//if tokenIdCount > 100 {
	//	return nil, errors.New(constant.ErrLimit)
	//}
	//
	//tokenIdAndTypes, err := dbModels.HoldTokenIdAndAddressType(req.ContractAddress, req.Address)
	//if err != nil {
	//	beego.Error("dbModels.HoldTokenIdAndAddressType error.", err)
	//}
	//tokenIdHistorys, err := dbModels.HoldTokenIdHistory(req.ContractAddress, req.Address)
	//if err != nil {
	//	beego.Error("dbModels.HoldTokenIdHistory error.", err)
	//}
	//contractTxs, err := dbModels.LongestHold(req.ContractAddress, req.Address)
	//if err != nil {
	//	beego.Error("dbModels.LongestHold error.", err)
	//	return nil, err
	//}
	//fmt.Println("contractTxs:", contractTxs)
	//var startTime time.Time
	//var endTime time.Time
	//var holdTime time.Duration
	//var holdNftNum int
	//var holdNftNumHistory int
	//
	//var holdTokenIds []apiModels.HoldTokenIds
	//
	//var holdTokenIdsHistory []apiModels.HoldTokenIdsHistory
	//
	//for _, v := range contractTxs {
	//
	//	var startTime time.Time
	//	var endTime time.Time
	//	var holdTime time.Duration
	//
	//	if v.To == req.Address {
	//		var holdTokenId apiModels.HoldTokenIds
	//		holdTokenId.HoldTokenId = v.TokenId
	//		holdTokenIds = append(holdTokenIds, holdTokenId)
	//		holdTokenIds = append(holdTokenIds[:i], holdTokenIds[i+1:]...)
	//		startTime = v.TxTime
	//		holdNftNum ++
	//	}
	//	if v.From == req.Address {
	//		endTime = v.TxTime
	//	}
	//	holdTime = endTime.Sub(startTime)
	//}
	//fmt.Println("holdTime--->",holdTime)
	//
	//var res *apiModels.RespNFTAddressDetail
	//if len(tokenIdAndTypes) > 0 {
	//	res.AddressType = tokenIdAndTypes[0].AccountType
	//}
	//res.HoldTokenIdList = tokenIdAndTypes
	//res.HoldTokenIdHistoryList = tokenIdHistorys
	//return res, nil
//}

// NFTTransferDetails NFT溯源-地址详情-NFT流转详情列表
func (this *NftService) NFTTransferDetailByAddress(req apiModels.ReqNFTTransferDetailsByAddress) ([]*apiModels.RespNFTTransferDetailsByAddress, error) {
	list, err := dbModels.NFTTransferDetails(req)
	if err != nil {
		return nil, err
	}
	//todo 需要对token_id转换为数字格式，且从小到大排序
	//todo 流转详情：显示的最后一次交易from to，统计当前合约地址，当前token_id交易次数
	for _, v := range list {

		fmt.Println("---->",v.TokenId)
		tokenid, err := utils.ConvertTokenID(v.TokenId)
		if err != nil {
			fmt.Println("--err-->",err)
		}
		fmt.Println("--tokenid-->",tokenid)
	}

	return list, nil
}

// NFTTxDetail NFT溯源-交易详情
func (this *NftService) NFTTxDetail(req apiModels.ReqNFTTxDetail) (*apiModels.RespNFTTxDetail, error) {
	contractTx, count, err := dbModels.NFTTxDetail(req)
	if err != nil {
		return nil, err
	}
	if contractTx == nil {
		return nil, nil
	}

	//处理数据 method
	contractTx.Method = utils.GetMethod(contractTx.Method)

	//处理数据 token_id
	contractTx.TokenId = common.HexToHash(contractTx.TokenId).Big().String()

	contractTx.TxHash = req.TxHash
	contractTx.TransferCount = count
	return contractTx, nil
}

// NFTDetail NFT溯源-NFT详情
func (this *NftService) NFTDetail(req apiModels.ReqNFTDetail) (apiModels.RespNFTDetail, error) {
	//TransferCount      string `json:"transfer_count"`
	//HistoryHolderCount string `json:"history_holder_count"`
	//MintTime           string `json:"mint_time"`
	res := apiModels.RespNFTDetail{}
	detail, err := dbModels.NFTDetail(req)
	if err != nil {
		beego.Error("dbModels.NFTDetail error.", err)
		//return nil, err
	}
	if detail != nil {
		res.TransferCount = detail.TransferCount
		res.HistoryHolderCount = detail.HistoryHolderCount
		res.MintTime = detail.MintTime
	}

	contract, err := dbModels.ContractInfo(req)
	if err != nil {
		beego.Error("dbModels.ContractInfo error.", err)
	}

	if contract != nil {
		res.Name = contract.Name
		res.Symbol = contract.Symbol
		res.TokenType = contract.TokenType
		res.Holder = contract.Holder
	}

	res.ContractAddress = req.ContractAddress

	//todo nebula获取
	//LongestHoldTime    string `json:"longest_hold_time"`
	//CurrentHoldTime    string `json:"current_hold_time"`
	return res, nil
}


// NFTTransferDetail NFT溯源-交易详情
func (this *NftService) NFTTransferDetail(req apiModels.ReqNFTTxDetail) (*apiModels.RespNFTTxDetail, error) {
	list, _, err := dbModels.NFTTxDetail(req)
	if err != nil {
		return nil, err
	}
	return list, nil
}


// NFTTransferDetailByTokenId NFT溯源-NFT详情-NFT流转详情列表
func (this *NftService) NFTTransferDetailByTokenId(req apiModels.ReqNFTTransferDetailsByTokenId) ([]*apiModels.RespNFTTransferDetailsByTokenId, error) {
	list, err := dbModels.NFTTransferDetailByTokenId(req)
	if err != nil {
		return nil, err
	}

	//处理method
	for _, v := range list {
		if len(v.Method) >= 10 {
			v.Method = utils.GetMethod(v.Method)
		}
	}
	return list, nil
}


//-----------------------------nebula------------------------------

// NFTStartAnalysis NFT溯源-NFT开始分析-交易图
//func (this *NftService) NFTStartAnalysis(req apiModels.ReqNFTStartAnalysis) ([]*apiModels.RespNFTTransferDetailsByTokenId, error) {
//	// todo token_id为查询条件，返回单条数据
//
//	list, err := nebulaModels.TraceNFTTxs(req)
//	if err != nil {
//		return nil, err
//	}
//	return list, nil
//}
//
//// NFTTrace NFT溯源-NFT追溯-交易图
//func (this *NftService) NFTTrace(req apiModels.ReqNFTTrace) ([]*apiModels.RespNFTTransferDetailsByTokenId, error) {
//
//	list, err := nebulaModels.TraceNFTTxs(req)
//	if err != nil {
//		return nil, err
//	}
//	return list, nil
//}


