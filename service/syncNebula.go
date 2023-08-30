package service

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/server/txs-analysis/constant"
	"github.com/server/txs-analysis/models/dbModels"
	"github.com/server/txs-analysis/models/nebulaModels"
	"github.com/server/txs-analysis/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	eachSyncNum        int
	syncDuringTime     int
	intervalBlockNum   int
	nebulaBlockControl int
)

var txsRankMap = make(map[string]int64)

func init() {
	syncNebulaFlag, _ := beego.AppConfig.Int("nebula::syncflag")
	syncNebulaTime, _ := beego.AppConfig.Int("nebula::synctime")
	eachSyncNum, _ = beego.AppConfig.Int("nebula::apieachsyncnum")
	syncDuringTime, _ = beego.AppConfig.Int("nebula::apisyncduringtime")
	intervalBlockNum, _ = beego.AppConfig.Int("nebula::intervalblocknum")

	if syncNebulaFlag == 1 {
		beego.Info("syncNebula start --> syncNebulaFlag is ", syncNebulaFlag)

		err := ReadTxsRank()
		if err != nil {
			beego.Error("ReadTxsRank error: ", err)
			return
		}
		fmt.Println(nebulaBlockControl)

		timer := time.NewTimer(time.Second * time.Duration(2))
		go func(t *time.Timer) {
			for range t.C {
				SyncDataToNebula()
				t.Reset(time.Second * time.Duration(syncNebulaTime))
			}
		}(timer)
	}
}

func ReadTxsRank() error {
	beego.Debug("ReadTxsRank start ==>")

	files, err := filepath.Glob(constant.DirPath + "*" + constant.RankSuffix)
	if len(files) == 0 {
		beego.Debug("ReadTxsRank count = 0")
		return nil
	}

	blockIdMax := 0
	for _, file := range files {
		blockId, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(file, constant.RankSuffix), constant.DirPath))
		if blockId > blockIdMax {
			blockIdMax = blockId
		}
	}
	nebulaBlockControl = blockIdMax
	fileName := constant.DirPath + strconv.Itoa(blockIdMax) + constant.RankSuffix

	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		beego.Error("ReadTxsRank open file failed. err:", err)
		return err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	i := 0
	mapLen := 0
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		str = strings.TrimRight(str, "\n")
		index := strings.Index(str, ":")
		key := str[:index]
		value := str[index+1:]
		valueInt, _ := strconv.ParseInt(value, 10, 64)

		if i == 1 {
			mapLen = int(valueInt)
		}
		if i > 1 {
			txsRankMap[key] = valueInt
		}
		i++
	}

	if len(txsRankMap) != mapLen {
		return errors.New("The number of data read from the file is incorrect.")
	}

	beego.Info("-----------ReadTxsRank end. fileName:", fileName, " count:", mapLen)
	return nil
}

func SyncDataToNebula() {
	ormer := orm.NewOrm()
	var blockHeight int
	//sql := "select block_id from tb_transaction order by block_id desc limit 1"
	sql := "SELECT last_analyze_block_id from tb_block_control WHERE control_type = 1 and `status` = 1 and is_deleted = 0"
	err := ormer.Raw(sql).QueryRow(&blockHeight)
	if err != nil {
		beego.Error("get blockHeight error", err)
		return
	}

	if blockHeight-intervalBlockNum <= nebulaBlockControl {
		return
	}

	startNum := nebulaBlockControl + 1
	endNum := blockHeight - intervalBlockNum

	for currentBlockId := startNum; currentBlockId <= endNum; currentBlockId += eachSyncNum {
		currentEndBlockId := currentBlockId + eachSyncNum - 1
		if currentEndBlockId > endNum {
			currentEndBlockId = endNum
		}

		start := strconv.Itoa(currentBlockId)
		end := strconv.Itoa(currentEndBlockId)

		errAddr := SyncAddressDataToNebula(start, end)
		if errAddr != nil {
			beego.Error("SyncAddressDataToNebula error", errAddr)
			return
		}

		errCoin := SyncCoinTxsDataToNebula(start, end)
		if errCoin != nil {
			beego.Error("SyncCoinTxsDataToNebula error", errCoin)
			return
		}

		errToken := SyncTokenTxsDataToNebula(start, end)
		if errToken != nil {
			beego.Error("SyncTokenTxsDataToNebula error", errToken)
			return
		}

		errNFT := SyncNFTTxsDataToNebula(start, end)
		if errNFT != nil {
			beego.Error("SyncNFTTxsDataToNebula error", errNFT)
			return
		}

		errUpdate := SaveBlockIdAndRank(strconv.Itoa(currentEndBlockId))
		if errUpdate != nil {
			beego.Error("SyncDataToNebula updateBlockControl error", errUpdate)
			return
		}
	}
	return
}

