package nebulaModels

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/models/apiModels"
	"github.com/server/txs-analysis/models/dbModels"
	"github.com/server/txs-analysis/utils"
	"github.com/vesoft-inc/nebula-go/v3/nebula"
	"github.com/zhihu/norm"
	"github.com/zhihu/norm/constants"
	"strconv"
	"strings"
)

//----------数据同步 start-------------------------

func PrepareTxs(db *norm.DB) error {
	err := prepareTxRank(db)
	if err != nil {
		return err
	}
	err1 := prepareCoinTxs(db)
	if err1 != nil {
		return err1
	}
	err2 := prepareTokenTxs(db)
	if err2 != nil {
		return err2
	}
	err3 := prepareNFTTxs(db)
	if err3 != nil {
		return err3
	}
	return nil
}

func prepareTxRank(db *norm.DB) error {
	createSchema :=
		"CREATE EDGE IF NOT EXISTS Tx_Rank(`key` string, tx_rank int64);" +
			"CREATE EDGE INDEX tx_rank_index on Tx_Rank();"
	_, err := db.Execute(createSchema)
	return err
}

func prepareCoinTxs(db *norm.DB) error {
	createSchema :=
		"CREATE EDGE IF NOT EXISTS Coin_Txs(block_id int, tx_hash string, tx_time datetime, from_address string, to_address string, `value` string, amount double, caller string, callee string);" +
			"CREATE EDGE INDEX coin_txs_index on Coin_Txs();"
	_, err := db.Execute(createSchema)
	return err
}

func prepareTokenTxs(db *norm.DB) error {
	createSchema :=
		"CREATE EDGE IF NOT EXISTS Token_Txs(block_id int, tx_hash string, tx_time datetime, from_address string, to_address string, `value` string, amount double, caller string, callee string, token_address string);" +
			"CREATE EDGE INDEX token_txs_index on Token_Txs();"
	_, err := db.Execute(createSchema)
	return err
}

func prepareNFTTxs(db *norm.DB) error {
	createSchema :=
		"CREATE EDGE IF NOT EXISTS NFT_Txs(block_id int, tx_hash string, tx_time datetime, from_address string, to_address string, `value` string, amount double, token_id string, tx_type int, caller string, callee string, token_address string);" +
			"CREATE EDGE INDEX nft_txs_index on NFT_Txs();"
	_, err := db.Execute(createSchema)
	return err
}

func InsertCoinTxs(nebulaDB *norm.DB, tx dbModels.ResSyncTransaction, rank int64) error {
	if len(tx.Caller) == 0 {
		tx.Caller = tx.From
	}
	tx.TxTimeStr = tx.TxTime.Format("2006-01-02 15:04:05")
	rankStr := strconv.Itoa(int(rank))
	err := InsertEdgeTypeGql(nebulaDB, constant.CoinTxs, tx, rankStr)
	if err != nil {
		beego.Error("InsertCoinTxs err: ", err, " txHash:", tx.TxHash)
	} else {
		beego.Info("InsertCoinTxs success! txHash:", tx.TxHash)
	}
	return err
}

func InsertTokenTxs(nebulaDB *norm.DB, tx dbModels.ResSyncTransaction, rank int64) error {
	tx.TxTimeStr = tx.TxTime.Format("2006-01-02 15:04:05")
	rankStr := strconv.Itoa(int(rank))
	err := InsertEdgeTypeGql(nebulaDB, constant.TokenTxs, tx, rankStr)
	if err != nil {
		beego.Error("InsertTokenTxs err: ", err, " txHash:", tx.TxHash)
	} else {
		beego.Info("InsertTokenTxs success! txHash:", tx.TxHash)
	}
	return err
}

func InsertNFTTxs(nebulaDB *norm.DB, tx dbModels.ResSyncTransaction, rank int64) error {
	if strings.HasPrefix(tx.TokenId, "0x") {
		tokenID, err := utils.TokenIDConvert(tx.TokenId)
		if err != nil {
			beego.Error("InsertNFTTxs TokenIDConvert err:", err, " tokenId:", tx.TokenId)
		}
		tx.TokenId = tokenID
	}

	tx.TxTimeStr = tx.TxTime.Format("2006-01-02 15:04:05")
	rankStr := strconv.Itoa(int(rank))
	err := InsertEdgeTypeGql(nebulaDB, constant.NftTxs, tx, rankStr)
	if err != nil {
		beego.Error("InsertNFTTxs err: ", err, " txHash:", tx.TxHash)
	} else {
		beego.Info("InsertNFTTxs success! txHash:", tx.TxHash)
	}
	return err
}

