package model

type Content struct {
	ContentType   int   `xorm:"<-" json:"content_type"`
	ContentID     int   `xorm:"<- 'content_id'" json:"content_id"`
	ContentAmount int64 `xorm:"<-" json:"content_amount"`
}

type RewardDrop struct { // unused
	DropColor int     `json:"drop_color"`
	Content   Content `json:"content"`
	IsRare    bool    `json:"is_rare"`
	BonusType *int    `json:"bonus_type"`
}
