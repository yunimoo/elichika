package handler

import (
	"elichika/config"

	"net/http"

	"github.com/gin-gonic/gin"
)

func LoveRankingFetch(ctx *gin.Context) {
	// TODO: fetch from db instead
	// probably needs to store the ranking somewhere, or just have very big sql
	resp := SignResp(ctx, GetData("loveRankingFetch.json"), config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