func InsertEdgeTypeGql(db *norm.DB, edgeType string, tx dbModels.ResSyncTransaction, rank string) error {
	nGql := ""
	if len(tx.TokenId) > 0 {
		nGql = fmt.Sprintf("Insert edge %s(block_id,tx_hash,tx_time,from_address,to_address,value,amount,caller,callee,token_id, tx_type,token_address) values '%s' -> '%s'@%s:(%s,'%s',datetime(\"%s\"),'%s','%s','%s',%s,'%s','%s','%s',%s,'%s')",
			edgeType, tx.From, tx.To, rank, tx.BlockId, tx.TxHash, tx.TxTimeStr, tx.From, tx.To, tx.Value, tx.Amount, tx.Caller, tx.ContractAddress, tx.TokenId, tx.TokenType, tx.TokenAddress)
	} else if len(tx.TokenAddress) > 0 {
		nGql = fmt.Sprintf("Insert edge %s(block_id,tx_hash,tx_time,from_address,to_address,value,amount,caller,callee,token_address) values '%s' -> '%s'@%s:(%s,'%s',datetime(\"%s\"),'%s','%s','%s',%s,'%s','%s','%s')",
			edgeType, tx.From, tx.To, rank, tx.BlockId, tx.TxHash, tx.TxTimeStr, tx.From, tx.To, tx.Value, tx.Amount, tx.Caller, tx.ContractAddress, tx.TokenAddress)
	} else {
		nGql = fmt.Sprintf("Insert edge %s(block_id,tx_hash,tx_time,from_address,to_address,value,amount,caller,callee) values '%s' -> '%s'@%s:(%s,'%s',datetime(\"%s\"),'%s','%s','%s',%s,'%s','%s')",
			edgeType, tx.From, tx.To, rank, tx.BlockId, tx.TxHash, tx.TxTimeStr, tx.From, tx.To, tx.Value, tx.Amount, tx.Caller, tx.ContractAddress)
	}
	_, err := db.Debug().Execute(nGql)
	if err != nil {
		return err
	}
	return nil
}

func InsertTxRank(db *norm.DB, from, to string, txRank int64) error {
	key := from + "|" + to
	txRankNebula := TxRank{
		EModel: norm.EModel{
			Src:       from,
			SrcPolicy: constants.PolicyNothing,
			Dst:       to,
			DstPolicy: constants.PolicyNothing,
			//Rank:      rank,
		},
		Key:    key,
		TxRank: txRank,
	}
	err := db.InsertEdge(&txRankNebula)
	if err != nil {
		beego.Error("InsertTxRank err: ", err, " key:", key, " rank:", txRank)
	} else {
		beego.Info("InsertTxRank success! ", " key:", key, " rank:", txRank)
	}
	return err
}

func GetTxsMaxRank(db *norm.DB, edgeType, startBlock, endBlock string) ([]map[string]interface{}, error) {
	//nGql := fmt.Sprintf("LOOKUP ON %s yield src(edge) as src, dst(edge) as dst, rank(edge) AS ranking | GROUP BY $-.src, $-.dst YIELD $-.src as `from`, $-.dst as `to`, MAX($-.ranking) as rank",
	//	edgeType)
	nGql := fmt.Sprintf("LOOKUP ON %s where %s.block_id>=%s and %s.block_id<=%s yield src(edge) as src, dst(edge) as dst, rank(edge) AS ranking | GROUP BY $-.src, $-.dst YIELD $-.src as `from`, $-.dst as `to`, MAX($-.ranking) as rank",
		edgeType, edgeType, startBlock, edgeType, endBlock)
	res, err := db.Debug().Execute(nGql)
	if err != nil {
		return nil, err
	} else {
		result := make([]map[string]interface{}, 0)
		err := UnmarshalResultSet(res, &result)
		if err != nil {
			return result, err
		}
		beego.Info("GetTxsMaxRank success! edgeType:", edgeType, " startBlock:", startBlock, " endBlock:", endBlock)
		return result, nil
	}
}

func QueryTxRank(db *norm.DB, key string) (int64, error) {
	nql := fmt.Sprintf("LOOKUP ON `%s` where %s.key == '%s' YIELD properties(edge).tx_rank as tx_rank",
		constant.TxRanK, constant.TxRanK, key)
	res, err := db.Debug().Execute(nql)
	if err != nil {
		return 0, err
	} else {
		result := make([]map[string]interface{}, 0)
		err := UnmarshalResultSet(res, &result)
		if err != nil {
			return 0, err
		}

		var txRank int64
		for _, vpath := range result {
			for _, v := range vpath {
				if n, ok := v.(int64); ok {
					txRank = n
				}
			}
		}
		return txRank, nil
	}
	return 0, nil
}

