package item

import (
	"elichika/enum"
	"elichika/model"
)

var (
	Gold = model.Content{
		ContentType:   enum.ContentTypeGameMoney,
		ContentID:     0,
		ContentAmount: 1,
	}
)
