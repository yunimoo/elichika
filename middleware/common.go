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

	userID, _ := strconv.Atoi(ctx.Query("u"))
	ctx.Set("user_id", userID)

	ctx.Set("ep", ctx.Request.URL.String())

	ctx.Next()
}