func DeleteTxs(edgeType, blockId string) error {
	db := Init()
	nql := fmt.Sprintf("LOOKUP ON `%s` where %s.block_id > %s yield src(edge) as src, dst(edge) as dst, rank(edge) AS ts| delete edge `%s` $-.src->$-.dst@$-.ts", edgeType, edgeType, blockId, edgeType)
	_, err := db.Debug().Execute(nql)
	if err != nil {
		return err
	}
	return nil
}

//----------数据同步 end-------------------------

// QueryTxRoute 查两步以内入账交易
func QueryTxRoute(address string) ([]*TxsRoute, error) {
	db := Init()
	nql := fmt.Sprintf("MATCH p=(v:address)-[e:transaction*1..2]->(v2:address{address:\"%s\"}) RETURN e AS p", address)
	result := make([]map[string]interface{}, 0)
	res, err := db.Debug().Execute(nql)
	if err != nil {
		return nil, err
	} else {
		err := UnmarshalResultSet(res, &result)
		if err != nil {
			return nil, err
		}
		paths := make([]*TxsRoute, 0, len(result))

		for _, vpath := range result {
			for _, v := range vpath {
				if path, ok := v.(*nebula.NList); ok {
					pathValue := path.GetValues()
					steps := ParseTxInfo(pathValue)
					tokenRoute := new(TxsRoute)
					tokenRoute.Steps = steps
					paths = append(paths, tokenRoute)
				}
			}
		}
		return paths, nil
	}
}

func ParseTxInfo(pathValue []*nebula.Value) []RouteTxStep {
	txs := make([]TransactionEdge, 0)
	steps := make([]RouteTxStep, 0)
	for _, value := range pathValue {
		if value.EVal != nil {
			tx := GetTxEdgeInfoFromProps(value.EVal)
			txs = append(txs, tx)
			routeStep := RouteTxStep{
				Transaction: txs,
			}
			steps = append(steps, routeStep)
		}
	}
	return steps
}

func GetTxEdgeInfoFromProps(edge *nebula.Edge) TransactionEdge {
	tx := TransactionEdge{}
	if txHash, exist := edge.Props[constant.TxHash]; exist {
		fmt.Println("txHash:", getValueofValue(txHash))
		tx.TxHash = getValueofValue(txHash)
	}
	if txTime, exist := edge.Props[constant.TxTime]; exist {
		fmt.Println("txTime:", txTime)
		tx.TxTime = getValueofValue(txTime)
	}
	if fromAddress, exist := edge.Props[constant.FromAddress]; exist {
		fmt.Println("fromAddress:", fromAddress)
		tx.FromAddress = getValueofValue(fromAddress)
	}
	if toAddress, exist := edge.Props[constant.ToAddress]; exist {
		fmt.Println("toAddress:", toAddress)
		tx.ToAddress = getValueofValue(toAddress)
	}
	if amount, exist := edge.Props[constant.Amount]; exist {
		fmt.Println("amount:", amount)
		tx.Amount = getValueofValue(amount)
	}
	return tx
}

// GetEdgeTypeTxs 根据边的属性查边和数量
func GetEdgeTypeTxs(db *norm.DB, edgeType string, tx *dbModels.ResSyncTransaction) (int64, error) {
	var nGql string
	if len(tx.ContractAddress) > 0 {
		//nGql = fmt.Sprintf("LOOKUP ON `%s` where %s.tx_hash == \"%s\" and %s.from_address == \"%s\" and %s.to_address==\"%s\" and %s.amount == \"%s\" and %s.contract_address == \"%s\" yield edge as tx, rank(edge) AS rank ", edgeType, edgeType, tx.TxHash, edgeType, tx.From, edgeType, tx.To, edgeType, tx.Amount, edgeType, tx.ContractAddress)
		nGql = fmt.Sprintf("LOOKUP ON `%s` where %s.tx_hash == \"%s\" and %s.from_address == \"%s\" and %s.to_address==\"%s\" and %s.amount == \"%s\" and %s.contract_address == \"%s\" yield edge as tx| YIELD COUNT(*) AS count", edgeType, edgeType, tx.TxHash, edgeType, tx.From, edgeType, tx.To, edgeType, tx.Amount, edgeType, tx.ContractAddress)
	} else {
		//nGql = fmt.Sprintf("LOOKUP ON `%s` where %s.tx_hash == \"%s\" and %s.from_address == \"%s\" and %s.to_address==\"%s\" and %s.amount == \"%s\"  yield edge as tx, rank(edge) AS rank ", edgeType, edgeType, tx.TxHash, edgeType, tx.From, edgeType, tx.To, edgeType, tx.Amount)
		nGql = fmt.Sprintf("LOOKUP ON `%s` where %s.tx_hash == \"%s\" and %s.from_address == \"%s\" and %s.to_address==\"%s\" and %s.amount == \"%s\"  yield edge as tx| YIELD COUNT(*) AS count ", edgeType, edgeType, tx.TxHash, edgeType, tx.From, edgeType, tx.To, edgeType, tx.Amount)
	}
	result := make([]map[string]interface{}, 0)
	res, err := db.Debug().Execute(nGql)
	if err != nil {
		return 0, err
	} else {
		err := UnmarshalResultSet(res, &result)
		if err != nil {
			return 0, err
		}

		var countTx int64
		for _, vpath := range result {
			for _, v := range vpath {
				if n, ok := v.(int64); ok {
					countTx = n
				}

			}
		}
		return countTx, nil
	}
}

