package locale

import (
	"elichika/dictionary"
	"elichika/gamedata"
	"elichika/utils"

	"xorm.io/xorm"
)

type Locale struct {
	// Loaded bool
	Path             string
	Language         string
	StartUpKey       string
	MasterVersion    string
	MasterdataEngine *xorm.Engine
	ServerdataEngine *xorm.Engine
	Gamedata         *gamedata.Gamedata
	Dictionary       *dictionary.Dictionary
}

func (locale *Locale) Load() {
	var err error
	locale.MasterdataEngine, err = xorm.NewEngine("sqlite", locale.Path+"masterdata.db")
	utils.CheckErr(err)
	locale.MasterdataEngine.SetMaxOpenConns(50)
	locale.MasterdataEngine.SetMaxIdleConns(10)
	locale.Dictionary = new(dictionary.Dictionary)
	locale.Dictionary.Init(locale.Path, locale.Language)
	locale.Gamedata = new(gamedata.Gamedata)
	locale.Gamedata.Init(locale.MasterdataEngine, locale.ServerdataEngine, locale.Dictionary)
}
