package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"
	"encoding/json"

	"fmt"
	// "math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	// "xorm.io/xorm"
)

// func FinishSession(ctx *gin.Context, session *serverdb.Session) {
// 	resp := SignResp(ctx.GetString("ep"), session.Finalize("user_model_diff"), config.SessionKey)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.String(http.StatusOK, resp)
// 	fmt.Println(resp)
// }

func SendCardInfoDiff(ctx *gin.Context, cardInfo *model.CardInfo) {
	userCardInfo := []any{}
	userCardInfo = append(userCardInfo, cardInfo.CardMasterID)
	userCardInfo = append(userCardInfo, *cardInfo)

	cardResp := GetUserData("userModelDiff.json")
	cardResp, _ = sjson.Set(cardResp, "user_model_diff.user_status", GetUserStatus())
	cardResp, _ = sjson.Set(cardResp, "user_model_diff.user_card_by_card_id", userCardInfo)
	resp := SignResp(ctx.GetString("ep"), cardResp, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	fmt.Println(resp)
}

func FetchTrainingTree(ctx *gin.Context) {
	signBody := GetData("fetchTrainingTree.json")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)


	// reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// type FetchTrainingTreeReq struct {
	// 	CardMasterId int `json:"card_master_id"`
	// }
	// req := FetchTrainingTreeReq{}
	// if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
	// 	panic(err)
	// }
	// user_id, err := strconv.Atoi(ctx.Query("u"))
	// CheckErr(err)
	// session := serverdb.GetSession(user_id)
	// signBody := `"{}"`
	// signBody, _ = sjson.Set(signBody, "user_card_training_tree_cell_list", session.GetTrainingTree(req.CardMasterId))
	// resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// ctx.Header("Content-Type", "application/json")
	// ctx.String(http.StatusOK, resp)
}

func LevelUpCard(ctx *gin.Context) {
	session := serverdb.GetSession(UserID)

	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]

	type LevelUpCardReq struct {
		CardMasterId int `json:"card_master_id"`
		AdditionalLevel int `json:"additional_level"`
	}

	req := LevelUpCardReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	cardInfo := session.GetCard(req.CardMasterId)
	cardInfo.Level += req.AdditionalLevel
	session.UpdateCard(cardInfo)
	signBody := session.Finalize(GetUserData("userModelDiff.json"), "user_model_diff")
	signBody, _ = sjson.Set(signBody, "user_model_diff.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	
	// TODO: Handle things like exp and gold cost

	// SendCardInfoDiff(ctx, &cardInfo)

	// SetUserData("userCard.json", "user_card_by_card_id." + key.String(), cardInfo)
	// SendCardInfoDiff(ctx, &cardInfo)
}

// func GradeUpCard(ctx *gin.Context) {
// 	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
// 	type GradeUpCardReq struct {
// 		CardMasterId int `json:"card_master_id"`
// 	}
// 	req := GradeUpCardReq{}
// 	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
// 		panic(err)
// 	}

// 	uid, _ := strconv.Atoi(ctx.Query("u"))
// 	session := serverdb.GetSession(uid)

// 	cardInfo := session.GetCard(req.CardMasterId)
// 	cardInfo.Grade += 1
// 	session.UpdateCard(cardInfo)

// 	// This correctly increase the limit, but doesn't increase bond level, and doesn't pop the screen after lb
// 	// FinishSession(ctx, &session)

// 	// we need to set user_info_trigger_card_grade_up_by_trigger_id

// 	type Trigger struct {
// 		TriggerId            int64 `json:"trigger_id"`
// 		CardMasterId         int   `json:"card_master_id"`
// 		BeforeLoveLevelLimit int   `json:"before_love_level_limit"`
// 		AfterLoveLevelLimit  int   `json:"after_love_level_limit"`
// 	}

// 	trigger := Trigger{}
// 	// this trigger show the pop up after limit break
// 	// TODO: we load the card again, the animation will be played again
// 	// this has something to do with the state of the game, as restarting fix this

// 	trigger.TriggerId = rand.Int63()
// 	trigger.CardMasterId = cardInfo.CardMasterID
// 	trigger.BeforeLoveLevelLimit = 500
// 	trigger.AfterLoveLevelLimit = 503

// 	session.AddCardGradeUpTrigger(trigger.TriggerId, trigger)

// 	resp := SignResp(ctx.GetString("ep"), session.Finalize("user_model_diff"), config.SessionKey)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.String(http.StatusOK, resp)
// 	fmt.Println(resp)

// 	// AccessCard(&cardInfo, req.CardMasterId, func (key, value *gjson.Result) {
// 	// 	cardInfo.Grade += 1
// 	// 	SetUserData("userCard.json", "user_card_by_card_id." + key.String(), cardInfo)
// 	// })
// 	// SendCardInfoDiff(ctx, &cardInfo)
// }

