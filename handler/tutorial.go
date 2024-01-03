package handler

import (
	"elichika/userdata"
	"elichika/config"

	"net/http"

	"github.com/gin-gonic/gin"
)

func CorePlayableEnd(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()

	if session.UserStatus.TutorialPhase != 99 {
		session.UserStatus.TutorialPhase += 1
	}

	signBody := session.Finalize("{}", "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
