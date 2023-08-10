package handler

import (
	"elichika/config"
	"elichika/enum"
	"elichika/klab"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"xorm.io/xorm"
)

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
	session := serverdb.GetSession(ctx, UserID)
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
		if len(partnerCards) == 0 {
			continue
		}
		for _, card := range partnerCards {
			for i := 1; i <= 7; i++ {
				if (card.LivePartnerCategories & (1 << i)) != 0 {
					partnerCardInfo := serverdb.GetPartnerCardFromUserCard(card)
					partner.CardByCategory = append(partner.CardByCategory, i)
					partner.CardByCategory = append(partner.CardByCategory, partnerCardInfo)
				}
			}
		}
		livePartners = append(livePartners, partner)
	}

	signBody := "{}"
	signBody, _ = sjson.Set(signBody, "partner_select_state.live_partners", livePartners)
	signBody, _ = sjson.Set(signBody, "partner_select_state.friend_count", len(livePartners))
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)
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
	session := serverdb.GetSession(ctx, UserID)

	session.UserStatus.LastLiveDifficultyID = req.LiveDifficultyID
	session.UserStatus.LatestLiveDeckID = req.DeckID

	// 保存请求包因为 /live/finish 接口的响应包里有部分字段不在该接口的请求包里
	// live is stored in db
	liveState := model.LiveState{}
	liveState.UserID = UserID
	liveState.PartnerUserID = req.PartnerUserID
	liveState.LiveID = time.Now().UnixNano()
	liveState.LiveType = 1 // not sure what this is
	liveState.IsPartnerFriend = true
	liveState.DeckID = req.DeckID
	liveState.CellID = req.CellID // cell id send player to the correct place after playing, normal live don't have cell id.

	liveNotes := utils.ReadAllText(fmt.Sprintf("assets/stages/%d.json", req.LiveDifficultyID))
	if liveNotes == "" {
		panic("歌曲情报信息不存在！(song doesn't exist)")
	}

	if err := json.Unmarshal([]byte(liveNotes), &liveState.LiveStage); err != nil {
		panic(err)
	}

	if req.IsAutoPlay {
		for k := range liveState.LiveStage.LiveNotes {
			liveState.LiveStage.LiveNotes[k].AutoJudgeType = 30
		}
	}

	if req.PartnerUserID != 0 {
		liveState.LivePartnerCard = serverdb.GetPartnerCardFromUserCard(
			serverdb.GetUserCard(req.PartnerUserID, req.PartnerCardMasterID))
	}

	liveStartResp := session.Finalize(GetData("userModelDiff.json"), "user_model_diff")
	liveStartResp, _ = sjson.Set(liveStartResp, "live", liveState)
	if req.PartnerUserID == 0 {
		liveStartResp, _ = sjson.Set(liveStartResp, "live.live_partner_card", nil)
	}
	serverdb.SaveLiveState(liveState)
	resp := SignResp(ctx.GetString("ep"), liveStartResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func LiveFinish(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type LiveFinishReq struct {
		LiveID           int64 `json:"live_id"`
		LiveFinishStatus int   `json:"live_finish_status"`
		LiveScore        struct {
			StartInfo                  any   `json:"start_info"`
			FinishInfo                 any   `json:"finish_info"`
			ResultDict                 []any `json:"result_dict"`
			WaveStatDict               []any `json:"wave_stat_dict"`
			TurnStatDict               []any `json:"turn_stat_dict"`
			CardStatDict               []any `json:"card_stat_dict"`
			TargetScore                int   `json:"target_score"`
			CurrentScore               int   `json:"current_score"`
			ComboCount                 int   `json:"combo_count"`
			ChangeSquadCount           int   `json:"change_squad_count"`
			HighestComboCount          int   `json:"highest_combo_count"`
			RemainingStamina           int   `json:"remaining_stamina"`
			IsPerfectLive              bool  `json:"is_perfect_live"`
			IsPerfectFullCombo         bool  `json:"is_perfect_full_combo"`
			UseVoltageActiveSkillCount int   `json:"use_voltage_active_skill_count"`
			UseHealActiveSkillCount    int   `json:"use_heal_active_skill_count"`
			UseDebufActiveSkillCount   int   `json:"use_debuf_active_skill_count"`
			UseBufActiveSkillCount     int   `json:"use_buf_active_skill_count"`
			UseSpSkillCount            int   `json:"use_sp_skill_count"`
			CompleteAppealChanceCount  int   `json:"complete_appeal_chance_count"`
			TriggerCriticalCount       int   `json:"triggered_critical_count"`
			LivePower                  int   `json:"live_power"`
			SpSkillScoreList           []int `json:"sp_skill_score_list"`
		} `json:"live_score"`
		ResumeFinishInfo any `json:"resume_finish_info"`
		RoomID           int `json:"room_id"`
	}
	req := LiveFinishReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	CheckErr(err)

	exists, liveState := serverdb.LoadLiveState(UserID)
	if !exists {
		panic("live doesn't exists")
	}

	db := ctx.MustGet("masterdata.db").(*xorm.Engine)
	session := serverdb.GetSession(ctx, UserID)

	liveState.DeckID = session.UserStatus.LatestLiveDeckID
	liveState.LiveStage.LiveDifficultyID = session.UserStatus.LastLiveDifficultyID

	// record this live
	liveRecord := session.GetLiveDifficultyRecord(session.UserStatus.LastLiveDifficultyID)
	lastPlayDeck := session.BuildLastPlayLiveDifficultyDeck(liveState.DeckID, liveState.LiveStage.LiveDifficultyID)

	liveRecord.PlayCount++
	lastPlayDeck.IsCleared = req.LiveFinishStatus == enum.LiveFinishStatusCleared
	if lastPlayDeck.IsCleared {
		// update clear record
		liveRecord.ClearCount++
		// and award items
		type LiveDifficultyMission struct {
			Position    int
			TargetValue int
			Reward      model.RewardByContent `xorm:"extends"`
		}
		missions := []LiveDifficultyMission{}

		db.Table("m_live_difficulty_mission").Where("live_difficulty_master_id = ?", session.UserStatus.LastLiveDifficultyID).
			Find(&missions)
		// TODO: set live_result_achievements
		for _, mission := range missions { // TODO: maybe work on JSON.
			if (mission.TargetValue == 1) || (mission.TargetValue <= req.LiveScore.CurrentScore) {
				if mission.Position == 1 {
					if liveRecord.ClearedDifficultyAchievement1 == nil {
						liveRecord.ClearedDifficultyAchievement1 = new(int)
						*liveRecord.ClearedDifficultyAchievement1 = 1
						session.AddRewardContent(mission.Reward)
					}
				} else if mission.Position == 2 {
					if liveRecord.ClearedDifficultyAchievement2 == nil {
						liveRecord.ClearedDifficultyAchievement2 = new(int)
						*liveRecord.ClearedDifficultyAchievement2 = 2
						session.AddRewardContent(mission.Reward)
					}
				} else if mission.Position == 3 {
					if liveRecord.ClearedDifficultyAchievement3 == nil {
						liveRecord.ClearedDifficultyAchievement3 = new(int)
						*liveRecord.ClearedDifficultyAchievement3 = 3
						session.AddRewardContent(mission.Reward)
					}
				}
			}
		}
	}

	mvpInfo := model.MvpInfo{
		CardMasterID:        0,
		GetVoltage:          -1,
		SkillTriggeredCount: 0,
		AppealCount:         0,
	}

	type RewardLovePoint struct {
		RewardLovePoint int `json:"reward_love_point"`
	}
	rewardLovePoint := RewardLovePoint{}
	if lastPlayDeck.IsCleared {
		_, err := db.Table("m_live_difficulty").Where("live_difficulty_id = ?", liveState.LiveStage.LiveDifficultyID).
			Cols("reward_base_love_point").Get(&rewardLovePoint.RewardLovePoint)
		CheckErr(err)
	}
	memberLoveStatuses := []any{}
	bondCardPosition := make(map[int]int)
	for i := 1; i <= 17; i += 2 {
		cardStat := req.LiveScore.CardStatDict[i].(map[string]any)
		cardMasterID := int(cardStat["card_master_id"].(float64))
		gotVoltage := int(cardStat["got_voltage"].(float64))
		skillCount := int(cardStat["skill_triggered_count"].(float64))
		appealCount := int(cardStat["appeal_count"].(float64))

		// calculate mvp
		if gotVoltage > mvpInfo.GetVoltage {
			mvpInfo.GetVoltage = gotVoltage
			mvpInfo.CardMasterID = cardMasterID
			mvpInfo.SkillTriggeredCount = skillCount
			mvpInfo.AppealCount = appealCount
		}

		// update card stat and member bond if cleared
		if lastPlayDeck.IsCleared {
			card := session.GetUserCard(cardMasterID)
			card.LiveJoinCount++
			card.ActiveSkillPlayCount += skillCount
			session.UpdateUserCard(card)
			// update member bond
			// TODO: figure out how to add bond for center(s)
			memberMasterID := klab.MemberMasterIDFromCardMasterID(cardMasterID)
			session.AddLovePoint(memberMasterID, rewardLovePoint.RewardLovePoint)
			pos, exists := bondCardPosition[memberMasterID]
			// only use 1 item or an idol might be shown multiple time
			if !exists {
				memberLoveStatuses = append(memberLoveStatuses, cardMasterID)
				bondCardPosition[memberMasterID] = len(memberLoveStatuses)
				memberLoveStatuses = append(memberLoveStatuses, rewardLovePoint)
			} else {
				memberLoveStatuses[pos] = RewardLovePoint{
					RewardLovePoint: memberLoveStatuses[pos].(RewardLovePoint).RewardLovePoint + rewardLovePoint.RewardLovePoint}
			}
		}
	}

	lastPlayDeck.Voltage = req.LiveScore.CurrentScore
	oldMaxScore := liveRecord.MaxScore
	if liveRecord.MaxScore < req.LiveScore.CurrentScore { // if new high score
		liveRecord.MaxScore = req.LiveScore.CurrentScore
	}
	if liveRecord.MaxCombo < req.LiveScore.HighestComboCount {
		liveRecord.MaxCombo = req.LiveScore.HighestComboCount
	}

	liveResult := model.LiveResultAchievementStatus{
		ClearCount:       1,
		GotVoltage:       int64(req.LiveScore.CurrentScore),
		RemainingStamina: int64(req.LiveScore.RemainingStamina),
	}

	liveFinishResp := GetData("liveFinish.json")
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.live_difficulty_master_id", liveState.LiveStage.LiveDifficultyID)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.live_deck_id", liveState.DeckID)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.mvp", mvpInfo)
	if liveState.PartnerUserID == 0 {
		liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.partner", nil)
	} else {
		liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.partner",
			session.GetOtherUserBasicProfile(liveState.PartnerUserID))
	}
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.live_result_achievement_status", liveResult)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.voltage", req.LiveScore.CurrentScore)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.last_best_voltage", oldMaxScore)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.before_user_exp", session.UserStatus.Exp)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.gain_user_exp", 0)
	liveFinishResp, _ = sjson.Set(liveFinishResp, "live_result.member_love_statuses", memberLoveStatuses)

	session.InsertOrUpdateLiveDifficultyRecord(liveRecord)
	session.SetLastPlayLiveDifficultyDeck(lastPlayDeck)
	liveFinishResp = session.Finalize(liveFinishResp, "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), liveFinishResp, config.SessionKey)
	fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
