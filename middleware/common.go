package middleware

import (
	"elichika/config"

	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Common(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}
	defer ctx.Request.Body.Close()
	ctx.Set("reqBody", string(body))

	lang, _ := ctx.GetQuery("l")
	if lang == "" {
		lang = "ja"
	}
	ctx.Set("locale", config.Locales[lang])
	ctx.Set("masterdata.db", config.Locales[lang].MasterdataEngine)
	ctx.Set("gamedata", config.Locales[lang].Gamedata)
	ctx.Set("dictionary", config.Locales[lang].Dictionary)

	// if lang == "" {
	// 	ctx.Set("gamedata", config.GamedataJp)
	// 	handler.IsGlobal = false
	// 	handler.MasterVersion = config.MasterVersionJp
	// 	handler.StartUpKey = "5f7IZY1QrAX0D49g"
	// } else {
	// 	ctx.Set("masterdata.db", config.MasterdataEngGl)
	// 	ctx.Set("gamedata", config.GamedataGl)
	// 	handler.IsGlobal = true
	// 	handler.MasterVersion = config.MasterVersionGl
	// 	handler.StartUpKey = "TxQFwgNcKDlesb93"
	// }

	userID, _ := strconv.Atoi(ctx.Query("u"))
	ctx.Set("user_id", userID)

	ctx.Set("ep", ctx.Request.URL.String())

	ctx.Next()
}
