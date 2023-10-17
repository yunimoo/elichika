package gamedata

import (
	"elichika/model"
	"elichika/utils"

	"xorm.io/xorm"
)

const GRADE_COUNT = 6 // from 0 - 5. Not hardcoding this is pretty messy

type AccessoryData struct {
	// from m_accessory
	MasterID int `xorm:"'id'"`
	// AccessoryName string  `xorm:"'name'"`
	// No int `xorm:"'accessory_no'"`
	// ThumbnailAssetPath int `xorm:"'thumbnail_asset_path'"`
	Type int `xorm:"'accessory_type'"`
	// MemberMasterID *int `xorm:"'member_master_id'"`
	RarityType int `xorm:"'rarity_type'"`
	Attribute  int `xorm:"attribute"`
	// Role int `xorm:"'role'"`
	MaxGrade int `xorm:"'max_grade'"`

	// from m_accessory_grade_up and m_accessory_melt
	Grade []struct {
		// m_accessory_grade_up
		Grade                 int  `xorm:"'grade'"`
		MaxLevel              int  `xorm:"'max_level'"`
		PassiveSkill1MasterID *int `xorm:"'accessory_passive_skill_1_master_id'"`
		PassiveSkill2MasterID *int `xorm:"'accessory_passive_skill_2_master_id'"`
		// m_accessory_melt
		MeltGroupMasterID int `xorm:"-"`
	} `xorm:"-"`

	// from m_accessory_level_up
	Level []struct {
		Level int `xorm:"level"`
		Exp   int `xorm:"exp"`
	} `xorm:"-"`

	// from m_accessory_rarity_up
	RarityUp struct {
		AfterAccessoryMasterID         int `xorm:"'after_accessory_master_id'"`
		AccessoryRarityUpGroupMasterID int `xorm:"'accessory_rarity_up_group_master_id'"`
	} `xorm:"-"`
}

type AccessoryMeltGroup struct {
	// from m_accessory_melt_group
	ID       int           `xorm:"'id'"`
	Resource model.Content `xorm:"extends"`
}

type AccessoryPassiveSkill struct {
	// from m_accessory_passive_skill
	ID int `xorm:"'id'"`
	// SkillType int `xorm:"'skill_type'"`
	// Rarity int `xorm:"'rarity'"`
	// IconAssetPath string `xorm:"'icon_asset_path'"`
	// ThumbnailAssetPath string `xorm:"'thumbnail_asset_path'"`
	MaxLevel int `xorm:"max_level"`
	// some random stuffs that doesn't need to be here

}

type AccessoryRarity struct {
	// from m_accessory_rarity_setting
	RarityType int `xorm:"'rarity_type'"`
	// Name string `xorm:"'name'"`

	// from m_accessory_frame_type
	FrameType int `xorm:"-"`

	Grade [GRADE_COUNT]struct {
		// from m_accessory_passive_skill_level_up_denominator
		SkillLevelUpDenominator []int
		// from m_accessory_passive_skill_level_up_setting
		SkillMaxLevel int
		// from m_accessory_grade_up_setting, this is the money used when consumed for grade up
		GradeUpMoney int
	} `xorm:"-"`

	// from m_accessory_passive_skill_level_up_plus_percent
	SkillLevel []struct {
		SkillLevel  int `xorm:"'skill_level'"`
		PlusPercent int `xorm:"'plus_percent'"`
	} `xorm:"-"`

	// from m_accessory_level_up_setting
	Level []struct {
		Level     int `xorm:"'level'"`
		PlusExp   int `xorm:"'plus_exp'"`
		GameMoney int `xorm:"'game_money'"`
	} `xorm:"-"`

	// from m_accessory_rarity_up_setting
	RarityUpMoney int `xorm:"-"`
}

type AccessoryLevelUpItem struct {
	// from m_accessory_level_up_item
	ID int `xorm:"'id'"`
	// Rarity int `xorm:"'rarity'"`
	// Attribute int `xorm:"'attribute'"`
	PlusExp   int `xorm:"'plus_exp'"`
	GameMoney int `xorm:"'game_money'"`
	// other unnecessary things
}

type AccessoryRarityUpGroup struct {
	// from m_accessory_rarity_up_group
	ID       int           `xorm:"'id'"`
	Resource model.Content `xorm:"extends"`
}

type Accessory struct {
	Accessory     map[int]AccessoryData
	MeltGroup     map[int]AccessoryMeltGroup
	PassiveSkill  map[int]AccessoryPassiveSkill
	Rarity        map[int]AccessoryRarity
	LevelUpItem   map[int]AccessoryLevelUpItem
	RarityUpGroup map[int]AccessoryRarityUpGroup
}

