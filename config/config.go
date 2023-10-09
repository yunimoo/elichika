package config

import (
	"elichika/masterdata"
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
	ServerdataDb   = "assets/db/serverdata.db"

	MasterVersionGl = "2d61e7b4e89961c7" // read from GL database, so user can update db just by changing that
	MasterVersionJp = "b66ec2295e9a00aa" // ditto

	MainEng         *xorm.Engine
	MasterdataEngGl *xorm.Engine
	MasterdataEngJp *xorm.Engine

	MasterdataGl *masterdata.Masterdata
	MasterdataJp *masterdata.Masterdata

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
	if err != nil {
		panic(err)
	}
	err = eng.Ping()
	if err != nil {
		panic(err)
	}
	MainEng = eng
	MainEng.SetMaxOpenConns(50)
	MainEng.SetMaxIdleConns(10)

	MasterdataEngGl, err = xorm.NewEngine("sqlite", GlDatabasePath+"masterdata.db")
	if err != nil {
		panic(err)
	}
	MasterdataEngGl.SetMaxOpenConns(50)
	MasterdataEngGl.SetMaxIdleConns(10)
	MasterVersionGl = readMasterdataManinest(GlDatabasePath + "masterdata_a_en")
	MasterdataGl = new(masterdata.Masterdata)
	MasterdataGl.Init(MasterdataEngGl)

	MasterdataEngJp, err = xorm.NewEngine("sqlite", JpDatabasePath+"masterdata.db")
	if err != nil {
		panic(err)
	}
	MasterdataEngJp.SetMaxOpenConns(50)
	MasterdataEngJp.SetMaxIdleConns(10)
	MasterdataJp = new(masterdata.Masterdata)
	MasterdataJp.Init(MasterdataEngJp)

	MasterVersionJp = readMasterdataManinest(JpDatabasePath + "masterdata_a_ja")
	fmt.Println("gl master version:", MasterVersionGl)
	fmt.Println("jp master version:", MasterVersionJp)
}
