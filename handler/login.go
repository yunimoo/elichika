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
	fmt.Println("User logins: ", UserID)
	session := serverdb.GetSession(UserID)

	// live decks
	dbLiveDecks := session.GetAllLiveDecks()
	if len(dbLiveDecks) == 0 {
		fmt.Println("importing live deck data to db")
		liveDeckInfo := model.UserLiveDeck{}
		liveDeckData := gjson.Parse(GetLiveDeckData())
		liveDeckData.Get("user_live_deck_by_id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &liveDeckInfo); err != nil {
					panic(err)
				}
				liveDeckInfo.UserID = UserID
				dbLiveDecks = append(dbLiveDecks, liveDeckInfo)
			}
			return true
		})
		session.InsertLiveDecks(dbLiveDecks)
	}
	userLiveDecks := []any{}
	for _, liveDeckInfo := range dbLiveDecks {
		userLiveDecks = append(userLiveDecks, liveDeckInfo.UserLiveDeckID)
		userLiveDecks = append(userLiveDecks, liveDeckInfo)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_live_deck_by_id", userLiveDecks)

	dbLiveParties := session.GetAllLiveParties()
	if len(dbLiveParties) == 0 {
		fmt.Println("importing live party data to db")
		livePartyInfo := model.UserLiveParty{}
		liveDeckData := gjson.Parse(GetLiveDeckData())
		var livePartyData []json.RawMessage
		decoder := json.NewDecoder(strings.NewReader(liveDeckData.Get("user_live_party_by_id").String()))
		decoder.UseNumber()
		err = decoder.Decode(&livePartyData)
		CheckErr(err)
		for i := 1; i < len(livePartyData); i += 2 {
			err := json.Unmarshal(livePartyData[i], &livePartyInfo)
			CheckErr(err)
			livePartyInfo.UserID = UserID
			dbLiveParties = append(dbLiveParties, livePartyInfo)
		}
		session.InsertLiveParties(dbLiveParties)
	}
	userLiveParties := []any{}
	for _, livePartyInfo := range dbLiveParties {
		userLiveParties = append(userLiveParties, livePartyInfo.PartyID)
		userLiveParties = append(userLiveParties, livePartyInfo)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_live_party_by_id", userLiveParties)

	// member settings
	dbMembers := session.GetAllMembers()
	if len(dbMembers) == 0 {
		// insert from json
		fmt.Println("importing member data to db")
		userMemberInfo := model.UserMemberInfo{}
		memberData := gjson.Parse(GetUserData("memberSettings.json"))
		memberData.Get("user_member_by_member_id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &userMemberInfo); err != nil {
					panic(err)
				}
				userMemberInfo.UserID = UserID
				dbMembers = append(dbMembers, userMemberInfo)
			}
			return true
		})
		session.InsertMembers(dbMembers)
	}
	var userMembers []any
	for _, memberInfo := range dbMembers {
		userMembers = append(userMembers, memberInfo.MemberMasterID)
		userMembers = append(userMembers, memberInfo)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_member_by_member_id", userMembers)

	// lesson decks
	dbLessonDecks := session.GetAllLessonDecks()
	if len(dbLessonDecks) == 0 {
		// insert from json
		fmt.Println("importing lesson deck data to db")
		lessonData := gjson.Parse(GetUserData("lessonDeck.json"))
		userLessonDeck := model.UserLessonDeck{}
		lessonData.Get("user_lesson_deck_by_id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &userLessonDeck); err != nil {
					panic(err)
				}
				userLessonDeck.UserID = UserID
				dbLessonDecks = append(dbLessonDecks, userLessonDeck)
			}
			return true
		})
		session.InsertLessonDecks(dbLessonDecks)
	}

	userLessonDecks := []any{}
	for _, userLessonDeck := range dbLessonDecks {
		userLessonDecks = append(userLessonDecks, userLessonDeck.UserLessonDeckID)
		userLessonDecks = append(userLessonDecks, userLessonDeck)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_lesson_deck_by_id", userLessonDecks)

	// user cards
	dbCards := session.GetAllCards()
	if len(dbCards) == 0 {
		// first time, parse the json and insert to db
		fmt.Println("importing json card data to db")
		cardData := gjson.Parse(GetUserData("userCard.json"))
		cardInfo := model.CardInfo{}
		cardData.Get("user_card_by_card_id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &cardInfo); err != nil {
					panic(err)
				}
				cardInfo.UserID = UserID
				dbCards = append(dbCards, cardInfo)
			}
			return true
		})
		session.InsertCards(dbCards)
	}

	userCards := []any{}
	for _, cardInfo := range dbCards {
		userCards = append(userCards, cardInfo.CardMasterID)
		userCards = append(userCards, cardInfo)
	}
	loginBody, _ = sjson.Set(loginBody, "user_model.user_card_by_card_id", userCards)

	// user accessory
	var UserAccessory []any
	decoder := json.NewDecoder(strings.NewReader(
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
