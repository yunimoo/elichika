package serverdb

import (
	"elichika/model"

	"fmt"
)

// fetch a card, use the value in diff is present, otherwise fetch from db
func (session *Session) GetUserCard(cardMasterID int) model.UserCard {
	card, exist := session.CardDiffs[cardMasterID]
	if exist {
		return card
	}
	card = model.UserCard{}
	exists, err := Engine.Table("s_user_card").
		Where("user_id = ? AND card_master_id = ?", session.UserStatus.UserID, cardMasterID).Get(&card)
	if err != nil {
		panic(err)
	}

	if !exists {
		panic("no user card")
	}
	return card
}

func (session *Session) GetAllCards() []model.UserCard {
	var cards []model.UserCard
	err := Engine.Table("s_user_card").Where("user_id = ?", session.UserStatus.UserID).Find(&cards)
	if err != nil {
		panic(err)
	}
	return cards
}

func (session *Session) UpdateUserCard(card model.UserCard) {
	session.CardDiffs[card.CardMasterID] = card
}

func (session *Session) FinalizeCardDiffs() []any {
	userCardByCardID := []any{}
	for cardMasterID, card := range session.CardDiffs {
		userCardByCardID = append(userCardByCardID, cardMasterID)
		userCardByCardID = append(userCardByCardID, card)
		// .AllCols() is necessary to all field
		affected, err := Engine.Table("s_user_card").
			Where("user_id = ? AND card_master_id = ?", session.UserStatus.UserID, cardMasterID).AllCols().Update(card)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userCardByCardID
}

// insert all the cards
func (session *Session) InsertCards(cards []model.UserCard) {
	affected, err := Engine.Table("s_user_card").Insert(&cards)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", affected, " cards")
}
