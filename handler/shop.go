package handler

import (
	"elichika/config"

	"net/http"
	// "encoding/json"

	"github.com/gin-gonic/gin"
)

// Not implementing this system seriously

func FetchShopTop(ctx *gin.Context) {
	resp := SignResp(ctx, GetData("fetchShopTop.json"), config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchShopPack(ctx *gin.Context) {
	resp := SignResp(ctx, GetData("fetchShopPack.json"), config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchShopSnsCoin(ctx *gin.Context) {
	resp := SignResp(ctx, GetData("fetchShopSnsCoin.json"), config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
