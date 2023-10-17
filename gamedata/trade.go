// trade or exchange
// trades are stored inside both serverdata.db and masterdata.db
// for the clients to read them, mismatch between serverdata.db and masterdata.db will be resolved toward masterdata.db side
// i.e. the masterdata.db ID will be used and such
// Requirements:
// - TradeID in serverdata.db/s_trade must be the same TradeID in masterdata.db/m_trade
// - TradeID reference in serverdata.db/s_trade_product must refer to the correct TradeID above
// - The number of products per TradeID in serverdata.db/s_trade_product should be the same as the
// amount of products per TradeID in m_trade_product. If they are not the same, some items will be ignored
// or filled with undetermined items

package gamedata

import (
	"elichika/model"
	"elichika/utils"

	"fmt"
	"os"

	"xorm.io/xorm"
)

type Trade struct {
	TradesByType [3][]model.Trade    // map from trade type to array of Trade
	Trades       map[int]model.Trade // map from TradeID to Trade
	Products     map[int]model.TradeProduct
}

func (trade *Trade) Load(masterdata_db, serverdata_db *xorm.Session) {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		return
	}
	trade.Trades = make(map[int]model.Trade)
	trade.Products = make(map[int]model.TradeProduct)
	serverTrades := make(map[int]model.Trade)
	err := serverdata_db.Table("s_trade").Find(&serverTrades)
	utils.CheckErr(err)
	for tradeID, t := range serverTrades {
		exists, err := masterdata_db.Table("m_trade").Where("id = ?", tradeID).
			Cols("trade_type", "source_content_type", "source_content_id").Get(
			&t.TradeType, &t.SourceContentType, &t.SourceContentID)
		utils.CheckErr(err)
		if !exists {
			fmt.Println("Warning: Skipped trade ", tradeID, " (did not exists in masterdata.db)")
			continue
		}
		serverProducts := []model.TradeProduct{}
		err = serverdata_db.Table("s_trade_product").Where("trade_id = ?", tradeID).
			OrderBy("product_id").Find(&serverProducts)
		utils.CheckErr(err)
		clientProductIDs := []int{}
		err = masterdata_db.Table("m_trade_product").Where("trade_master_id = ?", tradeID).
			OrderBy("id").Cols("id").Find(&clientProductIDs)
		utils.CheckErr(err)
		n := len(serverProducts)
		m := len(clientProductIDs)
		for ; n < m; n++ {
			serverProducts = append(
				serverProducts,
				model.TradeProduct{
					TradeID:      tradeID,
					SourceAmount: 1,
					StockAmount: nil,
					ActualContent: model.Content{
						ContentType:   10,
						ContentID:     0,
						ContentAmount: 1,
					},
				})
		}
		serverProducts = serverProducts[0:m]
		for i := 0; i < m; i++ {
			serverProducts[i].ProductID = clientProductIDs[i]
			serverProducts[i].DummyID = clientProductIDs[i]
			serverProducts[i].Contents = append(serverProducts[i].Contents, serverProducts[i].ActualContent)
			trade.Products[clientProductIDs[i]] = serverProducts[i]
		}
		t.Products = serverProducts
		trade.TradesByType[t.TradeType] = append(trade.TradesByType[t.TradeType], t)
		trade.Trades[t.TradeID] = t
	}
}
