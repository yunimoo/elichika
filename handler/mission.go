package handler

import (
	"elichika/config"
	"elichika/userdata"

	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchMission(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, UserID)
	defer session.Close()
	signBody := session.Finalize(GetData("fetchMission.json"), "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ClearMissionBadge(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, UserID)
	defer session.Close()
	signBody := session.Finalize(GetData("clearMissionBadge.json"), "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
