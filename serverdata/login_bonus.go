package serverdata

import (
	// "elichika/config"
	"elichika/enum"
	"elichika/model"
	"elichika/utils"

	"xorm.io/xorm"
)

func InitialiseLoginBonus(session *xorm.Session) {
	const BeginnerLoginBonusId = 1000001
	const NormalLoginBonusId = 1000002
	const BirthDayLoginBonusId = 1000003
	// normal login bonus
	loginBonuses := []model.LoginBonus{}
	loginBonusRewardDays := []model.LoginBonusRewardDay{}
	loginBonusRewardContents := []model.LoginBonusRewardContent{}

	loginBonuses = append(loginBonuses, model.LoginBonus{
		LoginBonusId:   NormalLoginBonusId,
		LoginBonusType: enum.LoginBonusTypeNormal,
		BackgroundId:   100100700,
		WhiteboardTextureAsset: &model.TextureStruktur{
			V: "/4n",
		},
		StartAt:                 0,
		EndAt:                   1<<31 - 1,
		LoginBonusHandler:       "normal_login_bonus",
		LoginBonusHandlerConfig: "",
	})
	for day := 1; day <= 7; day++ {
		loginBonusRewardDays = append(loginBonusRewardDays, model.LoginBonusRewardDay{
			LoginBonusId: NormalLoginBonusId,
			Day:          day,
			ContentGrade: enum.LoginBonusContentGradeNormal,
		})
	}
	loginBonusRewardContents = append(loginBonusRewardContents, model.LoginBonusRewardContent{
		LoginBonusId: NormalLoginBonusId,
		Day:          1,
		Content: model.Content{
			ContentType:   enum.ContentTypeRecoveryLp,
			ContentID:     1300,
			ContentAmount: 1,
		},
	})
	loginBonusRewardContents = append(loginBonusRewardContents, model.LoginBonusRewardContent{
		LoginBonusId: NormalLoginBonusId,
		Day:          2,
		Content: model.Content{
			ContentType:   enum.ContentTypeSnsCoin,
			ContentID:     0,
			ContentAmount: 10,
		},
	})
	loginBonusRewardContents = append(loginBonusRewardContents, model.LoginBonusRewardContent{
		LoginBonusId: NormalLoginBonusId,
		Day:          3,
		Content: model.Content{
			ContentType:   enum.ContentTypeRecoveryAp,
			ContentID:     2200,
			ContentAmount: 1,
		},
	})
	loginBonusRewardContents = append(loginBonusRewardContents, model.LoginBonusRewardContent{
		LoginBonusId: NormalLoginBonusId,
		Day:          4,
		Content: model.Content{
			ContentType:   enum.ContentTypeSnsCoin,
			ContentID:     0,
			ContentAmount: 20,
		},
	})
	loginBonusRewardContents = append(loginBonusRewardContents, model.LoginBonusRewardContent{
		LoginBonusId: NormalLoginBonusId,
		Day:          5,
		Content: model.Content{
			ContentType:   enum.ContentTypeRecoveryLp,
			ContentID:     1300,
			ContentAmount: 1,
		},
	})
	loginBonusRewardContents = append(loginBonusRewardContents, model.LoginBonusRewardContent{
		LoginBonusId: NormalLoginBonusId,
		Day:          6,
		Content: model.Content{
			ContentType:   enum.ContentTypeStoryEventUnlock,
			ContentID:     17001,
			ContentAmount: 1,
		},
	})
	loginBonusRewardContents = append(loginBonusRewardContents, model.LoginBonusRewardContent{
		LoginBonusId: NormalLoginBonusId,
		Day:          7,
		Content: model.Content{
			ContentType:   enum.ContentTypeSnsCoin,
			ContentID:     0,
			ContentAmount: 30,
		},
	})

	// birthday login bonus
	loginBonuses = append(loginBonuses, model.LoginBonus{
		LoginBonusId:            1000003,
		LoginBonusType:          enum.LoginBonusTypeBirthday,
		StartAt:                 0,
		EndAt:                   1<<31 - 1,
		LoginBonusHandler:       "birthday_login_bonus",
		LoginBonusHandlerConfig: "random", // can be set to latest to use the latest pair
	})

	for _, loginBonus := range loginBonuses {
		_, err := session.Table("s_login_bonus").Insert(loginBonus)
		utils.CheckErr(err)
	}

	for _, day := range loginBonusRewardDays {
		_, err := session.Table("s_login_bonus_reward_day").Insert(day)
		utils.CheckErr(err)
	}

	for _, content := range loginBonusRewardContents {
		_, err := session.Table("s_login_bonus_reward_content").Insert(content)
		utils.CheckErr(err)
	}

}
