package nebulaModels

import (
	"github.com/astaxie/beego"
	"github.com/zhihu/norm"
)

func PrepareAddress(db *norm.DB) error {
	createSchema :=
		"CREATE TAG IF NOT EXISTS address(address string, type int);" +
			"CREATE TAG INDEX address_index on address();"
	_, err := db.Execute(createSchema)
	return err
}

func InsertAddress(db *norm.DB, address *Address) error {
	err := db.InsertVertex(address)
	if err != nil {
		beego.Error("InsertAddress err: ", err, " address:", address.Address)
	} else {
		beego.Info("InsertAddress success! address:", address.Address)
	}
	return err
}
