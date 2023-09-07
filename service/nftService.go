package service

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
	"github.com/server/txs-analysis/models/nebulaModels"
	"github.com/server/txs-analysis/utils"
	"math/big"
	"sort"
	"time"
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
func (this *NftService) NFTAddressDetail(req apiModels.ReqNFTAddressDetail) (*apiModels.RespNFTAddressDetail, error) {
	var res apiModels.RespNFTAddressDetail

	tokenIdAndTypes, err := dbModels.HoldTokenIdAndAddressType(req.ContractAddress, req.Address)
	if err != nil {
		beego.Error("dbModels.HoldTokenIdAndAddressType error.", err)
	}
	if len(tokenIdAndTypes) > 0 {
		res.AddressType = tokenIdAndTypes[0].AccountType //	地址类型
		for _, v := range tokenIdAndTypes {
			var tokenId apiModels.HoldTokenId
			tokenId.TokenId = common.HexToHash(v.TokenId).Big().String()
			res.HoldTokenIds = append(res.HoldTokenIds, tokenId) //	持有的 Token ID
		}
	}

	tokenIdsHistory, err := dbModels.HoldTokenIdHistory(req.ContractAddress, req.Address)
	if err != nil {
		beego.Error("dbModels.HoldTokenIdHistory error.", err)
	}

	if len(tokenIdsHistory) == 0 {
		return nil, errors.New("该地址没有所选的NFT交易记录")
	}
	for _, v := range tokenIdsHistory {
		var tokenIdHis apiModels.HoldTokenIdHistory
		tokenIdHis.HoldTokenIdHistory = common.HexToHash(v.HoldTokenIdHistory).Big().String()
		res.HoldTokenIdsHistory = append(res.HoldTokenIdsHistory, tokenIdHis) // 历史持有过的 Token ID
	}

	//	最长持有 Token ID
	//	最长持有时间
	contractTxs, err := dbModels.LongestHold(req.ContractAddress, req.Address, "")
	if err != nil {
		beego.Error("dbModels.LongestHold error.", err)
		return nil, err
	}

	maxHoldTime, maxHoldToken := getLongestHold(contractTxs, req.Address)
	res.LongestHoldTokenId = maxHoldToken
	// 处理最长持有时间
	maxHoldTimeBigInt := big.NewInt(maxHoldTime)
	maxHoldDay := maxHoldTimeBigInt.Div(maxHoldTimeBigInt, big.NewInt(86400))
	res.LongestHoldTime = maxHoldDay.String() + "天"

	return &res, nil
}

func getLongestHold(contractTxs []*dbModels.TbContractTransaction, accountAddress string) (int64, string) {
	mainStartMap := map[string][]map[int]int64{}
	mainEndMap := map[string][]map[int]int64{}

	for i := 0; i < len(contractTxs); i++ {
		subStartMap := make(map[int]int64)
		subEndMap := make(map[int]int64)

		if contractTxs[i].To == accountAddress {
			tokenId := common.HexToHash(contractTxs[i].TokenId).Big().String()
			subStartMap[len(mainStartMap[tokenId])] = contractTxs[i].TxTime.Unix()
			//subStartMaps = append(subStartMaps, subStartMap)
			//mainStartMap[tokenId] = subStartMaps
			mainStartMap[tokenId] = append(mainStartMap[tokenId], subStartMap)
		}
		if contractTxs[i].From == accountAddress {
			tokenId := common.HexToHash(contractTxs[i].TokenId).Big().String()
			subEndMap[len(mainEndMap[tokenId])] = contractTxs[i].TxTime.Unix()
			//subEndMaps = append(subEndMaps, subEndMap)
			mainEndMap[tokenId] = append(mainEndMap[tokenId], subEndMap)
		}
	}

	//mainStartMap := map[int]map[string]int64{}
	//mainEndMap := map[int]map[string]int64{}
	//
	//for i := 0; i < len(contractTxs); i++ {
	//	subStartMap := make(map[string]int64)
	//	subEndMap := make(map[string]int64)
	//
	//	if contractTxs[i].To == req.Address {
	//		tokenId := common.HexToHash(contractTxs[i].TokenId).Big().String()
	//		subStartMap[tokenId] = contractTxs[i].TxTime.Unix()
	//		mainStartMap[len(mainStartMap)] = subStartMap
	//	}
	//	if contractTxs[i].From == req.Address {
	//		tokenId := common.HexToHash(contractTxs[i].TokenId).Big().String()
	//		subEndMap[tokenId] = contractTxs[i].TxTime.Unix()
	//		mainEndMap[len(mainStartMap)] = subEndMap
	//	}
	//}

	holdMap := make(map[int64]string)
	for k, v := range mainStartMap {
		var startTime int64
		var endTime int64
		var holdTime int64
		if _, ok := mainEndMap[k]; !ok {
			//mainEndMap不存在当前tokenId记录，则此tokenId为当前地址持有
			endTime = time.Now().Unix()
			startTime = v[0][0]
			holdTime = endTime - startTime
			holdMap[holdTime] = k
		} else {
			// 遍历数组
			for i := 0; i < len(v); i++ { //v []map[int]int64
				// 遍历map
				for key, value := range v[i] {
					if _, ok := mainEndMap[k][i][key]; !ok {
						endTime = time.Now().Unix()
						startTime = value
						holdTime = endTime - startTime
						holdMap[holdTime] = k
					} else {
						endTime = mainEndMap[k][i][key]
						startTime = value
						holdTime = endTime - startTime
						holdMap[holdTime] = k
					}
				}
			}
		}
	}

	//遍历holdMap, 计算最长持有时间
	var maxHoldTime int64
	var maxHoldToken string
	for k, v := range holdMap {

		if k > maxHoldTime {
			maxHoldTime = k
			maxHoldToken = v
		}
	}
	return maxHoldTime, maxHoldToken
}