func (accessory *Accessory) Load(db, _ *xorm.Session) {
	var err error
	var exists bool

	accessory.Accessory = make(map[int]AccessoryData)
	allAccessories := []AccessoryData{}
	err = db.Table("m_accessory").Find(&allAccessories)
	utils.CheckErr(err)
	for _, item := range allAccessories {
		err = db.Table("m_accessory_grade_up").Where("accessory_master_id = ?", item.MasterID).
			OrderBy("grade").Find(&item.Grade)
		utils.CheckErr(err)
		if len(item.Grade) != GRADE_COUNT {
			panic("accessory must have correct grade count")
		}
		for g := range item.Grade {
			exists, err = db.Table("m_accessory_melt").
				Where("accessory_master_id = ? AND grade = ?", item.MasterID, item.Grade[g].Grade).
				Cols("accessory_melt_group_master_id").
				Get(&item.Grade[g].MeltGroupMasterID)
			utils.CheckErrMustExist(err, exists)
		}

		err = db.Table("m_accessory_level_up").Where("accessory_master_id = ?", item.MasterID).
			OrderBy("level").Find(&item.Level)
		utils.CheckErr(err)
		item.Level = append(item.Level[:1], item.Level...)

		exists, err = db.Table("m_accessory_rarity_up").Where("accessory_master_id = ?", item.MasterID).
			Get(&item.RarityUp)
		utils.CheckErr(err) // allowed to not exists
		accessory.Accessory[item.MasterID] = item

	}

	accessory.MeltGroup = make(map[int]AccessoryMeltGroup)
	allMeltGroups := []AccessoryMeltGroup{}
	err = db.Table("m_accessory_melt_group").Find(&allMeltGroups)
	utils.CheckErr(err)
	for _, group := range allMeltGroups {
		accessory.MeltGroup[group.ID] = group
	}

	accessory.PassiveSkill = make(map[int]AccessoryPassiveSkill)
	allPassiveSkills := []AccessoryPassiveSkill{}
	err = db.Table("m_accessory_passive_skill").Find(&allPassiveSkills)
	for _, skill := range allPassiveSkills {
		accessory.PassiveSkill[skill.ID] = skill
	}

	accessory.Rarity = make(map[int]AccessoryRarity)
	allRarities := []AccessoryRarity{}
	err = db.Table("m_accessory_rarity_setting").Find(&allRarities)
	utils.CheckErr(err)

	for _, rarity := range allRarities {
		exists, err = db.Table("m_accessory_frame_type").Where("rarity_type = ?", rarity.RarityType).
			Cols("frame_type").Get(&rarity.FrameType)
		utils.CheckErrMustExist(err, exists)

		for grade := range rarity.Grade {
			err = db.Table("m_accessory_passive_skill_level_up_denominator").
				Where("rarity = ? AND grade = ?", rarity.RarityType, grade).OrderBy("skill_level").Cols("denominator").
				Find(&rarity.Grade[grade].SkillLevelUpDenominator)
			utils.CheckErr(err)
			rarity.Grade[grade].SkillLevelUpDenominator =
				append(rarity.Grade[grade].SkillLevelUpDenominator[:1], rarity.Grade[grade].SkillLevelUpDenominator...)
			exists, err = db.Table("m_accessory_passive_skill_level_up_setting").
				Where("rarity = ? AND grade = ?", rarity.RarityType, grade).
				Cols("max_level").Get(&rarity.Grade[grade].SkillMaxLevel)
			utils.CheckErrMustExist(err, exists)
			exists, err = db.Table("m_accessory_grade_up_setting").
				Where("rarity = ? AND grade = ?", rarity.RarityType, grade).
				Cols("game_money").Get(&rarity.Grade[grade].GradeUpMoney)
			utils.CheckErrMustExist(err, exists)
		}

		err = db.Table("m_accessory_passive_skill_level_up_plus_percent").Where("rarity = ?", rarity.RarityType).
			Find(&rarity.SkillLevel)
		utils.CheckErr(err)
		rarity.SkillLevel = append(rarity.SkillLevel[:1], rarity.SkillLevel...)

		err = db.Table("m_accessory_level_up_setting").Where("rarity = ?", rarity.RarityType).
			Find(&rarity.Level)
		utils.CheckErr(err)
		rarity.Level = append(rarity.Level[:1], rarity.Level...)

		exists, err = db.Table("m_accessory_rarity_up_setting").Where("rarity = ?", rarity.RarityType).
			Cols("game_money").Get(&rarity.RarityUpMoney)
		utils.CheckErr(err)

		accessory.Rarity[rarity.RarityType] = rarity
	}

	accessory.LevelUpItem = make(map[int]AccessoryLevelUpItem)
	allItems := []AccessoryLevelUpItem{}
	err = db.Table("m_accessory_level_up_item").Find(&allItems)
	utils.CheckErr(err)
	for _, item := range allItems {
		accessory.LevelUpItem[item.ID] = item
	}

	accessory.RarityUpGroup = make(map[int]AccessoryRarityUpGroup)
	allGroups := []AccessoryRarityUpGroup{}
	err = db.Table("m_accessory_rarity_up_group").Find(&allGroups)
	utils.CheckErr(err)
	for _, group := range allGroups {
		accessory.RarityUpGroup[group.ID] = group
	}

}
