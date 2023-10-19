// this package represents the gamedata
// i.e. settings, parameters of the games that are shared between all users
// the data are stored in both Gamedata.db and serverdata.db
// some data are loaded only from one of the 2, but some data need boths
//
// eventually no handling function should interact with the db and only interact with this package
// this is done both to reduce the time necessary to look into db, as well as to simplify accessing data to the most
// relevant id, instead of having to access multiple tables or use magic id system
// for example, everything related to a single card / accessory will use that card / accessory master id as id
// everything related to all card / accessory of a rarity will use that rarity as id
// relation is defined by Gamedata.db or serverdata.db
// i.e. some setting might be the same across all accessory of a rarity, but as long as it's store separately in the db,
// it's stored separately here
package gamedata

import (
	"elichika/dictionary"

	"xorm.io/xorm"
)

type Gamedata struct {
	Accessory Accessory
	Trade     Trade
	LiveParty LiveParty
}

func (Gamedata *Gamedata) Init(masterdata *xorm.Engine, serverdata *xorm.Engine, dictionary *dictionary.Dictionary) {
	masterdata_session := masterdata.NewSession()
	serverdata_session := serverdata.NewSession()
	defer masterdata_session.Close()
	defer serverdata_session.Close()
	Gamedata.Accessory.Load(masterdata_session, serverdata_session, dictionary)
	Gamedata.Trade.Load(masterdata_session, serverdata_session, dictionary)
	Gamedata.LiveParty.Load(masterdata_session, serverdata_session, dictionary)
}
