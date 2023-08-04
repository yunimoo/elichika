package handler

import (
	"bytes"
	"elichika/config"
	// "elichika/database"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func SaveDeckAll(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	req := model.SaveDeckReq{}
	decoder := json.NewDecoder(strings.NewReader(reqBody.String()))
	decoder.UseNumber()
	err := decoder.Decode(&req)
	CheckErr(err)

	session := serverdb.GetSession(UserID)
	deckInfo := session.GetUserLiveDeck(req.DeckID)

	for i := 0; i < 9; i++ {
		if req.CardWithSuit[i*2+1] == 0 {
			req.CardWithSuit[i*2+1] = GetMemberDefaultSuitByCardMasterId(req.CardWithSuit[i*2])
		}
	}

	deckByte, _ := json.Marshal(deckInfo)
	deckJson := string(deckByte)
	for i := 0; i < 9; i++ {
		deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("card_master_id_%d", i+1), req.CardWithSuit[i*2])
		deckJson, _ = sjson.Set(deckJson, fmt.Sprintf("suit_master_id_%d", i+1), req.CardWithSuit[i*2+1])
	}

	if err := json.Unmarshal([]byte(deckJson), &deckInfo); err != nil {
		panic(err)
	}
	// fmt.Println(deckInfo)

	session.UpdateUserLiveDeck(deckInfo)

	for k, v := range req.SquadDict {
		if k%2 == 0 {
			partyId, err := v.(json.Number).Int64()
			if err != nil {
				panic(err)
			}
			// fmt.Println("Party ID:", partyId)

			rDictInfo, err := json.Marshal(req.SquadDict[k+1])
			CheckErr(err)

			dictInfo := model.DeckSquadDict{}
			decoder := json.NewDecoder(bytes.NewReader(rDictInfo))
			decoder.UseNumber()
			err = decoder.Decode(&dictInfo)
			CheckErr(err)
			// fmt.Println("Party Info:", dictInfo)

			roleIds := []int{}
			err = MainEng.Table("m_card").
				Where("id IN (?,?,?)", dictInfo.CardMasterIds[0], dictInfo.CardMasterIds[1], dictInfo.CardMasterIds[2]).
				Cols("role").Find(&roleIds)
			CheckErr(err)
			// fmt.Println("roleIds:", roleIds)

			partyInfo := model.UserLiveParty{}
			partyInfo.UserID = UserID
			partyInfo.PartyID = int(partyId)
			partyIcon, partyName := GetPartyInfoByRoleIds(roleIds)
			partyInfo.Name.DotUnderText = GetRealPartyName(partyName)
			partyInfo.UserLiveDeckID = req.DeckID
			partyInfo.IconMasterID = partyIcon
			partyInfo.CardMasterID1 = dictInfo.CardMasterIds[0]
			partyInfo.CardMasterID2 = dictInfo.CardMasterIds[1]
			partyInfo.CardMasterID3 = dictInfo.CardMasterIds[2]
			partyInfo.UserAccessoryID1 = dictInfo.UserAccessoryIds[0]
			partyInfo.UserAccessoryID2 = dictInfo.UserAccessoryIds[1]
			partyInfo.UserAccessoryID3 = dictInfo.UserAccessoryIds[2]
			session.UpdateUserLiveParty(partyInfo)
		}
	}

	respBody := session.Finalize(GetData("saveDeckAll.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), respBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchLiveMusicSelect(ctx *gin.Context) {
	now := time.Now()
	year, month, day := now.Year(), now.Month(), now.Day()
	tomorrow := time.Date(year, month, day+1, 0, 0, 0, 0, now.Location()).Unix()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	liveDailyList := []model.LiveDaily{}
	err := MainEng.Table("m_live_daily").Where("weekday = ?", weekday).Cols("id,live_id").Find(&liveDailyList)
	CheckErr(err)
	for k := range liveDailyList {
		liveDailyList[k].EndAt = int(tomorrow)
		liveDailyList[k].RemainingPlayCount = 5
		liveDailyList[k].RemainingRecoveryCount = 10
	}

	signBody := GetData("fetchLiveMusicSelect.json")
	signBody, _ = sjson.Set(signBody, "weekday_state.weekday", weekday)
	signBody, _ = sjson.Set(signBody, "weekday_state.next_weekday_at", tomorrow)
	signBody, _ = sjson.Set(signBody, "live_daily_list", liveDailyList)
	session := serverdb.GetSession(UserID)
	signBody = session.Finalize(signBody, "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchLivePartners(ctx *gin.Context) {
	// a set of partners player (i.e. friends and others), then fetch the card for them
	// this set include the current user, so we can use our own cards.
	// currently only have current user
	// note that all card are available, but we need to use the filter functionality to actually get them to show up.
	partnerIDs := []int{}
	partnerIDs = append(partnerIDs, UserID)
	livePartners := []model.LiveStartLivePartner{}
	for _, partnerID := range partnerIDs {
		partner := model.LiveStartLivePartner{}
		partner.IsFriend = true
		serverdb.FetchDBProfile(partnerID, &partner)
		partnerCards := serverdb.FetchPartnerCards(partnerID) // model.UserCard
		for i := 1; i <= 7; i++ {
			partner.CardByCategory = append(partner.CardByCategory, i)
			partner.CardByCategory = append(partner.CardByCategory, model.PartnerCardInfo{})
		}
		for _, card := range partnerCards {
			for i := 1; i <= 7; i++ {
				if (card.LivePartnerCategories & (1 << i)) != 0 {
					partnerCardInfo := serverdb.GetPartnerCardFromUserCard(card)
					partner.CardByCategory[i*2-1] = partnerCardInfo
				}
			}
		}
		livePartners = append(livePartners, partner)
	}

	signBody := "{}"
	signBody, _ = sjson.Set(signBody, "partner_select_state.live_partners", livePartners)
	signBody, _ = sjson.Set(signBody, "partner_select_state.friend_count", len(livePartners))
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchLiveDeckSelect(ctx *gin.Context) {
	signBody := GetData("fetchLiveDeckSelect.json")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func LiveStart(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())
	req := model.LiveStartReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}
	session := serverdb.GetSession(UserID)

	session.UserStatus.LastLiveDifficultyID = req.LiveDifficultyID
	session.UserStatus.LatestLiveDeckID = req.DeckID

	// 保存请求包因为 /live/finish 接口的响应包里有部分字段不在该接口的请求包里
	// live is stored in db
	live := model.LiveState{}
	live.UserID = UserID
	live.PartnerUserID = req.PartnerUserID
	live.LiveID = time.Now().UnixNano()
	live.LiveType = 1 // not sure what this is
	live.IsPartnerFriend = true
	live.DeckID = req.DeckID
	live.CellID = req.CellID  // cell id send player to the correct place after playing, normal live don't have cell id.

	liveNotes := utils.ReadAllText(fmt.Sprintf("assets/stages/%d.json", req.LiveDifficultyID))
	if liveNotes == "" {
		panic("歌曲情报信息不存在！(song doesn't exist)")
	}

	if err := json.Unmarshal([]byte(liveNotes), &live.LiveStage); err != nil {
		panic(err)
	}

	if req.IsAutoPlay {
		for k := range live.LiveStage.LiveNotes {
			live.LiveStage.LiveNotes[k].AutoJudgeType = 30
		}
	}
	
	if req.PartnerUserID != 0 {
		live.LivePartnerCard = serverdb.GetPartnerCardFromUserCard(
			serverdb.GetUserCard(req.PartnerUserID, req.PartnerCardMasterID))
	}

	liveStartResp := session.Finalize(GetData("userModelDiff.json"), "user_model_diff")
	liveStartResp, _ = sjson.Set(liveStartResp, "live", live)
	if req.PartnerUserID == 0 {
		liveStartResp, _ = sjson.Set(liveStartResp, "live.live_partner_card", nil)
	}
	serverdb.SaveLiveState(live)
	resp := SignResp(ctx.GetString("ep"), liveStartResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func LiveFinish(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	// fmt.Println(reqBody.String())

	var cardMasterId, maxVolt, skillCount, appealCount int64
	liveFinishReq := gjson.Parse(reqBody.String())
	liveFinishReq.Get("live_score.card_stat_dict").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			volt := value.Get("got_voltage").Int()
			if volt > maxVolt {
				maxVolt = volt

				cardMasterId = value.Get("card_master_id").Int()
				skillCount = value.Get("skill_triggered_count").Int()
				appealCount = value.Get("appeal_count").Int()
			}
		}
		return true
	})

	session := serverdb.GetSession(UserID)

	mvpInfo := model.MvpInfo{
		CardMasterID:        cardMasterId,
		GetVoltage:          maxVolt,
		SkillTriggeredCount: skillCount,
		AppealCount:         appealCount,
	}

	exists, live := serverdb.LoadLiveState(UserID)
	if !exists {
		panic("live doesn't exists")
	}
	live.DeckID = session.UserStatus.LatestLiveDeckID
	live.LiveStage.LiveDifficultyID = session.UserStatus.LastLiveDifficultyID

	liveResult := model.LiveResultAchievementStatus{
		ClearCount:       1,
		GotVoltage:       liveFinishReq.Get("live_score.current_score").Int(),
		RemainingStamina: liveFinishReq.Get("live_score.remaining_stamina").Int(),
	}

	liveFinishResp := GetData("liveFinish.json")
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.live_difficulty_master_id", live.LiveStage.LiveDifficultyID)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.live_deck_id", live.DeckID)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.mvp", mvpInfo)
	if live.PartnerUserID == 0 {
		liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.partner", nil)
	} else {
		liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.partner", 
		session.GetOtherUserBasicProfile(live.PartnerUserID))
	}
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.live_result_achievement_status", liveResult)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.voltage", liveFinishReq.Get("live_score.current_score").Int())
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.last_best_voltage", liveFinishReq.Get("live_score.current_score").Int())
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.before_user_exp", session.UserStatus.Exp)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.gain_user_exp", 0)


	liveFinishResp = session.Finalize(liveFinishResp, "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), liveFinishResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

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
	// fmt.Println(reqData)

	saveReq := model.LiveSaveDeckReq{}
	err := json.Unmarshal([]byte(reqData.String()), &saveReq)
	if err != nil {
		panic(err)
	}
	// fmt.Println(saveReq)

	userLiveMvDeckInfo := model.UserLiveMvDeckInfo{
		LiveMasterID: saveReq.LiveMasterID,
	}

	memberInfoList := map[int]model.UserMemberInfo{}
	memberIds := map[int]int{}
	for k, v := range saveReq.MemberMasterIDByPos {
		if k%2 == 0 {
			memberId := saveReq.MemberMasterIDByPos[k+1]
			memberIds[v] = memberId

			switch v {
			case 1:
				userLiveMvDeckInfo.MemberMasterID1 = memberId
			case 2:
				userLiveMvDeckInfo.MemberMasterID2 = memberId
			case 3:
				userLiveMvDeckInfo.MemberMasterID3 = memberId
			case 4:
				userLiveMvDeckInfo.MemberMasterID4 = memberId
			case 5:
				userLiveMvDeckInfo.MemberMasterID5 = memberId
			case 6:
				userLiveMvDeckInfo.MemberMasterID6 = memberId
			case 7:
				userLiveMvDeckInfo.MemberMasterID7 = memberId
			case 8:
				userLiveMvDeckInfo.MemberMasterID8 = memberId
			case 9:
				userLiveMvDeckInfo.MemberMasterID9 = memberId
			case 10:
				userLiveMvDeckInfo.MemberMasterID10 = memberId
			case 11:
				userLiveMvDeckInfo.MemberMasterID11 = memberId
			case 12:
				userLiveMvDeckInfo.MemberMasterID12 = memberId
			}

			memberInfoList[v] = GetMemberInfo(memberId)
		}
	}
	// fmt.Println(memberIds)
	// fmt.Println(memberInfoList)

	suitIds := map[int]int{}
	for k, v := range saveReq.SuitMasterIDByPos {
		if k%2 == 0 {
			suitId := saveReq.SuitMasterIDByPos[k+1]
			suitIds[v] = suitId

			switch v {
			case 1:
				userLiveMvDeckInfo.SuitMasterID1 = suitId
			case 2:
				userLiveMvDeckInfo.SuitMasterID2 = suitId
			case 3:
				userLiveMvDeckInfo.SuitMasterID3 = suitId
			case 4:
				userLiveMvDeckInfo.SuitMasterID4 = suitId
			case 5:
				userLiveMvDeckInfo.SuitMasterID5 = suitId
			case 6:
				userLiveMvDeckInfo.SuitMasterID6 = suitId
			case 7:
				userLiveMvDeckInfo.SuitMasterID7 = suitId
			case 8:
				userLiveMvDeckInfo.SuitMasterID8 = suitId
			case 9:
				userLiveMvDeckInfo.SuitMasterID9 = suitId
			case 10:
				userLiveMvDeckInfo.SuitMasterID10 = suitId
			case 11:
				userLiveMvDeckInfo.SuitMasterID11 = suitId
			case 12:
				userLiveMvDeckInfo.SuitMasterID12 = suitId
			}
		}
	}
	// fmt.Println(suitIds)

	var newMemberInfoList []any
	for k, v := range saveReq.ViewStatusByPos {
		if k%2 == 0 {
			memberInfo := memberInfoList[v]
			memberInfo.ViewStatus = saveReq.ViewStatusByPos[k+1]

			newMemberInfoList = append(newMemberInfoList, memberInfo.MemberMasterID)
			newMemberInfoList = append(newMemberInfoList, memberInfo)
			// fmt.Printf("k => %d, v => %d, val => %d\n", k, v, saveReq.ViewStatusByPos[k+1])
		}
	}
	// fmt.Println(newMemberInfoList)

	var userLiveMvDeckCustomByID []any
	userLiveMvDeckCustomByID = append(userLiveMvDeckCustomByID, saveReq.LiveMasterID)
	userLiveMvDeckCustomByID = append(userLiveMvDeckCustomByID, userLiveMvDeckInfo)
	// fmt.Println(userLiveMvDeckCustomByID)

	session := serverdb.GetSession(UserID)
	signBody := GetData("liveMvSaveDeck.json")
	signBody = session.Finalize(signBody, "user_model")
	// signBody, _ = sjson.Set(signBody, "user_model.user_status", GetUserStatus())
	signBody, _ = sjson.Set(signBody, "user_model.user_live_mv_deck_custom_by_id", userLiveMvDeckCustomByID)
	signBody, _ = sjson.Set(signBody, "user_model.user_member_by_member_id", newMemberInfoList)

	resp := SignResp(ctx.GetString("ep"), string(signBody), config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SaveSuit(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)

	req := gjson.Parse(reqBody).Array()[0]
	deckId := req.Get("deck_id").Int()
	cardId := req.Get("card_index").Int()
	suitId := req.Get("suit_master_id").Int()

	deckIndex := deckId*2 - 1
	keyLiveDeck := fmt.Sprintf("user_live_deck_by_id.%d", deckIndex)
	// fmt.Println("keyLiveDeck:", keyLiveDeck)
	liveDeck := gjson.Parse(GetLiveDeckData()).Get(keyLiveDeck).String()
	// fmt.Println(liveDeck)
	keyLiveDeckInfo := fmt.Sprintf("suit_master_id_%d", cardId)
	liveDeck, _ = sjson.Set(liveDeck, keyLiveDeckInfo, suitId)
	// fmt.Println(liveDeck)

	var deckInfo model.UserLiveDeck
	if err := json.Unmarshal([]byte(liveDeck), &deckInfo); err != nil {
		panic(err)
	}

	SetLiveDeckData(keyLiveDeck, deckInfo)
	session := serverdb.GetSession(UserID)
	signBody := session.Finalize(GetData("saveSuit.json"), "user_model")
	signBody, _ = sjson.Set(signBody, "user_model.user_live_deck_by_id.0", deckId)
	signBody, _ = sjson.Set(signBody, "user_model.user_live_deck_by_id.1", deckInfo)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SaveDeck(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)

	req := gjson.Parse(reqBody).Array()[0]
	deckId := req.Get("deck_id")
	// fmt.Println("deckId:", deckId)

	position := req.Get("card_master_ids.0")
	cardMasterId := req.Get("card_master_ids.1")
	// fmt.Println("cardMasterId:", cardMasterId)

	var deckInfo, partyInfo string
	var oldCardMasterId int64
	var partyId int64
	var savePartyInfo model.UserLiveParty
	deckList := GetLiveDeckData()
	gjson.Parse(deckList).Get("user_live_deck_by_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() && value.Get("user_live_deck_id").String() == deckId.String() {
			deckInfo = value.String()
			// fmt.Println("deckInfo:", deckInfo)

			oldCardMasterId = gjson.Parse(deckInfo).Get("card_master_id_" + position.String()).Int()
			deckInfo, _ = sjson.Set(deckInfo, "card_master_id_"+position.String(), cardMasterId.Int())
			deckInfo, _ = sjson.Set(deckInfo, "suit_master_id_"+position.String(), cardMasterId.Int())
			// fmt.Println("New deckInfo:", deckInfo)

			SetLiveDeckData("user_live_deck_by_id."+key.String(), gjson.Parse(deckInfo).Value())

			return false
		}
		return true
	})
	gjson.Parse(deckList).Get("user_live_party_by_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() && (value.Get("party_id").String() == deckId.String()+"01" ||
			value.Get("party_id").String() == deckId.String()+"02" ||
			value.Get("party_id").String() == deckId.String()+"03") {
			value.ForEach(func(kk, vv gjson.Result) bool {
				if vv.Int() == oldCardMasterId {
					partyInfo = value.String()
					// fmt.Println("partyInfo:", partyInfo)

					partyInfo, _ = sjson.Set(partyInfo, kk.String(), cardMasterId.Int())
					// fmt.Println("New partyInfo:", partyInfo)

					newPartyInfo := gjson.Parse(partyInfo)
					partyId = newPartyInfo.Get("party_id").Int()

					roleIds := []int{}
					err := MainEng.Table("m_card").
						Where("id IN (?,?,?)", newPartyInfo.Get("card_master_id_1").Int(),
							newPartyInfo.Get("card_master_id_2").Int(),
							newPartyInfo.Get("card_master_id_3").Int()).
						Cols("role").Find(&roleIds)
					CheckErr(err)
					// fmt.Println("roleIds:", roleIds)

					partyIcon, partyName := GetPartyInfoByRoleIds(roleIds)
					realPartyName := GetRealPartyName(partyName)
					partyInfo, _ = sjson.Set(partyInfo, "name.dot_under_text", realPartyName)
					partyInfo, _ = sjson.Set(partyInfo, "icon_master_id", partyIcon)
					// fmt.Println("New partyInfo 2:", partyInfo)

					decoder := json.NewDecoder(strings.NewReader(partyInfo))
					decoder.UseNumber()
					err = decoder.Decode(&savePartyInfo)
					CheckErr(err)
					SetLiveDeckData("user_live_party_by_id."+key.String(), savePartyInfo)

					return false
				}
				return true
			})
		}
		return true
	})

	session := serverdb.GetSession(UserID)
	signBody := session.Finalize(GetData("SaveDeck.json"), "user_model")
	signBody, _ = sjson.Set(signBody, "user_model.user_live_deck_by_id.0", deckId.Int())
	signBody, _ = sjson.Set(signBody, "user_model.user_live_deck_by_id.1", gjson.Parse(deckInfo).Value())
	signBody, _ = sjson.Set(signBody, "user_model.user_live_party_by_id.0", partyId)
	signBody, _ = sjson.Set(signBody, "user_model.user_live_party_by_id.1", savePartyInfo)
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