// func ActivateTrainingTreeCell(ctx *gin.Context) {
// 	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
// 	type ActivateTrainingTreeCellReq struct {
// 		CardMasterId  int   `json:"card_master_id"`
// 		CellMasterIds []int `json:"cell_master_ids"`
// 		PayType       int   `json:"pay_type"`
// 	}
// 	req := ActivateTrainingTreeCellReq{}
// 	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
// 		panic(err)
// 	}

// 	uid, _ := strconv.Atoi(ctx.Query("u"))
// 	session := serverdb.GetSession(uid)

// 	db := GetMasterdataDb()

// 	type TrainingTreeMapping struct {
// 		// Id int  // card master id
// 		TrainingTreeMappingMId                  int // the training tree
// 		// TrainingTreeCardParamMId int  // same as card master id
// 		TrainingTreeCardPassiveSkillIncreaseMId int // 1 to 2
// 	}

// 	treeMapping := TrainingTreeMapping{}
// 	exists, err := db.Table("m_training_tree").Where("id = ?", req.CardMasterId).
// 		Cols("training_tree_mapping_m_id", "training_tree_card_passive_skill_increase_m_id").Get(&treeMapping)
// 	if (err != nil) || (!exists) {
// 		panic(err)
// 	}

// 	type TrainingTreeCell struct {
// 		// Id int // tree id
// 		CellId               int // cell id
// 		TrainingTreeCellType int // type of the cell
// 		TrainingContentNo    int // the content of the cell
// 		// RequiredGrade int // the limit break required, no need to read this here
// 		TrainingTreeCellItemSetMId int // the set of items used to unlock this cell
// 		// SnsCoin int // always 1, maybe we could earn gem by unlocking cell?
// 	}

// 	refCells := []TrainingTreeCell{}
// 	db.ShowSQL(true)
// 	err = db.Table("m_training_tree_cell_content").Where("id = ?", treeMapping.TrainingTreeMappingMId).
// 		Cols("cell_id", "training_tree_cell_type", "training_content_no", "training_tree_cell_item_set_m_id").
// 		OrderBy("cell_id").Find(&refCells)

// 	if (err != nil) {
// 		panic(err)
// 	}

// 	// stats reward is in "m_training_tree_card_param"
// 	type TrainingTreeCardParam struct {
// 		TrainingContentType int
// 		Value               int
// 	}
// 	cellParams := []TrainingTreeCardParam{}
// 	db.Table("m_training_tree_card_param").
// 		Where("id = ?", req.CardMasterId).
// 		Cols("training_content_type", "value").OrderBy("training_content_no").Find(&cellParams)

// 	increasedStats := [5]int{}
// 	card := session.GetCard(req.CardMasterId)
// 	for _, cId := range req.CellMasterIds {
// 		switch refCells[cId].TrainingTreeCellType {
// 		case 2:
// 			// param cells, reference the 
// 			paramCell := &cellParams[refCells[cId].TrainingContentNo - 1]
// 			increasedStats[paramCell.TrainingContentType] +=  paramCell.Value
// 		case 3:
// 			// voice
// 			// TODO: Unlock voice
// 		case 4:
// 			// story
// 			// TODO: Unlock story
// 		case 5:
// 			// idolize
// 			card.IsAwakening = true
// 			card.IsAwakeningImage = true
// 			// idolize stats is calculated by the flag?
// 		case 7:
// 			// skill
// 			card.ActiveSkillLevel += 1
// 		case 8:
// 			// insight
// 			card.MaxFreePassiveSkill += 1
// 		case 9:
// 			// ability
// 			card.PassiveSkillALevel += 1
// 		}
// 	}

// 	card.TrainingLife += increasedStats[2]
// 	card.TrainingAttack += increasedStats[3]
// 	card.TrainingDexterity += increasedStats[4]
// 	card.TrainingActivatedCellCount += len(req.CellMasterIds)

// 	if card.TrainingActivatedCellCount + 1 == len(refCells) {
// 		card.IsAllTrainingActivated = true
// 	}
	
// 	session.UpdateCard(card)

// 	// set "user_card_training_tree_cell_list" to the cell unlocked and insert the cell to db
// 	unlockedCells := []dbmodel.DbTrainingTreeCell{}
// 	t, _ := strconv.ParseInt(ctx.Query("t"), 10, 64)
// 	for _, cId := range req.CellMasterIds {
// 		cell := dbmodel.DbTrainingTreeCell{}
// 		cell.UserId = uid
// 		cell.CardMasterId = req.CardMasterId
// 		cell.CellId = cId
// 		cell.ActivatedAt = t
// 		unlockedCells = append(unlockedCells, cell)
// 	}
	


// 	session.InsertTrainingCells(&unlockedCells)
// 	jsonResp := session.Finalize("user_model_diff")
// 	jsonResp, _ = sjson.Set(jsonResp, "user_card_training_tree_cell_list", session.GetTrainingTree(req.CardMasterId))
// 	resp := SignResp(ctx.GetString("ep"), jsonResp, config.SessionKey)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.String(http.StatusOK, resp)
// 	fmt.Println(resp)

// }
