package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"net/http"
	"encoding/json"
	// "strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func FetchProfile(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type FetchProfileReq struct {
		UserID int `json:"user_id"`
	}
	req := FetchProfileReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(req.UserID)
	
	signBody := GetUserData("fetchProfile.json")
	signBody, _ = sjson.Set(signBody, "profile_info.basic_info.name.dot_under_text",
		session.UserInfo.Name.DotUnderText)
	signBody, _ = sjson.Set(signBody, "profile_info.basic_info.introduction_message.dot_under_text",
		session.UserInfo.Message.DotUnderText)
	signBody, _ = sjson.Set(signBody, "profile_info.basic_info.emblem_id",
		session.UserInfo.EmblemID)
	signBody, _ = sjson.Set(signBody, "profile_info.basic_info.user_id", req.UserID)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetProfile(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	session := serverdb.GetSession(UserID)
	fmt.Println(reqBody)

	req := gjson.Parse(reqBody).Array()[0]
	if req.Get("name").String() != "" {
		session.UserInfo.Name.DotUnderText = gjson.Parse(reqBody).Array()[0].Get("name").String()
	} else if req.Get("nickname").String() != "" {
		session.UserInfo.Nickname.DotUnderText = gjson.Parse(reqBody).Array()[0].Get("nickname").String()
	} else if req.Get("message").String() != "" {
		session.UserInfo.Message.DotUnderText = gjson.Parse(reqBody).Array()[0].Get("message").String()
	}

	signBody := session.Finalize(GetData("setProfile.json"),	"user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetRecommendCard(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	fmt.Println(reqBody)
	session := serverdb.GetSession(UserID)
	cardMasterId := int(gjson.Parse(reqBody).Array()[0].Get("card_master_id").Int())
	session.UserInfo.RecommendCardMasterID = cardMasterId
	cardInfo := session.GetCard(cardMasterId)
	// set profile for basic profile
	SetUserData("fetchProfile.json", "profile_info.basic_info.recommend_card_master_id", cardMasterId)
	SetUserData("fetchProfile.json", "profile_info.basic_info.is_recommend_card_image_awaken", cardInfo.IsAwakeningImage)

	signBody := session.Finalize(GetData("setRecommendCard.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
