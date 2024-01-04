package handler

import (
	"elichika/userdata"
	"elichika/config"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CorePlayableEnd(ctx *gin.Context) {
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

func PhaseEnd(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()

	// This should go to the next tutorial phase
	if session.UserStatus.TutorialPhase != 99 {
		session.UserStatus.TutorialPhase = 99
		session.UserStatus.TutorialEndAt = int(time.Now().Unix())
	}

	signBody := session.Finalize("{}", "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func TutorialSkip(ctx *gin.Context) {
	PhaseEnd(ctx)
}
