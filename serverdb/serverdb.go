package serverdb

import (
	"elichika/config"
	"elichika/model"

	"fmt"

	"xorm.io/xorm"
)

var (
	Engine *xorm.Engine
)

func InitTable(tableName string, structure interface{}) {
	exist, err := Engine.Table(tableName).IsTableExist(tableName)
	if err != nil {
		panic(err)
	}

	if !exist {
		fmt.Println("Creating new table: ", tableName)
		// err = Engine.Table(tableName).CreateTable(&dbmodel.UserInfo{})
		err = Engine.Table(tableName).CreateTable(structure)
		if err != nil {
			panic(err)
		}
		// Engine.Table(tableName).Commit()
	}
}

func InitTables() {
	InitTable("s_user_info", model.UserInfo{})
	InitTable("s_user_card", model.CardInfo{})
	// InitTable("s_user_member", dbmodel.DbUserMember{})
	// InitTable("s_user_training_tree_cell", dbmodel.DbTrainingTreeCell{})
}

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite", config.ServerdataDb)
	if err != nil {
		panic(err)
	}
	Engine.SetMaxOpenConns(50)
	Engine.SetMaxIdleConns(10)
	// Engine.ShowSQL(true)

	InitTables()
}
