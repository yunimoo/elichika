package model

type TriggerBasic struct {
	TriggerID       int64   `json:"trigger_id"` // use nano timestamp
	InfoTriggerType int     `json:"info_trigger_type"`
	LimitAt         *int64  `json:"limit_at"` // seems like some sort of timed timestamp, probably for event popup
	Description     *string `json:"description"`
	ParamInt        int     `json:"param_int"`
}

type TriggerCardGradeUp struct {
	TriggerID            int64 `json:"trigger_id"` // use nano timestamp
	CardMasterID         int   `json:"card_master_id"`
	BeforeLoveLevelLimit int   `json:"before_love_level_limit"`
	AfterLoveLevelLimit  int   `json:"after_love_level_limit"`
}

type TriggerReadReq struct {
	TriggerID int64 `json:"trigger_id"` // same for all trigger, for now
}
