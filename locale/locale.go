package locale

import (
	"elichika/config"
	"elichika/dictionary"
	"elichika/gamedata"
	"elichika/serverdata"
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

var (
	Locales                map[string](*Locale)
	sharedServerdataEngine *xorm.Engine
)

func addLocale(path, language, masterVersion, startUpKey string) {
	locale := Locale{
		Path:             path,
		Language:         language,
		MasterVersion:    masterVersion,
		StartUpKey:       startUpKey,
		ServerdataEngine: serverdata.Engine,
	}
	locale.Load()
	Locales[language] = &locale
}

func init() {
	Locales = make(map[string](*Locale))
	addLocale(config.GlDatabasePath, "en", config.MasterVersionGl, config.GlStartUpKey)
	addLocale(config.GlDatabasePath, "zh", config.MasterVersionGl, config.GlStartUpKey)
	addLocale(config.GlDatabasePath, "ko", config.MasterVersionGl, config.GlStartUpKey)
	addLocale(config.JpDatabasePath, "ja", config.MasterVersionJp, config.JpStartUpKey)
}
