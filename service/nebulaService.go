package service

import (
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/nebulaModels"
)

type NebulaService struct {
}

//GetEntryTxsByAddress 根据指定地址，查询所有入账记录
func (this *NebulaService) GetEntryTxsByAddress(address string) ([]*nebulaModels.TxsRoute, error) {
	nebulaDB := nebulaModels.Init()
	txsRoute, err := nebulaModels.QueryTxRoute(nebulaDB, address)
	if err != nil {
		return nil, err
	}
	return txsRoute, nil
}

//NFTTrace NFT溯源-NFT追溯-交易图
func (this *NebulaService) TraceNFTTxs(req apiModels.ReqNFTTrace) ([]*nebulaModels.TxsRoute, error) {
	txsRoute, err := nebulaModels.TraceNFTTxs(req.Address)
	if err != nil {
		return nil, err
	}
	return txsRoute, nil
}

//根据指定条件(tokenID address txHash)，查询所有入账记录
func (this *NebulaService) NFTStartAnalysis(contractAddress, input string) ([]*nebulaModels.TxsRoute, error) {

	//todo 怎么确定是否是最后一个或第一个
	//from= 0x0000000000000000000000000000000000000000 是开始， to 是owner
	txsRoute, err := nebulaModels.QueryNFTTxsPath(contractAddress, input)
	if err != nil {
		return nil, err
	}
	return txsRoute, nil
}


//根据指定条件(tokenID address txHash)，查询所有入账记录
func (this *NebulaService) GetNFTTxsPath(req apiModels.ReqNFTTrace) ([]*nebulaModels.TxsRoute, error) {
	nebulaDB := nebulaModels.Init()

	//todo 怎么确定是否是最后一个或第一个
	//from= 0x0000000000000000000000000000000000000000 是开始， to 是owner
	//contractAddress, fromAddress, tokenId, steps  方向
	txsRoute, err := nebulaModels.GetNFTTxsPath(nebulaDB, req.ContractAddress, req.Address, req.TokenID, req.Count, req.Direction)
	if err != nil {
		return nil, err
	}
	return txsRoute, nil
}
