package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"encoding/json"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func ReadCardGradeUp(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]

	type ReadCardGradeUpReq struct {
		TriggerID int64 `json:"trigger_id"`
	}

	session := serverdb.GetSession(UserID)

	triggerReq := ReadCardGradeUpReq{}

	if err := json.Unmarshal([]byte(reqBody.String()), &triggerReq); err != nil {
		panic(err)
	}

	session.AddCardGradeUpTrigger(triggerReq.TriggerID, nil)
	resp := session.Finalize(GetUserData("userModelDiff.json"), "user_model_diff")
	resp, _ = sjson.Set(resp, "user_model_diff.user_status", GetUserStatus())
	resp = SignResp(ctx.GetString("ep"), resp, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	fmt.Println(resp)
}
