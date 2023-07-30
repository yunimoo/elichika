package handler

// import (
// 	"elichika/config"
// 	"elichika/serverdb"

// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/tidwall/gjson"
// 	// "github.com/tidwall/sjson"
// 	// "xorm.io/xorm"
// )

// func ReadCardGradeUp(ctx *gin.Context) {
// 	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]

// 	type ReadCardGradeUpReq struct {
// 		TriggerID int64 `json:"trigger_id"`
// 	}

// 	uid, _ := strconv.Atoi(ctx.Query("u"))
// 	session := serverdb.GetSession(uid)

// 	triggerReq := ReadCardGradeUpReq{}

// 	if err := json.Unmarshal([]byte(reqBody.String()), &triggerReq); err != nil {
// 		panic(err)
// 	}

// 	session.AddCardGradeUpTrigger(triggerReq.TriggerID, nil)
// 	resp := SignResp(ctx.GetString("ep"), session.Finalize("user_model_diff"), config.SessionKey)
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.String(http.StatusOK, resp)
// 	fmt.Println(resp)
// }
