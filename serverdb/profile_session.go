package serverdb

import (
	"elichika/model"

	"strconv"
	"strings"
	"encoding/json"

	"github.com/tidwall/gjson"
)

func (session *Session) FetchPartnerCards(otherUserID int) []model.CardInfo {
	partnerCards := []model.CardInfo{}
	err := Engine.Table("s_user_card").
		Where("user_id = ? AND live_partner_categories != 0", otherUserID).
		Find(&partnerCards)
	if err != nil {
			panic(err)
	}
	return partnerCards

}

// fetch profile of another user, from session.UserStatus.UserID's perspective
// it's possible that otherUserID == session.UserStatus.UserID
func (session *Session) FetchProfile(otherUserID int) model.Profile {
	profile := model.Profile{}

	exists, err := Engine.Table("s_user_info").Where("user_id = ?", otherUserID).Get(&profile)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("user doesn't exist")
	}
	// recommend card

	recommendCard := model.CardInfo{}
	exists, err = Engine.Table("s_user_card").Where("user_id = ? AND card_master_id = ?", 
		otherUserID, profile.ProfileInfo.BasicInfo.RecommendCardMasterID).
		Get(&recommendCard)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("card doesn't exist")
	}
	profile.ProfileInfo.BasicInfo.RecommendCardLevel = recommendCard.Level
	profile.ProfileInfo.BasicInfo.IsRecommendCardImageAwaken = recommendCard.IsAwakeningImage
	profile.ProfileInfo.BasicInfo.IsRecommendCardAllTrainingActivated = recommendCard.IsAllTrainingActivated

	// friend system, not implemented
	profile.ProfileInfo.BasicInfo.FriendApprovedAt = nil
	profile.ProfileInfo.BasicInfo.RequestStatus = 1
	profile.ProfileInfo.BasicInfo.IsRequestPending = false

	// other user's members
	members := []model.UserMemberInfo{}
	err = Engine.Table("s_user_member").Where("user_id = ?", otherUserID).OrderBy("love_point DESC").Find(&members)
	if err != nil {
		panic(err)
	}
	profile.ProfileInfo.TotalLovePoint = 0
	for _, member := range members {
		profile.ProfileInfo.TotalLovePoint += member.LovePoint
	}
	for i:=0; i < 3; i++ {
		profile.ProfileInfo.LoveMembers[i].MemberMasterID = members[i].MemberMasterID
		profile.ProfileInfo.LoveMembers[i].LovePoint = members[i].LovePoint
	}

	// need to return this in order, so make the array then write to it
	// this also prevent the game from freezing if 2 cards or more have the same bit
	for i := 0; i < 7; i++ {
		profile.GuestInfo.LivePartnersCards = append(profile.GuestInfo.LivePartnersCards, model.LivePartnerCard{})
		profile.GuestInfo.LivePartnersCards[i].LivePartnerCategoryMasterId = i + 1
	}

	partnerCards := session.FetchPartnerCards(otherUserID)
	for _, card := range partnerCards {
		// this is a just convention, might be better to check db
		memberId := (card.CardMasterID / 10000) % 1000
		
		partnerCard := model.PartnerCardInfo{}

		jsonByte, err := json.Marshal(card)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(jsonByte, &partnerCard)
		if err != nil {
			panic(err)
		}
		for _, member := range members {
			if member.MemberMasterID == memberId {
				partnerCard.LoveLevel = member.LoveLevel
				break
			}
		}
		partnerCard.PassiveSkillLevels = []int{}
		partnerCard.PassiveSkillLevels = append(partnerCard.PassiveSkillLevels, card.PassiveSkillALevel)
		partnerCard.PassiveSkillLevels = append(partnerCard.PassiveSkillLevels, card.PassiveSkillBLevel)
		partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill1ID)
		partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill2ID)
		partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill3ID)
		partnerCard.AdditionalPassiveSkillIds = append(partnerCard.AdditionalPassiveSkillIds, card.AdditionalPassiveSkill4ID)
		Engine.Table("s_user_member_love_panel").
			Where("user_id = ? AND member_id = ?", otherUserID, memberId).
			Cols("member_love_panel_cell_id").Find(&partnerCard.MemberLovePanels)
		livePartner := model.LivePartnerCard{}
		livePartner.PartnerCard = partnerCard
		for i := 1; i <= 7; i++ {
			if (card.LivePartnerCategories & (1 << i)) != 0 {
				profile.GuestInfo.LivePartnersCards[i - 1].PartnerCard = partnerCard
			}
		}
	}
	liveStats := model.DBUserProfileLiveStats{}
	_, err = Engine.Table("s_user_info").Where("user_id = ?", otherUserID).Get(&liveStats)
	if err != nil {
		panic(err)
	}
	jsonByte, err := json.Marshal(liveStats)
	if err != nil {
		panic(err)
	}
	gjson.Parse(string(jsonByte)).ForEach( func (key, value gjson.Result) bool {
		songRarity, _ := strconv.Atoi(key.String()[len(key.String()) - 2:])
		if strings.Contains(key.String(), "LivePlayCount") {
			profile.PlayInfo.LivePlayCount = append(profile.PlayInfo.LivePlayCount, songRarity)
			profile.PlayInfo.LivePlayCount = append(profile.PlayInfo.LivePlayCount, int(value.Int()))
		} else if strings.Contains(key.String(), "LiveClearCount") {
			profile.PlayInfo.LiveClearCount = append(profile.PlayInfo.LiveClearCount, songRarity)
			profile.PlayInfo.LiveClearCount = append(profile.PlayInfo.LiveClearCount, int(value.Int()))
		}
		return true
	})

	Engine.Table("s_user_card").Where("user_id = ?", otherUserID).
	OrderBy("live_join_count DESC").Limit(3).Find(&profile.PlayInfo.JoinedLiveCardRanking)
	Engine.Table("s_user_card").Where("user_id = ?", otherUserID).
	OrderBy("active_skill_play_count DESC").Limit(3).Find(&profile.PlayInfo.PlaySkillCardRanking)

	// can get from members[] to save sql
	Engine.Table("s_user_member").Where("user_id = ?", otherUserID).Find(&profile.MemberInfo.UserMembers)
	return profile
}