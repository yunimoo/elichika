package serverdb

// give reward based on content type, content id, and content amount
// depend on the actual handler, negative amount can lead to a decrease as well

import (
	"elichika/model"

	"fmt"
)

var (
	RewardHandler map[int]func(*Session, int, int)
)

func (session *Session) AddRewardContent(reward model.RewardByContent) {
	RewardHandler[reward.ContentType](session, reward.ContentID, reward.ContentAmount)
}

func init() {
	RewardHandler = make(map[int]func(*Session, int, int))
	RewardHandler[1] = CurrencyRewardHandler
	RewardHandler[7] = SuitRewardHandler
}

func SuitRewardHandler(session *Session, suitMasterID, _ int) {
	session.InsertUserSuit(model.UserSuit{
		UserID:       session.UserStatus.UserID,
		SuitMasterID: suitMasterID,
		IsNew:        true})
}

func CurrencyRewardHandler(session *Session, currencyType, amount int) {
	fmt.Println("TODO: Add economy")
}