func SaveBlockIdAndRank(blockId string) error {
	fileName := constant.DirPath + blockId + constant.RankSuffix
	exist, err := utils.PathFileExists(fileName)
	if exist {
		err := os.Remove(fileName)
		if err != nil {
			beego.Error("SaveBlockIdAndRank os.Remove error. fileName:", fileName)
			return err
		}
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	write := bufio.NewWriter(f)

	write.WriteString("blockId:" + blockId + "\n")
	write.WriteString("length(txsRankMap):" + strconv.Itoa(len(txsRankMap)) + "\n")

	for k, v := range txsRankMap {
		str := k + ":" + strconv.FormatInt(v, 10)
		write.WriteString(str + "\n")
	}

	err = write.Flush()
	if err != nil {
		return err
	}
	defer f.Close()
	beego.Debug("SaveBlockIdAndRank fileName:", fileName)

	//保留3个文件，其余的删除
	files, err := filepath.Glob(constant.DirPath + "*" + constant.RankSuffix)
	if err != nil {
		beego.Error("SaveBlockIdAndRank filepath.Glob err:", err)
		return err
	}
	if len(files) <= 3 {
		return nil
	}

	blockIdInt, _ := strconv.Atoi(blockId)
	blockIdMin := blockIdInt
	for _, file := range files {
		blockIdOld, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(file, constant.DirPath), constant.RankSuffix))
		if blockIdOld < blockIdInt {
			if blockIdMin > blockIdOld {
				blockIdMin = blockIdOld
			}
		}
	}

	latestFile := constant.DirPath + strconv.Itoa(blockIdMin) + constant.RankSuffix
	err = os.Remove(latestFile)
	if err != nil {
		beego.Error("SaveBlockIdAndRank os.Remove error. fileName:", latestFile)
		return err
	}
	return nil
}

func SyncAddressDataToNebula(startNum, endNum string) error {
	beego.Debug("SyncAddressDataToNebula==> start:", startNum, " end:", endNum)
	nebulaDB := nebulaModels.Init()

	//同步address
	accountList, err := dbModels.GetSyncAddressData(startNum, endNum)
	if err != nil {
		beego.Error("SyncDataToNebula dbModels.GetAccountType error: ", err)
		return err
	}

	if len(accountList) > 0 {
		for _, v := range accountList {
			address := &nebulaModels.Address{
				Address: v.AccountAddress,
				Type:    v.AccountType,
			}
			err := nebulaModels.InsertAddress(nebulaDB, address)
			if err != nil {
				beego.Error("SyncDataToNebula nebulaModels.InsertAddress error: ", err, " address:", v.AccountAddress)
				if strings.Contains(err.Error(), "RPC failure, probably timeout") {
					nebulaDB := nebulaModels.Init()
					for {
						beego.Debug("insert address Contains----timeout-------->")
						errNew := nebulaModels.InsertAddress(nebulaDB, address)
						if errNew == nil {
							break
						}
					}
					beego.Debug("insert address Contains----timeout--break------>")
				} else {
					return err
				}
			}
		}
		time.Sleep(time.Duration(syncDuringTime) * time.Second)
	}

	nebulaDB.Close()
	return nil
}

