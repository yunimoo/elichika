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
		err = Engine.Table(tableName).CreateTable(structure)
		if err != nil {
			panic(err)
		}
	}
}

func InitTables() {
	type DbUser struct {
		model.UserStatus `xorm:"extends"`
		model.DBUserProfileLiveStats `xorm:"extends"`
	}
	InitTable("s_user_info", DbUser{})
	type DbCard struct {
		model.CardInfo `xorm:"extends"`
		model.DBCardPlayInfo `xorm:"extends"`
	}
	InitTable("s_user_card", DbCard{})
	InitTable("s_user_training_tree_cell", model.TrainingTreeCell{})
	InitTable("s_user_member", model.UserMemberInfo{})
	InitTable("s_user_member_love_panel", model.UserMemberLovePanel{})
	InitTable("s_user_lesson_deck", model.UserLessonDeck{})
	InitTable("s_user_live_deck", model.UserLiveDeck{})
	InitTable("s_user_live_party", model.UserLiveParty{})
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
