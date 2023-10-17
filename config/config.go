package config

import (
	"elichika/gamedata"
	"elichika/utils"

	"fmt"
	"os"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var (
	// Taken from:
	// https://github.com/RayFirefist/SukuStar_Datamine/blob/master/lib/sifas_api/sifas.py#L120
	// https://github.com/RayFirefist/SukuStar_Datamine/blob/master/lib/sifas_api/sifas.py#L400
	ServerEventReceiverKey = "31f1f9dc7ac4392d1de26acf99d970e425b63335b461e720c73d6914020d6014"
	JaKey                  = "78d53d9e645a0305602174e06b98d81f638eaf4a84db19c756866fddac360c96"

	SessionKey = "12345678123456781234567812345678"

	MainDb         = "assets/main.db"
	GlDatabasePath = "assets/db/gl/"
	JpDatabasePath = "assets/db/jp/"
	// TODO: split into userdata.db and serverdata.db for gl / jp
	ServerdataDb = "assets/db/serverdata.db"

	MasterVersionGl = "2d61e7b4e89961c7" // read from GL database, so user can update db just by changing that
	MasterVersionJp = "b66ec2295e9a00aa" // ditto

	MainEng         *xorm.Engine
	MasterdataEngGl *xorm.Engine
	MasterdataEngJp *xorm.Engine

	GamedataGl *gamedata.Gamedata
	GamedataJp *gamedata.Gamedata

	Conf = &AppConfigs{}

	PresetDataPath = "assets/preset/"
	UserDataPath   = "assets/userdata/"
)

func readMasterdataManinest(path string) string {
	file, err := os.Open(path)
	utils.CheckErr(err)
	buff := make([]byte, 1024)
	count, err := file.Read(buff)
	utils.CheckErr(err)
	if count < 21 {
		panic("file too short")
	}
	length := int(buff[20])
	version := string(buff[21 : 21+length])
	return version
}

func init() {
	Conf = Load("./config.json")

	eng, err := xorm.NewEngine("sqlite", MainDb)
	utils.CheckErr(err)
	err = eng.Ping()
	utils.CheckErr(err)
	MainEng = eng
	MainEng.SetMaxOpenConns(50)
	MainEng.SetMaxIdleConns(10)

	ServerdataEng, err := xorm.NewEngine("sqlite", ServerdataDb)
	utils.CheckErr(err)
	ServerdataEng.SetMaxOpenConns(50)
	ServerdataEng.SetMaxIdleConns(10)

	MasterdataEngGl, err = xorm.NewEngine("sqlite", GlDatabasePath+"masterdata.db")
	utils.CheckErr(err)
	MasterdataEngGl.SetMaxOpenConns(50)
	MasterdataEngGl.SetMaxIdleConns(10)
	MasterVersionGl = readMasterdataManinest(GlDatabasePath + "masterdata_a_en")
	GamedataGl = new(gamedata.Gamedata)
	GamedataGl.Init(MasterdataEngGl, ServerdataEng)

	MasterdataEngJp, err = xorm.NewEngine("sqlite", JpDatabasePath+"masterdata.db")
	utils.CheckErr(err)
	MasterdataEngJp.SetMaxOpenConns(50)
	MasterdataEngJp.SetMaxIdleConns(10)
	MasterVersionJp = readMasterdataManinest(JpDatabasePath + "masterdata_a_ja")
	GamedataJp = new(gamedata.Gamedata)
	GamedataJp.Init(MasterdataEngJp, ServerdataEng)

	fmt.Println("gl master version:", MasterVersionGl)
	fmt.Println("jp master version:", MasterVersionJp)
}
