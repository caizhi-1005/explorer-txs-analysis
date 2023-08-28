package dbModels

import (
	"github.com/astaxie/beego/orm"
	"github.com/server/txs-analysis/constant"
	"time"
)

type TbUserAddressNote struct {
	ID         int64     `orm:"column(id);pk"`
	UserID     int64     `orm:"column(user_id);null" description:"用户ID"`
	Address    string    `orm:"column(address);size(42);null" description:"地址"`
	Tag        string    `orm:"column(tag);size(35);null" description:"标签"`
	Note       string    `orm:"column(note);size(500);null" description:"地址备注"`
	IsFavorite int       `orm:"column(is_favorite);null" description:"是否标记了喜欢,0否1是"`
	Status     int       `orm:"column(status);null" description:"状态 1-正常 2-冻结"`
	IsDeleted  int       `orm:"column(is_deleted);null" description:"删除状态 0-正常 1-删除"`
	CreateTime time.Time `orm:"column(create_time);type(timestamp);null;auto_now_add" description:"创建时间"`
	UpdateTime time.Time `orm:"column(update_time);type(timestamp);null;auto_now_add" description:"更新时间"`
}

func (t *TbUserAddressNote) TableName() string {
	return "tb_user_address_note"
}

func init() {
	orm.RegisterModel(new(TbUserAddressNote))
}

// ExistUserAddressNote 判断用户是否添加了此地址的备注
func ExistUserAddressNote(userID int, address string) (exist bool) {
	o := orm.NewOrm()
	o.Using(constant.USER_DB_ALIAS)
	exist = o.QueryTable(new(TbUserAddressNote).TableName()).
		//Filter("user_id", userID).
		Filter("address", address).
		Filter("is_deleted", 0).
		Exist()
	return
}

// GetUserAddressNote 查询用户地址备注
func GetUserAddressNote(address string) (note *TbUserAddressNote, err error) {
	note = new(TbUserAddressNote)
	o := orm.NewOrm()
	o.Using(constant.USER_DB_ALIAS)
	err = o.QueryTable(new(TbUserAddressNote).TableName()).
		//Filter("user_id", userID).
		Filter("address", address).
		Filter("is_deleted", 0).
		One(note)
	if err == orm.ErrNoRows {
		err = nil
	}
	return
}
