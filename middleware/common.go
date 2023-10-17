package middleware

import (
	"elichika/config"
	"elichika/handler"

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
		ctx.Set("masterdata.db", config.MasterdataEngJp)
		ctx.Set("gamedata", config.GamedataJp)
		handler.IsGlobal = false
		handler.MasterVersion = config.MasterVersionJp
		handler.StartUpKey = "5f7IZY1QrAX0D49g"
	} else {
		ctx.Set("masterdata.db", config.MasterdataEngGl)
		ctx.Set("gamedata", config.GamedataGl)
		handler.IsGlobal = true
		handler.MasterVersion = config.MasterVersionGl
		handler.StartUpKey = "TxQFwgNcKDlesb93"
	}

	UserID, _ := strconv.Atoi(ctx.Query("u"))
	ctx.Set("user_id", UserID)

	ctx.Set("ep", ctx.Request.URL.String())

	ctx.Next()
}
