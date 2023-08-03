package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"

	"fmt"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func UpdateCardNewFlag(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	fmt.Println(reqBody.String())
	session := serverdb.GetSession(UserID)
	
	signBody := session.Finalize(GetData("updateCardNewFlag.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ChangeIsAwakeningImage(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	req := model.CardAwakeningReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(UserID)
	cardInfo := session.GetCard(req.CardMasterID)
	cardInfo.IsAwakeningImage = req.IsAwakeningImage
	session.UpdateCard(cardInfo)

	cardResp := session.Finalize(GetData("changeIsAwakeningImage.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), cardResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ChangeFavorite(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	req := model.CardFavoriteReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(UserID)
	cardInfo := session.GetCard(req.CardMasterID)
	cardInfo.IsFavorite = req.IsFavorite
	session.UpdateCard(cardInfo)

	cardResp := session.Finalize(GetData("changeFavorite.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), cardResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func GetOtherUserCard(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	userCardReq := model.UserCardReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &userCardReq); err != nil {
		panic(err)
	}
	// fmt.Println(liveStartReq)

	var newUserCardInfo model.PartnerCard
	var cardInfo string
	partnerList := gjson.Parse(GetData("fetchLivePartners.json")).Get("partner_select_state.live_partners")
	partnerList.ForEach(func(k, v gjson.Result) bool {
		userId := v.Get("user_id").Int()
		if userId == userCardReq.UserID {
			v.Get("card_by_category").ForEach(func(kk, vv gjson.Result) bool {
				if vv.IsObject() {
					cardId := vv.Get("card_master_id").Int()
					if cardId == userCardReq.CardMasterID {
						cardInfo = vv.String()
						// fmt.Println(cardInfo)
						return false
					}
				}
				return true
			})
			return false
		}
		return true
	})

	if err := json.Unmarshal([]byte(cardInfo), &newUserCardInfo); err != nil {
		panic(err)
	}

	userCardResp := GetData("getOtherUserCard.json")
	userCardResp, _ = sjson.Set(userCardResp, "other_user_card", newUserCardInfo)
	resp := SignResp(ctx.GetString("ep"), userCardResp, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
