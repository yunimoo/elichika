package gacha

import (
	"elichika/klab"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"math/rand"

	"github.com/gin-gonic/gin"
)

var (
	cachedGroupWeight map[int]int64
)

func init() {
	cachedGroupWeight = make(map[int]int64)
}

// fetch from db with cache
func GetGroupWeight(groupID int) int64 {
	_, exists := cachedGroupWeight[groupID]
	if !exists {
		// fetch all from db
		groups := []model.GachaGroup{}
		err := serverdb.Engine.Table("s_gacha_group").Find(&groups)
		utils.CheckErr(err)
		for _, group := range groups {
			cachedGroupWeight[group.GroupMasterID] = group.GroupWeight
		}
	}
	return cachedGroupWeight[groupID]
}

// it's not too bad to call this function multiple time, but maybe it's better to have a function that return multiple
func ChooseRandomCard(cards []model.GachaCard) int {
	if len(cards) == 0 { // no card
		return 0
	}
	groups := map[int]([]int){}
	totalWeight := int64(0)
	for _, card := range cards {
		_, exists := groups[card.GroupMasterID]
		if !exists {
			totalWeight += GetGroupWeight(card.GroupMasterID)
			groups[card.GroupMasterID] = []int{}
		}
		groups[card.GroupMasterID] = append(groups[card.GroupMasterID], card.CardMasterID)
	}
	groupRand := rand.Int63n(totalWeight)
	for groupID, cardIDs := range groups {
		if GetGroupWeight(groupID) > groupRand { // this group
			return cardIDs[rand.Intn(len(cardIDs))]
		} else {
			groupRand -= GetGroupWeight(groupID)
		}
	}
	panic("this shouldn't happen")
}

func MakeResultCard(session *serverdb.Session, cardMasterID int, isGuaranteed bool) model.ResultCard {
	card := session.GetUserCard(cardMasterID)
	cardRarity := klab.CardRarityFromCardMasterID(cardMasterID)
	member := session.GetMember(klab.MemberMasterIDFromCardMasterID(cardMasterID))
	resultCard := model.ResultCard{
		GachaLotType:         1,
		CardMasterID:         cardMasterID,
		Level:                1,
		BeforeGrade:          card.Grade,
		AfterGrade:           card.Grade + 1,
		Content:              nil,
		LimitExceeded:        false,
		BeforeLoveLevelLimit: klab.BondLevelFromBondValue(member.LovePointLimit),
		AfterLoveLevelLimit:  0,
	}
	if isGuaranteed {
		resultCard.GachaLotType = 2
	}
	if resultCard.AfterGrade == 6 { // maxed out card
		resultCard.AfterGrade = 5
		resultCard.Content = &model.RewardByContent{
			ContentType:   13,
			ContentID:     1800,
			ContentAmount: 1,
		}
		// 30 20 10 for UR, SR, R
		for i := cardRarity; i > 10; i -= 10 {
			resultCard.Content.ContentAmount *= 5
		}
	} else {
		if resultCard.AfterGrade == 0 { // entirely new card
			resultCard.BeforeGrade = 0
		}
		card.Grade++ // new grade
		resultCard.AfterLoveLevelLimit = resultCard.BeforeLoveLevelLimit + cardRarity/10
		member.LovePointLimit = klab.BondRequiredTotal(resultCard.AfterLoveLevelLimit)
		// update the card and member
		session.UpdateUserCard(card)
		session.UpdateMember(member)
		// TODO: Add a card grade up trigger
	}
	return resultCard
}

func GachaFromGachaMasterID(gachaMasterID int) model.Gacha {
	gacha := model.Gacha{}
	exists, err := serverdb.Engine.Table("s_gacha").Where("gacha_master_id = ?", gachaMasterID).Get(&gacha)
	utils.CheckErrMustExist(err, exists)
	for _, gachaAppealMasterID := range gacha.DbGachaAppeals {
		gachaAppeal := model.GachaAppeal{}
		exists, err := serverdb.Engine.Table("s_gacha_appeal").Where("gacha_appeal_master_id = ?", gachaAppealMasterID).
			Get(&gachaAppeal)
		utils.CheckErrMustExist(err, exists)
		gacha.GachaAppeals = append(gacha.GachaAppeals, gachaAppeal)
	}
	for _, gachaDrawMasterID := range gacha.DbGachaDraws {
		gachaDraw := model.GachaDraw{}
		exists, err := serverdb.Engine.Table("s_gacha_draw").Where("gacha_draw_master_id = ?", gachaDrawMasterID).
			Get(&gachaDraw)
		utils.CheckErrMustExist(err, exists)
		gacha.GachaDraws = append(gacha.GachaDraws, gachaDraw)
	}
	return gacha
}

func HandleGacha(ctx *gin.Context, req model.GachaDrawReq) (model.Gacha, []model.ResultCard) {
	session := ctx.MustGet("session").(*serverdb.Session)
	draw := model.GachaDraw{}
	exists, err := serverdb.Engine.Table("s_gacha_draw").
		Where("gacha_draw_master_id = ?", req.GachaDrawMasterID).Get(&draw)
	utils.CheckErrMustExist(err, exists)
	gacha := GachaFromGachaMasterID(draw.GachaMasterID)
	cardPool := []model.GachaCard{}
	for _, group := range gacha.DbGachaGroups {
		groupPool := []model.GachaCard{}
		err := serverdb.Engine.Table("s_gacha_card").Where("group_master_id = ?", group).Find(&groupPool)
		utils.CheckErr(err)
		cardPool = append(cardPool, groupPool...)
		// allow 1 card to be in multiple group
	}
	ctx.Set("gacha_card_pool", cardPool)

	// TODO: gacha recovery and economy
	// for now just get this to work
	resultCards := []model.ResultCard{}
	for _, guranteeID := range draw.Guarantees {
		gachaGuarantee := model.GachaGuarantee{}
		exists, err := serverdb.Engine.Table("s_gacha_guarantee").
			Where("gacha_guarantee_master_id = ?", guranteeID).Get(&gachaGuarantee)
		utils.CheckErrMustExist(err, exists)
		cardMasterID := GuaranteeHandlers[gachaGuarantee.GuaranteeHandler](ctx, gachaGuarantee.GuaranteeParams)
		if cardMasterID == 0 {
			continue
		}
		resultCards = append(resultCards, MakeResultCard(session, cardMasterID, true))
	}
	for i := len(resultCards); i < draw.DrawCount; i++ {
		resultCards = append(resultCards, MakeResultCard(session, ChooseRandomCard(cardPool), false))
	}
	return gacha, resultCards
}
