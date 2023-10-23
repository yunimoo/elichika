package config

import (
	"elichika/locale"
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

	GlDatabasePath = "assets/db/gl/"
	JpDatabasePath = "assets/db/jp/"
	// TODO: split into userdata.db and serverdata.db for gl / jp
	ServerdataDb = "assets/db/serverdata.db"

	MasterVersionGl = "2d61e7b4e89961c7" // read from GL database, so user can update db just by changing that
	MasterVersionJp = "b66ec2295e9a00aa" // ditto

	MasterdataEngGl *xorm.Engine
	MasterdataEngJp *xorm.Engine
	Locales         map[string](*locale.Locale)

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

func addLocale(path, language, masterVersion, startUpKey string, serverdataEngine *xorm.Engine) {
	locale := locale.Locale{
		Path:             path,
		Language:         language,
		MasterVersion:    masterVersion,
		StartUpKey:       startUpKey,
		ServerdataEngine: serverdataEngine,
	}
	locale.Load()
	Locales[language] = &locale
}

func init() {
	Conf = Load("./config.json")

	ServerdataEng, err := xorm.NewEngine("sqlite", ServerdataDb)
	utils.CheckErr(err)

	MasterVersionGl = readMasterdataManinest(GlDatabasePath + "masterdata_a_en")
	MasterVersionJp = readMasterdataManinest(JpDatabasePath + "masterdata_a_ja")

	Locales = make(map[string](*locale.Locale))
	addLocale(GlDatabasePath, "en", MasterVersionGl, "TxQFwgNcKDlesb93", ServerdataEng)
	addLocale(GlDatabasePath, "zh", MasterVersionGl, "TxQFwgNcKDlesb93", ServerdataEng)
	addLocale(GlDatabasePath, "ko", MasterVersionGl, "TxQFwgNcKDlesb93", ServerdataEng)
	addLocale(JpDatabasePath, "ja", MasterVersionJp, "5f7IZY1QrAX0D49g", ServerdataEng)

	fmt.Println("gl master version:", MasterVersionGl)
	fmt.Println("jp master version:", MasterVersionJp)
}
