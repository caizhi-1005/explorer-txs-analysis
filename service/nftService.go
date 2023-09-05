package service

import (
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
	"github.com/server/txs-analysis/models/nebulaModels"
	"github.com/server/txs-analysis/utils"
	"sort"
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
	var res []string
	for _, v := range list {
		//处理数据 token_id
		v = common.HexToHash(v).Big().String()
		res = append(res, v)
	}
	return res, nil
}

// NFTTransferDetails NFT溯源-地址详情-NFT流转详情列表
func (this *NftService) NFTTransferDetailByAddress(req apiModels.ReqNFTTransferDetailsByAddress) ([]*apiModels.RespNFTTransferDetailsByAddress, error) {
	list, err := dbModels.NFTTransferDetails(req)
	if err != nil {
		return nil, err
	}

	var keys []string
	mapData := make(map[string]*apiModels.RespNFTTransferDetailsByAddress)
	for _, v := range list {
		tokenId, err := utils.ConvertTokenID(v.TokenId)
		if err != nil {
			beego.Error("NFTTransferDetailByAddress utils.ConvertTokenID error.", err)
		}
		keys = append(keys, tokenId)
		mapData[tokenId] = v
	}
	sort.Strings(keys)

	var res []*apiModels.RespNFTTransferDetailsByAddress
	for i := 0; i < len(keys); i++ {
		res = append(res, mapData[keys[i]])
	}

	return res, nil
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

//NFTAddressDetail NFT溯源-地址详情
//func (this *NftService) NFTAddressDetail(req apiModels.ReqNFTAddressDetail) (*apiModels.RespNFTAddressDetail, error) {
//	//todo 怎么确定是否是最后一个或第一个
//	//from= 0x0000000000000000000000000000000000000000 是开始， to 是owner
//	//contractAddress, fromAddress, tokenId, steps  方向
//	txsRoute, err := nebulaModels.GetNFTHoldRecord(req.ContractAddress, req.Address)
//	if err != nil {
//		return nil, err
//	}
//	return txsRoute, nil
//}

// NFTAddressDetail NFT溯源-地址详情

func (this *NftService) NFTAddressDetail(req apiModels.ReqNFTAddressDetail) (*apiModels.RespNFTAddressDetail, error) {
	var res apiModels.RespNFTAddressDetail
	res.AddressType = 1
	res.LongestHoldTokenId = "3738"
	res.LongestHoldTime = "24.97 天"
	var hold apiModels.HoldTokenId
	var holdHis apiModels.HoldTokenIdHistory
	hold.HoldTokenId = "3738"
	holdHis.HoldTokenIdHistory = "3738"
	res.HoldTokenIds = append(res.HoldTokenIds, hold)
	res.HoldTokenIdsHistory = append(res.HoldTokenIdsHistory, holdHis)
	return &res, nil
}

// todo 改为nebula查询
//func (this *NftService) NFTAddressDetailNew(req apiModels.ReqNFTAddressDetail) (*apiModels.RespNFTAddressDetail, error) {
//
//	//todo 该地址没有所选的NFT交易记录
//	var res apiModels.RespNFTAddressDetail
//
//	tokenIdAndTypes, err := dbModels.HoldTokenIdAndAddressType(req.ContractAddress, req.Address)
//	if err != nil {
//		beego.Error("dbModels.HoldTokenIdAndAddressType error.", err)
//	}
//	if len(tokenIdAndTypes) > 0 {
//		res.AddressType = tokenIdAndTypes[0].AccountType //	地址类型
//		for _, v := range tokenIdAndTypes {
//			var tokenId apiModels.HoldTokenId
//			tokenId.HoldTokenId = v.HoldTokenId
//			res.HoldTokenIds = append(res.HoldTokenIds, tokenId) //	持有的 Token ID
//		}
//	}
//
//	tokenIdsHistory, err := dbModels.HoldTokenIdHistory(req.ContractAddress, req.Address)
//	if err != nil {
//		beego.Error("dbModels.HoldTokenIdHistory error.", err)
//	}
//
//	if len(tokenIdsHistory) == 0 {
//		return nil, errors.New("该地址没有所选的NFT交易记录")
//	}
//	res.HoldTokenIdsHistory = tokenIdsHistory //	历史持有过的 Token ID
//
//	//	最长持有 Token ID
//	//	最长持有时间
//	contractTxs, err := dbModels.LongestHold(req.ContractAddress, req.Address)
//	if err != nil {
//		beego.Error("dbModels.LongestHold error.", err)
//		return nil, err
//	}
//	var startTime time.Time
//	var endTime time.Time
//	var holdTime time.Duration
//	var holdNftNum int
//	var holdNftNumHistory int
//
//	for _, v := range contractTxs {
//		var startTime time.Time
//		var endTime time.Time
//		var holdTime time.Duration
//		if v.To == req.Address {
//			var holdTokenId apiModels.HoldTokenIds
//			holdTokenId.HoldTokenId = v.TokenId
//			holdTokenIds = append(holdTokenIds, holdTokenId)
//			holdTokenIds = append(holdTokenIds[:i], holdTokenIds[i+1:]...)
//			startTime = v.TxTime
//			holdNftNum++
//		}
//		if v.From == req.Address {
//			endTime = v.TxTime
//		}
//		holdTime = endTime.Sub(startTime)
//	}
//
//	if len(tokenIdAndTypes) > 0 {
//		res.AddressType = tokenIdAndTypes[0].AccountType
//	}
//
//	return &res, nil
//}

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

// NFTStartAnalysis NFT溯源-NFT开始分析-交易图
func (this *NftService) NFTStartAnalysis(req apiModels.ReqNFTStartAnalysis) (*nebulaModels.RespGraph, error) {
	// todo token_id为查询条件，返回单条数据
	res, err := nebulaModels.NFTStartAnalysis(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// NFTTrace NFT溯源-NFT追溯-交易图
func (this *NftService) NFTTrace(req apiModels.ReqNFTTrace) (*nebulaModels.RespGraph, error) {
	res, err := nebulaModels.TraceNFTTxs(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
