package dbModels

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/server/txs-analysis/models/apiModels"
	"math/big"
	"time"

	"github.com/astaxie/beego/orm"
)

type TbContract struct {
	Id                  int64     `orm:"column(id);pk"`
	BlockId             int64     `orm:"column(block_id);null" description:"区块号"`
	BlockHash           string    `orm:"column(block_hash);size(255);null" description:"区块hash"`
	TxHash              string    `orm:"column(tx_hash);size(255);null" description:"交易hash"`
	ContractCreator     string    `orm:"column(contract_creator);size(255);null" description:"创建者"`
	ContractAddress     string    `orm:"column(contract_address);size(255);null" description:"合约地址"`
	ContractType        int       `orm:"column(contract_type);null" description:"合约类型：1-erc20 2-erc721"`
	Name                string    `orm:"column(name);size(255);null" description:"合约名称"`
	Symbol              string    `orm:"column(symbol);size(255);null" description:"合约标识"`
	Logo                string    `orm:"column(logo);size(255);null" description:"合约logo"`
	Website             string    `orm:"column(website);size(255);null" description:"网站"`
	SocialProfiles      string    `orm:"column(social_profiles);size(255);null" description:"社交档案"`
	Decimals            int64     `orm:"column(decimals);null" description:"合约精度"`
	TotalSupply         string    `orm:"column(total_supply);size(255);null" description:"总量"`
	ContractInfo        string    `orm:"column(contract_info);null" description:"合约信息"`
	ContractCode        string    `orm:"column(contract_code);null" description:"字节码"`
	Abi                 string    `orm:"column(abi);null" description:"abi"`
	SourceCode          string    `orm:"column(source_code);null" description:"源码"` // 23/6/26 update IsJsonFormatCode为1时,使用 [[filename,sourcecode], [filename2: sourcecode2]] json格式保存
	IsJsonFormatCode    int       `orm:"column(is_json_format_code);null" description:"是否是json格式的源码"`
	Settings            string    `orm:"column(settings);null" description:"standard-json设置项"`
	Arguments           string    `orm:"column(arguments);null" description:"构造函数参数"`
	LicenseType         string    `orm:"column(license_type);size(255);null" description:"license类型"`
	CompilerType        string    `orm:"column(compiler_type);size(255);null" description:"编译类型"`
	CompilerVersion     string    `orm:"column(compiler_version);size(255);null" description:"编译器版本"`
	OptimizationEnabled int8      `orm:"column(optimization_enabled);null" description:"编译优化开启状态 0-未开启 1-开启"`
	OptimizationRuns    int       `orm:"column(optimization_runs);null" description:"编译优化时间"`
	EvmVersion          string    `orm:"column(evm_version);size(255);null" description:"EVM版本"`
	Status              int       `orm:"column(status);null" description:"状态,1-正常 2-冻结"`
	IsVerfied           int8      `orm:"column(is_verfied);null" description:"验证状态 0-未验证 1-验证通过 2-验证失败"`
	Is165interface      int8      `orm:"column(is_165interface);null" description:"165接口状态 0-否 1-是"`
	IsMetadata          int8      `orm:"column(is_metadata);null" description:"支持Metadata 0-否 1-是"`
	IsEnumerable        int8      `orm:"column(is_enumerable);null" description:"支持Enumerable 0-否 1-是"`
	IsImplements        int8      `orm:"column(is_implements);null" description:"是否实现对应的合约标准 0-未实现 1-实现"`
	VerfiedTime         time.Time `orm:"column(verfied_time);type(timestamp);null;" description:"验证时间"`
	IsDeleted           int8      `orm:"column(is_deleted);null" description:"删除状态 0-正常 1-删除"`
	SyncTime            time.Time `orm:"column(sync_time);type(timestamp);null;auto_now_add" description:"同步时间"`
	CreateTime          time.Time `orm:"column(create_time);type(timestamp);null;auto_now_add" description:"创建时间"`
	UpdateTime          time.Time `orm:"column(update_time);type(timestamp);null;auto_now_add" description:"更新时间"`
}

func (t *TbContract) TableName() string {
	return "tb_contract"
}

func init() {
	orm.RegisterModel(new(TbContract))
}

// QueryContractList 获取可查合约列表
func QueryContractList(txType, name string) ([]*apiModels.RespContractList, error) {
	list := make([]*apiModels.RespContractList, 0)
	orm := orm.NewOrm()

	condition := " where status = 1 and contract_type = " + txType + " and symbol != '' AND ( LOWER(name) LIKE LOWER('%" + name + "%') OR LOWER(symbol) LIKE LOWER('%" + name + "%')) ORDER BY block_id DESC"
	_, err := orm.Raw("SELECT `name`, symbol, logo, contract_address FROM " + new(TbContract).TableName() + condition).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// NFTList 获取NFT列表
func NFTList() ([]*apiModels.RespNftList, error) {
	list := make([]*apiModels.RespNftList, 0)
	orm := orm.NewOrm()
	condition := " where status = 1 and contract_type = 2 and symbol != '' ORDER BY block_id DESC"
	_, err := orm.Raw("SELECT symbol, logo, contract_address FROM " + new(TbContract).TableName() + condition).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func ContractInfo(req apiModels.ReqNFTDetail) (*apiModels.RespNFTDetail, error) {
	var res *apiModels.RespNFTDetail
	orm := orm.NewOrm()

	//token_id 转换
	tokenId, ok := new(big.Int).SetString(req.TokenId, 10)
	if !ok {
		return nil, errors.New("convert token_id error.")
	}
	req.TokenId = common.BigToHash(tokenId).String()

	sqlStr := "SELECT `name`, symbol, logo, c.contract_type as token_type, a.account_address as holder from tb_contract c left join tb_contract_account a on c.contract_address = a.contract_address where c.contract_address = '" + req.ContractAddress + "' and a.token_id = '" + req.TokenId + "'"
	err := orm.Raw(sqlStr).QueryRow(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
