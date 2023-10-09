package db

import (
	"elichika/config"
	"elichika/serverdb"
	"elichika/utils"

	"fmt"
	"strconv"
	"time"
)

func InsertAll(args []string) {
	masterdata := config.MasterdataGl
	if len(args) < 1 {
		fmt.Println("Invalid params:", args)
		fmt.Println("Required: gl/jp [optional user_id]")
	}
	if args[0] == "jp" {
		masterdata = config.MasterdataJp
	}
	userID := 588296696
	if len(args) > 1 {
		uid, err := strconv.Atoi(args[1])
		utils.CheckErr(err)
		userID = uid
	}
	session := serverdb.GetSession(nil, userID)

	for i, a := range masterdata.Accessory.Accessory {
		accessory := session.GetUserAccessory(time.Now().UnixNano() + int64(i))
		accessory.AccessoryMasterID = a.MasterID
		accessory.Level = 1
		accessory.Exp = 0
		accessory.Grade = 0
		accessory.Attribute = a.Attribute
		accessory.PassiveSkill1ID = *a.Grade[0].PassiveSkill1MasterID
		accessory.PassiveSkill2ID = a.Grade[0].PassiveSkill2MasterID
		session.UpdateUserAccessory(accessory)
	}
	session.Finalize("", "")
}

func Accessory(args []string) {
	switch args[0] {
	case "insert_all":
		InsertAll(args[1:])
	}
}
