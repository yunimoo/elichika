package model

type UserInfo struct {
	UserId int `xorm:"pk" json:"-"`
	Name   struct {
		DotUnderText string `xorm:"name" json:"dot_under_text"`
	} `xorm:"extends" json:"name"` // player name
	Nickname struct {
		DotUnderText string `xorm:"nickname" json:"dot_under_text"`
	} `xorm:"extends" json:"nickname"` // nickname in story
	LastLoginAt int `json:"last_login_at"` // in unix second
	Rank        int `json:"rank"`          // rank
	Exp         int `json:"exp"`           // total exp
	Message     struct {
		DotUnderText string `xorm:"message" json:"dot_under_text"`
	} `xorm:"extends" json:"message"` // introduction message
	RecommendCardMasterId                     int    `json:"recommend_card_master_id"` // featured / partner card
	MaxFriendNum                              int    `json:"max_friend_num"`
	LivePointFullAt                           int    `json:"live_point_full_at"`                              // in unix second
	LivePointBroken                           int    `json:"live_point_broken"`                               // max LP
	LivePointSubscriptionRecoveryDailyCount   int    `json:"live_point_subscription_recovery_daily_count"`    // membership LP 1 or 0?
	LivePointSubscriptionRecoveryDailyResetAt int    `json:"live_point_subscription_recovery_daily_reset_at"` // in unix second
	ActivityPointCount                        int    `json:"activity_point_count"`
	ActivityPointResetAt                      int    `json:"activity_point_reset_at"`                        // in unix time
	ActivityPointPaymentRecoveryDailyCount    int    `json:"activity_point_payment_recovery_daily_count"`    // how many AP recover used today?
	ActivityPointPaymentRecoveryDailyResetAt  int    `json:"activity_point_payment_recovery_daily_reset_at"` // in unix time 1684854000?
	GameMoney                                 int64  `json:"game_money"`                                     // Gold
	CardExp                                   int64  `json:"card_exp"`                                       // exp currency
	FreeSnsCoin                               int    `json:"free_sns_coin"`                                  // free gem
	AppleSnsCoin                              int    `json:"apple_sns_coin"`                                 // paid gem ios
	GoogleSnsCoin                             int    `json:"google_sns_coin"`                                // paid gem android
	SubscriptionCoin                          int    `json:"subscription_coin"`                              // member coin
	BirthDate                                 int    `json:"birth_date"`                                     // yyyymm ?
	BirthMonth                                int    `json:"birth_month"`
	BirthDay                                  int    `json:"birth_day"`
	LatestLiveDeckId                          int    `json:"latest_live_deck_id"` // last used live formation
	MainLessonDeckId                          int    `json:"main_lesson_deck_id"` // last used training formation
	FavoriteMemberId                          int    `json:"favorite_member_id"`  // partner id
	LastLiveDifficultyId                      int    `json:"last_live_difficulty_id"`
	LpMagnification                           int    `json:"lp_magnification"`                    // unused feature, always 1
	EmblemId                                  int    `json:"emblem_id"`                           // title
	DeviceToken                               string `json:"device_token"`                        //  some sort of salted encryption?, used to prevent using multiple device
	TutorialPhase                             int    `json:"tutorial_phase"`                      // 99 = done
	TutorialEndAt                             int    `json:"tutorial_end_at"`                     // in unix second
	LoginDays                                 int    `json:"login_days"`                          // amount of days logged in
	NaviTapCount                              int    `json:"navi_tap_count"`                      // number of partner tap that will increase bond
	NaviTapRecoverAt                          int    `json:"navi_tap_recover_at"`                 // in unix time
	IsAutoMode                                bool   `json:"is_auto_mode"`                        // is autoplay enabled, for restarting app without finishing lives
	MaxScoreLiveDifficultyMasterId            int    `json:"max_score_live_difficulty_master_id"` // not the one featured in profile?
	LiveMaxScore                              int    `json:"live_max_score"`
	MaxComboLiveDifficultyMasterId            int    `json:"max_combo_live_difficulty_master_id"`
	LiveMaxCombo                              int    `json:"live_max_combo"`
	LessonResumeStatus                        int    `json:"lesson_resume_status"`          // for quitting while training, the number probably mean the phase of the training
	AccessoryBoxAdditional                    int    `json:"accessory_box_additional"`      // additional accessory slot, max is 400 in official
	TermsOfUseVersion                         int    `json:"terms_of_use_version"`          // 3 mean nothing to accept
	BootstrapSifidCheckAt                     int    `json:"bootstrap_sifid_check_at"`      // not really what it sound like, probably safe to ignore it
	GdprVersion                               int    `json:"gdpr_version"`                  // set to 0 for outside the EU
	MemberGuildMemberMasterId                 int    `json:"member_guild_member_master_id"` // member id of the channel joined
	MemberGuildLastUpdatedAt                  int    `json:"member_guild_last_updated_at"`  // unix time stamp, last time joining the channel (used to allow only 1 channel changing)
	Cash                                      int    `json:"cash"`
}
