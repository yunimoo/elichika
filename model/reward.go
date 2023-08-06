package model

type RewardByContent struct {
	ContentType   int `xorm:"<-"`
	ContentID     int `xorm:"<- 'content_id'"`
	ContentAmount int `xorm:"<-"`
}
