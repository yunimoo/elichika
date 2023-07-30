package serverdb

import (
	"elichika/model"

	"fmt"
	"encoding/json"

	"github.com/tidwall/gjson"
)

// fetch a card, use the value in diff is present, otherwise fetch from db
func (session *Session) GetCard(cardMasterId int) model.CardInfo {
	card, exist := session.CardDiffs[cardMasterId]
	if exist {
		return card
	}
	card = model.CardInfo{}
	exists, err := Engine.Table("s_user_card").Where("user_id = ? AND card_master_id = ?", session.UserInfo.UserId, cardMasterId).Get(&card)
	if err != nil {
		panic(err)
	}

	// if not in db then fetch from json
	if !exists {
		// insert this card, from json for now
		cardData := DbGetUserData("userCard.json")
		gjson.Parse(cardData).Get("user_card_by_card_id").ForEach(func(key, value gjson.Result) bool {
			if value.IsObject() {
				if err := json.Unmarshal([]byte(value.String()), &card); err != nil {
					panic(err)
				}
				if card.CardMasterID == cardMasterId {
					exists = true
					return false
				}
			}
			return true
		})
		if !exists {
			panic("Card doesn't exist")
		}
		card.UserId = session.UserInfo.UserId
		fmt.Println("Insert new card, user_id: ", card.UserId, ", card_master_id: ", cardMasterId)
		_, err := Engine.Table("s_user_card").Insert(&card)
		if err != nil {
			panic(err)
		}
	}
	return card
}

func (session *Session) GetAllCards() []model.CardInfo {
	var cards[] model.CardInfo
	err := Engine.Table("s_user_card").Where("user_id = ?", session.UserInfo.UserId).Find(&cards)
	if err != nil {
		panic(err)
	}
	return cards
}

func (session *Session) UpdateCard(card model.CardInfo) {
	session.CardDiffs[card.CardMasterID] = card
}

func (session *Session) FinalizeCardDiffs() []any {
	userCardByCardId := []any{}
	for cardMasterId, card := range session.CardDiffs {
		userCardByCardId = append(userCardByCardId, cardMasterId)
		userCardByCardId = append(userCardByCardId, card)
		// .AllCols() is necessary to all field
		affected, err := Engine.Table("s_user_card").
			Where("user_id = ? AND card_master_id = ?", session.UserInfo.UserId, cardMasterId).AllCols().Update(card)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userCardByCardId
}

// insert all the cards
func (session *Session) InsertCards(cards []model.CardInfo) {
	affected, err := Engine.Table("s_user_card").Insert(&cards)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", affected, " cards")
}
