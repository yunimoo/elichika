// this package represents masterdata.db in its entirity
// eventually no handling function should interact with masterdata.db, and only interact with this package
// this is done both to reduce the time necessary to look into db, as well as to simplify accessing data to the most
// relevant id, instead of having to access multiple tables or use magic id system
// for example, everything related to a single card / accessory will use that card / accessory master id as id
// everything related to all card / accessory of a rarity will use that rarity as id
// relation is defined by the stock masterdata.db
// i.e. some setting might be the same across all accessory of a rarity, but as long as it's store separately in the db,
// it's stored separately here
package masterdata

import (
	"xorm.io/xorm"
)

type Masterdata struct {
	Accessory Accessory
}

func (masterdata *Masterdata) Init(engine *xorm.Engine) {
	session := engine.NewSession()
	defer session.Close()
	masterdata.Accessory.Load(session)
}