func SyncCoinTxsDataToNebula(startNum, endNum string) error {
	beego.Debug("SyncCoinTxsDataToNebula==> start:", startNum, " end:", endNum)
	nebulaDB := nebulaModels.Init()

	//同步 tb_transaction 表数据
	txList, err := dbModels.GetSyncTxData(startNum, endNum)
	if err != nil {
		beego.Error("SyncCoinTxsDataToNebula dbModels.GetSyncTxData error: ", err)
		return err
	}

	if len(txList) > 0 {
		for _, v := range txList {
			//获取rank
			key := constant.CoinTxs + v.From + v.To
			coinRank, ok := txsRankMap[key]
			if !ok {
				coinRank = 0
				txsRankMap[key] = 0
			} else {
				coinRank = coinRank + 1
				txsRankMap[key] = coinRank
			}
			err := nebulaModels.InsertCoinTxs(nebulaDB, *v, coinRank)
			if err != nil {
				beego.Error("SyncCoinTxsDataToNebula insert txs error: ", err)
				if strings.Contains(err.Error(), "RPC failure, probably timeout") {
					nebulaDB := nebulaModels.Init()
					for {
						beego.Debug("SyncCoinTxsDataToNebula insert tx Contains----timeout-------->")
						errNew := nebulaModels.InsertCoinTxs(nebulaDB, *v, coinRank+1)
						if errNew == nil {
							break
						}
					}
					beego.Debug("SyncCoinTxsDataToNebula insert tx Contains----timeout--break------>")
				} else {
					return err
				}
			}
		}
	}
	time.Sleep(time.Duration(syncDuringTime) * time.Second)

	//同步 tb_internal_transaction 表数据
	txInternalList, err := dbModels.GetSyncInternalTxData(startNum, endNum)
	if err != nil {
		beego.Error("SyncDataToNebula dbModels.GetSyncTxData error: ", err)
		return err
	}

	if len(txInternalList) > 0 {
		for _, v := range txInternalList {
			key := constant.CoinTxs + v.From + v.To
			coinRank, ok := txsRankMap[key]
			if !ok {
				coinRank = 0
				txsRankMap[key] = 0
			} else {
				coinRank = coinRank + 1
				txsRankMap[key] = coinRank
			}

			err := nebulaModels.InsertCoinTxs(nebulaDB, *v, coinRank)
			if err != nil {
				beego.Error("SyncDataToNebula insert internal txs error: ", err)
				if strings.Contains(err.Error(), "RPC failure, probably timeout") {
					nebulaDB := nebulaModels.Init()
					for {
						beego.Debug("SyncDataToNebula insert internal txs Contains----timeout-------->")
						errNew := nebulaModels.InsertCoinTxs(nebulaDB, *v, coinRank)
						if errNew == nil {
							break
						}
					}
					beego.Debug("SyncDataToNebula insert internal txs Contains----timeout--break------>")
				} else {
					return err
				}
			}
		}
	}
	nebulaDB.Close()
	return nil
}

func SyncTokenTxsDataToNebula(startNum, endNum string) error {
	beego.Debug("SyncTokenTxsDataToNebula==> start:", startNum, " end:", endNum)
	nebulaDB := nebulaModels.Init()
	//同步token transaction
	txList, err := dbModels.GetSyncTokenTxsData(startNum, endNum)
	if err != nil {
		beego.Error("SyncDataToNebula dbModels.GetSyncTxData error: ", err)
		return err
	}

	if len(txList) > 0 {
		for _, v := range txList {
			key := constant.TokenTxs + v.From + v.To
			tokenRank, ok := txsRankMap[key]
			if !ok {
				tokenRank = 0
				txsRankMap[key] = 0
			} else {
				tokenRank = tokenRank + 1
				txsRankMap[key] = tokenRank
			}

			err := nebulaModels.InsertTokenTxs(nebulaDB, *v, tokenRank)
			if err != nil {
				beego.Error("SyncDataToNebula nebulaModels.InsertTxn error: ", err)
				if strings.Contains(err.Error(), "RPC failure, probably timeout") {
					nebulaDB := nebulaModels.Init()
					for {
						beego.Debug("insert tx Contains----timeout-------->")
						errNew := nebulaModels.InsertTokenTxs(nebulaDB, *v, tokenRank)
						if errNew == nil {
							break
						}
					}
					beego.Debug("insert tx Contains----timeout--break------>")
				} else {
					return err
				}
			}
		}
		time.Sleep(time.Duration(syncDuringTime) * time.Second)
	}
	nebulaDB.Close()
	return nil
}

func SyncNFTTxsDataToNebula(startNum, endNum string) error {
	beego.Debug("SyncNFTTxsDataToNebula==> start:", startNum, " end:", endNum)
	nebulaDB := nebulaModels.Init()
	//同步NFT交易
	txList, err := dbModels.GetSyncNFTTxsData(startNum, endNum)
	if err != nil {
		beego.Error("SyncDataToNebula dbModels.GetSyncTxData error: ", err)
		return err
	}

	if len(txList) > 0 {
		for _, v := range txList {
			key := constant.NftTxs + v.From + v.To
			nftRank, ok := txsRankMap[key]
			if !ok {
				nftRank = 0
				txsRankMap[key] = 0
			} else {
				nftRank = nftRank + 1
				txsRankMap[key] = nftRank
			}

			err := nebulaModels.InsertNFTTxs(nebulaDB, *v, nftRank)
			if err != nil {
				beego.Error("SyncDataToNebula nebulaModels.InsertTxn error: ", err)
				if strings.Contains(err.Error(), "RPC failure, probably timeout") {
					nebulaDB := nebulaModels.Init()
					for {
						beego.Debug("insert tx Contains----timeout-------->")
						errNew := nebulaModels.InsertNFTTxs(nebulaDB, *v, nftRank)
						if errNew == nil {
							break
						}
					}
					beego.Debug("insert tx Contains----timeout--break------>")
				} else {
					return err
				}
			}
		}
		time.Sleep(time.Duration(syncDuringTime) * time.Second)
	}
	nebulaDB.Close()
	return nil
}
