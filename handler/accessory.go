package handler

import (
	"elichika/config"
	// "elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func UpdateIsLock(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type UpdateIsLockReq struct {
		UserAccessoryID int64 `xorm:"'user_accessory_id' pk" json:"user_accessory_id"`
		IsLock          bool  `xorm:"'is_lock'" json:"is_lock"`
	}
	req := UpdateIsLockReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)
	userID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, userID)
	accessory := session.GetUserAccessory(req.UserAccessoryID)
	accessory.IsLock = req.IsLock
	session.UpdateUserAccessory(accessory)

	signBody := session.Finalize(GetData("userModel.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func AccessoryMelt(ctx *gin.Context) {
	// disassemble
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type MeltReq struct {
		UserAccessoryIDs []int64 `json:"user_accessory_ids"`
	}
	req := MeltReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)
	userID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, userID)
	for _, userAccessoryID := range req.UserAccessoryIDs {
		accessory := session.GetUserAccessory(userAccessoryID)
		// TODO: award items
		accessory.AccessoryMasterID = 0 // marked for delete
		session.UpdateUserAccessory(accessory)
	}

	signBody := session.Finalize(GetData("userModel.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
