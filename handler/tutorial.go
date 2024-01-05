package handler

import (
	"elichika/userdata"
	"elichika/config"
	"elichika/enum"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CorePlayableEnd(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()

	if session.UserStatus.TutorialPhase == enum.TutorialPhaseStory2 {
		session.UserStatus.TutorialPhase = enum.TutorialPhaseStory3
	} else if session.UserStatus.TutorialPhase == enum.TutorialPhaseStory4 {
		session.UserStatus.TutorialPhase = enum.TutorialPhaseFavoriateMember
	}

	signBody := session.Finalize("{}", "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

//uhh - what does this do?
func TimingAdjusterEnd(ctx *gin.Context) {
	CorePlayableEnd(ctx)
}

func PhaseEnd(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()

	if session.UserStatus.TutorialPhase == enum.TutorialPhaseFinal {
		// I think it ends here? Not 100% sure
		session.UserStatus.TutorialPhase = enum.TutorialFinished
		session.UserStatus.TutorialEndAt = int(time.Now().Unix())
	}

	signBody := session.Finalize("{}", "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func TutorialSkip(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	session := userdata.GetSession(ctx, userID)
	defer session.Close()

	if session.UserStatus.TutorialPhase != enum.TutorialFinished {
		session.UserStatus.TutorialPhase = enum.TutorialFinished
		session.UserStatus.TutorialEndAt = int(time.Now().Unix())
	}

	signBody := session.Finalize("{}", "user_model")
	resp := SignResp(ctx, signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
