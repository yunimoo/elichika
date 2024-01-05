package request

type ReadLoginBonusRequest struct {
	LoginBonusType int `json:"login_bonus_type"`
	LoginBonusId   int `json:"login_bonus_id"`
}