// NFTDetail NFT溯源-NFT详情
func (this *NftService) NFTDetail(req apiModels.ReqNFTDetail) (*apiModels.RespNFTDetail, error) {
	res := apiModels.RespNFTDetail{}
	detail, err := dbModels.NFTDetail(req)
	if err != nil {
		beego.Error("dbModels.NFTDetail error.", err)
	}
	if detail != nil {
		res.TransferCount = detail.TransferCount           // 流转次数
		res.HistoryHolderCount = detail.HistoryHolderCount // 历史持有者数量
		res.MintTime = detail.MintTime                     // 铸造时间
	}

	contract, err := dbModels.ContractInfo(req)
	if err != nil {
		beego.Error("dbModels.ContractInfo error.", err)
	}

	if contract != nil {
		res.Name = contract.Name
		res.Symbol = contract.Symbol
		res.Logo = contract.Logo
		res.TokenType = contract.TokenType
		res.Holder = contract.Holder
	}

	res.ContractAddress = req.ContractAddress

	//持有者、当前开始持有时间
	contractAccount, err := dbModels.HoldTokenInfo(req.ContractAddress, req.TokenID)
	if err != nil {
		beego.Error("dbModels.HoldTokenInfo error.", err)
	}
	if contractAccount != nil {
		res.Holder = contractAccount.To // 持有者
		// 处理当前持有时间
		res.CurrentHoldTime = utils.Float64String(time.Since(contractAccount.TxTime).Hours() / 24)
	}

	contractTxs, err := dbModels.LongestHold(req.ContractAddress, "", req.TokenID)
	if err != nil {
		beego.Error("NFTDetail dbModels.LongestHold error.", err)
		//return nil, err
	}

	if len(contractTxs) > 0 {
		//计算最长持有时间
		var holdTimes []int64
		for i := 0; i < len(contractTxs); i++ {
			var holdTime int64
			if len(contractTxs) == 1 {
				holdTime = time.Now().Unix() - contractTxs[i].TxTime.Unix()
			}else {
				if i == 0 {
					continue
				}
				if i == len(contractTxs)-1 {
					holdTime = time.Now().Unix() - contractTxs[i].TxTime.Unix()
				} else {
					holdTime = contractTxs[i].TxTime.Unix() - contractTxs[i-1].TxTime.Unix()
				}
			}
			holdTimes = append(holdTimes, holdTime)
		}

		//遍历holdMap, 计算最长持有时间
		var maxHoldTime int64
		for _, v := range holdTimes {
			if v > maxHoldTime {
				maxHoldTime = v
			}
		}

		// 处理最长持有时间
		maxHoldTimeBigInt := big.NewInt(maxHoldTime)
		maxHoldDay := maxHoldTimeBigInt.Div(maxHoldTimeBigInt, big.NewInt(86400))
		res.LongestHoldTime = maxHoldDay.String() + "天"
	}

	return &res, nil
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