// GetEdgeTypeTxsByTxHash 根据边的属性(tx_hash)查边和数量
func GetEdgeTypeTxsByTxHash(db *norm.DB, txHash, edgeType string) (int, error) {
	nGql := fmt.Sprintf("LOOKUP ON `%s` where %s.tx_hash == \"%s\"  yield edge as tx | YIELD COUNT(*) AS count ", edgeType, edgeType, txHash)

	result := make([]map[string]interface{}, 0)
	res, err := db.Debug().Execute(nGql)
	if err != nil {
		return 0, err
	} else {
		err := UnmarshalResultSet(res, &result)
		if err != nil {
			return 0, err
		}
		paths := make([]*TxsRoute, 0, len(result))

		for _, vpath := range result {
			for _, v := range vpath {
				if path, ok := v.(*nebula.NList); ok {
					pathValue := path.GetValues()
					steps := ParseTxInfo(pathValue)
					tokenRoute := new(TxsRoute)
					tokenRoute.Steps = steps
					paths = append(paths, tokenRoute)
				}
			}
		}
		return 0, nil
	}
}

//---------------------------- 交易图谱查询---------------------------------

// dbExecute 执行方法
func dbExecute(nGql string) ([]map[string]interface{}, error) {
	db := Init()
	result := make([]map[string]interface{}, 0)
	res, err := db.Debug().Execute(nGql)
	if err != nil {
		return nil, err
	} else {
		err := UnmarshalResultSet(res, &result)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

// AddressTxAnalysis 地址分析-地址交易图
func AddressTxAnalysis(req apiModels.ReqAddressTxGraph) (*RespGraph, error) {
	//todo 支持查转入 转出 全部
	tagConditionV := ""
	tagConditionV2 := ""
	eCondition := ""
	eCondition = constant.CoinTxs
	if len(req.ContractAddress) > 0 {
		eCondition = constant.TokenTxs + "{token_address:'" + req.ContractAddress + "',"
	}
	if req.Direction == "out" {
		tagConditionV = "{address:'" + req.Address + "'}"
	}
	if req.Direction == "in" {
		tagConditionV2 = "{address:'" + req.Address + "'}"
	}

	nGQL := fmt.Sprintf("MATCH (v:address%s)-[e:%s]->(v2:address%s) RETURN v,e,v2, sum(properties(e).amount) AS total_amount, count(*) as tx_count order by total_amount desc limit %s ", tagConditionV, eCondition, tagConditionV2, req.Count)

	result, err := dbExecute(nGQL)
	if err != nil || len(result) == 0 {
		return nil, err
	}
	res := ParseNebulaResult(result)
	return &res, nil
}

// NFTStartAnalysis NFT溯源-开始分析
func NFTStartAnalysis(req apiModels.ReqNFTStartAnalysis) (*RespGraph, error) {
	nGQL := fmt.Sprintf("MATCH p=(v)-[e:NFT_Txs]->(v2) where e.token_address==\"%s\"", req.ContractAddress)
	if req.Field == "tx_hash" {
		nGQL += fmt.Sprintf(" and e.tx_hash==\"%s\" and e.token_id== \"%s\" RETURN v, e ,v2 limit 1", req.Value, req.TokenId)
		//nGQL =fmt.Sprintf("LOOKUP ON NFT_Txs WHERE NFT_Txs.tx_hash == \"%s\" YIELD properties(edge).from_address AS from_address, properties(edge).to_address AS to_address,properties(edge).tx_hash as tx_hash", input)
	}

	if req.Field == "address" {
		//nGQL = fmt.Sprintf("MATCH (v)-[e:NFT_Txs]->(v2),(v3)-[e2:NFT_Txs]->(v4) where e.token_address==\"%s\" and e.to_address==\"%s\" and e2.token_address==\"%s\" and e2.from_address==\"%s\" e.token_id== \"%s\" RETURN v,e,e2,v4 limit 1", req.ContractAddress, req.Value, req.ContractAddress, req.Value, req.TokenId)
		nGQL = fmt.Sprintf("MATCH (v)-[e:NFT_Txs{token_address:'%s',to_address:'%s',token_id:'%s'}]->(v2),(v3)-[e2:NFT_Txs{token_address:'%s',from_address:'%s',token_id:'%s'}]->(v4)  RETURN v,e,e2,v4 limit 1", req.ContractAddress, req.Value, req.TokenId, req.ContractAddress, req.Value, req.TokenId)
	}

	if req.Field == "token_id" {
		//nGQL =fmt.Sprintf("LOOKUP ON NFT_Txs WHERE NFT_Txs.token_id == \"%s\" YIELD properties(edge).from_address AS from_address, properties(edge).to_address AS to_address,properties(edge).tx_hash as tx_hash", input)
		nGQL += fmt.Sprintf(" and e.from_address==\"0x0000000000000000000000000000000000000000\" and e.token_id==\"%s\"", req.Value)
		nGQL += " RETURN v,e,v2 limit 1"
	}

	result, err := dbExecute(nGQL)
	if err != nil || len(result) == 0 {
		return nil, err
	}
	res := ParseNebulaResult(result)
	return &res, nil
}

// TraceNFTTxs NFT溯源-追溯交易
func TraceNFTTxs(req apiModels.ReqNFTTrace) (*RespGraph, error) {
	//nGQL := fmt.Sprintf("MATCH p=(v:address)-[e:transaction*1..2]->(v2:address{address:\"%s\"}) RETURN e AS p", req.Address)
	//nGQL := fmt.Sprintf("GET SUBGRAPH WITH PROP 1 STEPS FROM \"0x18e548550a81e318b5b4ac97e26ed1958c8f12e4\" OUT NFT_Txs where NFT_Txs.token_id == \"15\" YIELD VERTICES as address,EDGES as e")
	nGQL := fmt.Sprintf("GET SUBGRAPH WITH PROP %s STEPS FROM \"%s\" %s NFT_Txs where NFT_Txs.token_address == \"%s\" and NFT_Txs.token_id == \"%s\" YIELD VERTICES as v, EDGES as e", req.Count, req.Address, req.Direction, req.ContractAddress, req.TokenID)
	result, err := dbExecute(nGQL)
	if err != nil || len(result) == 0 {
		return nil, err
	}
	res := ParseNebulaResult(result)
	return &res, nil
}

// GetNFTHoldRecord 根据地址或者token Id获取NFT的持有记录
func GetNFTHoldRecord(contractAddress, input string) (*RespGraph, error) {
	nGQL := fmt.Sprintf("MATCH p=(v)-[e:NFT_Txs]->(v2) where e.token_address==\"%s\"", contractAddress)
	if strings.HasPrefix(input, "0x") {
		//nGQL = fmt.Sprintf("MATCH (v)-[e:NFT_Txs]->(v2),(v3)-[e2:NFT_Txs]->(v4) where e.token_address==\"%s\" and e.to_address==\"%s\" and e2.token_address==\"%s\" and e2.from_address==\"%s\" RETURN v,e,e2,v4 limit 1", contractAddress, input, contractAddress, input)
		//nGQL = fmt.Sprintf("GET SUBGRAPH WITH PROP 1 STEPS FROM \"%s\" BOTH NFT_Txs where NFT_Txs.token_address == \"%s\" YIELD VERTICES as address, EDGES as e", input, contractAddress)
		nGQL = fmt.Sprintf("GET SUBGRAPH WITH PROP 1 STEPS FROM \"%s\" BOTH NFT_Txs where NFT_Txs.token_address == \"%s\" YIELD EDGES as e", input, contractAddress)

	} else {
		//todo 查询条件为token id
		//nGQL =fmt.Sprintf("LOOKUP ON NFT_Txs WHERE NFT_Txs.token_id == \"%s\" YIELD properties(edge).from_address AS from_address, properties(edge).to_address AS to_address,properties(edge).tx_hash as tx_hash", input)
		nGQL += fmt.Sprintf(" and e.from_address==\"0x0000000000000000000000000000000000000000\" and e.token_id==\"%s\"", input)
		nGQL += " RETURN v,e,v2 limit 1"
	}

	result, err := dbExecute(nGQL)
	if err != nil || len(result) == 0 {
		return nil, err
	}
	res := ParseNebulaResult(result)
	return &res, nil
}
