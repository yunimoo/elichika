package handler

import (
	"elichika/config"
	"elichika/encrypt"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func StartUp(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")

	var mask64 string
	req := gjson.Parse(reqBody)
	req.ForEach(func(key, value gjson.Result) bool {
		if value.Get("mask").String() != "" {
			mask64 = value.Get("mask").String()
			return false
		}
		return true
	})
	// fmt.Println("Request data:", req.String())
	// fmt.Println("Mask:", mask64)

	mask, err := base64.StdEncoding.DecodeString(mask64)
	if err != nil {
		panic(err)
	}
	randomBytes := encrypt.RSA_DecryptOAEP(mask, "privatekey.pem")
	// fmt.Println("Random Bytes:", randomBytes)

	newKey := utils.Xor(randomBytes, []byte(config.SessionKey))
	newKey64 := base64.StdEncoding.EncodeToString(newKey)
	// fmt.Println("Session Key:", newKey64)

	startupBody := GetData("startup.json")
	startupBody, _ = sjson.Set(startupBody, "authorization_key", newKey64)
	resp := SignResp(ctx.GetString("ep"), startupBody, StartUpKey)
	// fmt.Println("Response:", resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func Login(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")

	var mask64 string
	req := gjson.Parse(reqBody)
	req.ForEach(func(key, value gjson.Result) bool {
		if value.Get("mask").String() != "" {
			mask64 = value.Get("mask").String()
			return false
		}
		return true
	})
	// fmt.Println("Request data:", req.String())
	// fmt.Println("Mask:", mask64)

	mask, err := base64.StdEncoding.DecodeString(mask64)
	if err != nil {
		panic(err)
	}
	randomBytes := encrypt.RSA_DecryptOAEP(mask, "privatekey.pem")
	// fmt.Println("Random Bytes:", randomBytes)

	serverEventReceiverKey, err := hex.DecodeString(config.ServerEventReceiverKey)
	if err != nil {
		panic(err)
	}

	jaKey, err := hex.DecodeString(config.JaKey)
	if err != nil {
		panic(err)
	}

	newKey := utils.Xor(randomBytes, []byte(config.SessionKey))
	newKey = utils.Xor(newKey, serverEventReceiverKey)
	newKey = utils.Xor(newKey, jaKey)
	newKey64 := base64.StdEncoding.EncodeToString(newKey)
	// fmt.Println("Session Key:", newKey64)

	loginBody := GetData("login.json")
	loginBody, _ = sjson.Set(loginBody, "session_key", newKey64)
	loginBody, _ = sjson.Set(loginBody, "user_model.user_status", GetUserStatus())

	/* ======== UserData ======== */
	// live decks
	liveDeckData := gjson.Parse(GetLiveDeckData())
	loginBody, _ = sjson.Set(loginBody, "user_model.user_live_deck_by_id", liveDeckData.Get("user_live_deck_by_id").Value())

	var liveParty []any
	decoder := json.NewDecoder(strings.NewReader(liveDeckData.Get("user_live_party_by_id").String()))
	decoder.UseNumber()
	err = decoder.Decode(&liveParty)
	CheckErr(err)
	loginBody, _ = sjson.Set(loginBody, "user_model.user_live_party_by_id", liveParty)

	// member settings
	memberData := gjson.Parse(GetUserData("memberSettings.json"))
	loginBody, _ = sjson.Set(loginBody, "user_model.user_member_by_member_id", memberData.Get("user_member_by_member_id").Value())

	// lesson decks
	lessonData := gjson.Parse(GetUserData("lessonDeck.json"))
	loginBody, _ = sjson.Set(loginBody, "user_model.user_lesson_deck_by_id", lessonData.Get("user_lesson_deck_by_id").Value())

	// user cards
	fmt.Println("User logins: ", UserID)
	session := serverdb.GetSession(UserID)
	var userCards []any
	dbCards := session.GetAllCards()

	for _, cardInfo := range dbCards {
		userCards = append(userCards, cardInfo.CardMasterID)
		userCards = append(userCards, cardInfo)
	}

	if len(dbCards) == 0 {
		fmt.Println("importing json card data to db")
		// first time, parse the json and insert to db
		cardData := gjson.Parse(GetUserData("userCard.json"))
		cardInfo := model.CardInfo{}
		var cards []model.CardInfo
		cardData.Get("user_card_by_card_id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &cardInfo); err != nil {
					panic(err)
				}
				cardInfo.UserID = UserID
				userCards = append(userCards, cardInfo.CardMasterID)
				userCards = append(userCards, cardInfo)
				cards = append(cards, cardInfo)
			}
			return true
		})
		session.InsertCards(cards)
	}

	loginBody, _ = sjson.Set(loginBody, "user_model.user_card_by_card_id", userCards)

	// user accessory
	var UserAccessory []any
	decoder = json.NewDecoder(strings.NewReader(
		gjson.Parse(GetUserAccessoryData()).Get("user_accessory_by_user_accessory_id").String()))
	decoder.UseNumber()
	err = decoder.Decode(&UserAccessory)
	CheckErr(err)
	loginBody, _ = sjson.Set(loginBody, "user_model.user_accessory_by_user_accessory_id", UserAccessory)
	/* ======== UserData ======== */

	resp := SignResp(ctx.GetString("ep"), loginBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
