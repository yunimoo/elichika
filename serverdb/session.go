package serverdb

import (
	"elichika/config"
	"elichika/model"
	"elichika/utils"

	"encoding/json"
	"fmt"
	// "github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func DbGetUserData(fileName string) string {
	userDataFile := config.UserDataPath + fileName
	if utils.PathExists(userDataFile) {
		return utils.ReadAllText(userDataFile)
	}

	presetDataFile := config.PresetDataPath + fileName
	if !utils.PathExists(presetDataFile) {
		panic("File not exists")
	}

	userData := utils.ReadAllText(presetDataFile)
	utils.WriteAllText(userDataFile, userData)

	return userData
}

// A session is a complete transation between server and client
// so 1 session per request
// A session fetch the data needs to be modified.
type Session struct {
	UserInfo  model.UserInfo
	CardDiffs map[int]model.CardInfo
	// MemberDiffs         map[int]dbmodel.DbUserMember
	// CardGradeUpTriggers []any
}

// Update the diff into the Db, and return the delta patch
// the delta patch is sometime called "user_model", sometime called "user_model_diff"
// that is the mainKey here
// they all have the same structure
func (session *Session) Finalize(mainKey string) string {
	session.FinalizeCardDiffs()
	return mainKey
	// modelDiff := DbGetUserData("userModelDiff.json")
	// if mainKey == "user_model" {
	// 	modelDiff = DbGetUserData("userModel.json")
	// }
	// modelDiff, _ = sjson.Set(modelDiff, mainKey+".user_status", session.UserInfo)
	// modelDiff, _ = sjson.Set(modelDiff, mainKey+".user_status.gdpr_version", 4)
	// // modelDiff, _ = sjson.Set(modelDiff, mainKey + ".user_member_by_member_id", session.FinalizeMemberDiffs())
	// modelDiff, _ = sjson.Set(modelDiff, mainKey+".user_card_by_card_id", session.FinalizeCardDiffs())
	// modelDiff, _ = sjson.Set(modelDiff, mainKey+".user_info_trigger_card_grade_up_by_trigger_id", session.FinalizeCardGradeUpTrigger())
	// return modelDiff
	// for memberMasterId, member := range session.MemberDiffs {
	// 	affected, err := Engine.Table("s_user_card").
	// 		Where("user_id = ? AND card_master_id = ?", session.UserInfo.UserId, cardMasterId).Upddate(&card)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if affected != 1 {
	// 		panic()
	// 	}
	// }
}

// fetch the user, this is always sent back to client
func (session *Session) InitUser(userId int) {
	session.UserInfo.UserId = userId
	exists, err := Engine.Table("s_user_info").Where("user_id = ?", userId).Get(&session.UserInfo)
	if err != nil {
		panic(err)
	}
	if !exists { // create one
		fmt.Println("Insert new user: ", userId)
		data := utils.ReadAllText("assets/userdata/userStatus.json")
		if err := json.Unmarshal([]byte(data), &session.UserInfo); err != nil {
			panic(err)
		}

		// insert into the db
		_, err := Engine.Table("s_user_info").AllCols().Insert(&session.UserInfo)
		if err != nil {
			panic(err)
		}
	}
}

func GetSession(userId int) Session {
	s := Session{}
	s.CardDiffs = make(map[int]model.CardInfo)
	// s.MemberDiffs = make(map[int]dbmodel.DbUserMember)
	s.InitUser(userId)
	return s
}
