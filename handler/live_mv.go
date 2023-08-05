package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func LiveMvStart(ctx *gin.Context) {
	session := serverdb.GetSession(UserID)
	signBody := session.Finalize(GetData("liveMvStart.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func LiveMvSaveDeck(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)

	reqData := gjson.Parse(reqBody).Array()[0]

	saveReq := model.LiveSaveDeckReq{}
	err := json.Unmarshal([]byte(reqData.String()), &saveReq)
	if err != nil {
		panic(err)
	}

	userLiveMvDeckInfo := model.UserLiveMvDeckInfo{
		LiveMasterID: saveReq.LiveMasterID,
	}
	deckJsonBytes, err := json.Marshal(userLiveMvDeckInfo)
	CheckErr(err)
	deckJson := string(deckJsonBytes)

	for k, v := range saveReq.MemberMasterIDByPos {
		if k%2 == 0 {
			memberId := saveReq.MemberMasterIDByPos[k+1]
			deckJson, err = sjson.Set(deckJson, fmt.Sprintf("member_master_id_%d", v), memberId)
		}
	}
	for k, v := range saveReq.SuitMasterIDByPos {
		if k%2 == 0 {
			suitId := saveReq.SuitMasterIDByPos[k+1]
			deckJson, err = sjson.Set(deckJson, fmt.Sprintf("suit_master_id_%d", v), suitId)
		}
	}
	err = json.Unmarshal([]byte(deckJson), &userLiveMvDeckInfo)
	CheckErr(err)

	var userLiveMvDeckCustomByID []any
	userLiveMvDeckCustomByID = append(userLiveMvDeckCustomByID, saveReq.LiveMasterID)
	userLiveMvDeckCustomByID = append(userLiveMvDeckCustomByID, userLiveMvDeckInfo)

	session := serverdb.GetSession(UserID)
	signBody := GetData("liveMvSaveDeck.json")
	signBody = session.Finalize(signBody, "user_model")
	if saveReq.LiveMvDeckType == 1 {
		signBody, _ = sjson.Set(signBody, "user_model.user_live_mv_deck_by_id", userLiveMvDeckCustomByID)
	} else {
		signBody, _ = sjson.Set(signBody, "user_model.user_live_mv_deck_custom_by_id", userLiveMvDeckCustomByID)
	}

	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
