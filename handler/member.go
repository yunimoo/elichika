package handler

import (
	"elichika/config"
	"elichika/klab"
	"elichika/model"
	"elichika/serverdb"

	"encoding/json"
	"net/http"
	// "time"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func OpenMemberLovePanel(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type OpenMemberLovePanelReq struct {
		MemberID               int   `json:"member_id"`
		MemberLovePanelID      int   `json:"member_love_panel_id"`
		MemberLovePanelCellIDs []int `json:"member_love_panel_cell_ids"`
	}
	req := OpenMemberLovePanelReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	CheckErr(err)
	session := serverdb.GetSession(UserID)
	panel := session.GetMemberLovePanel(req.MemberID)
	panel.MemberLovePanelCellIDs = append(panel.MemberLovePanelCellIDs, req.MemberLovePanelCellIDs...)
	session.UpdateMemberLovePanel(panel)

	// if is full panel, then we have to send a basic info trigger to actually open up the next panel
	// we actually have to check for bond level, otherwise it freeze
	// this mean we also have to implement opening board when bond level catch up too
	n := len(panel.MemberLovePanelCellIDs)
	if n >= 5 {
		tier := panel.MemberLovePanelCellIDs[n-1]/10000 + 1
		if tier == panel.MemberLovePanelCellIDs[n-5]/10000+1 {
			member := session.GetMember(panel.MemberID)
			maxTier := klab.MaxLovePanelTierFromLoveLevel(member.LoveLevel)
			if tier < maxTier {
				// unlock the next board
				session.AddTriggerBasic(&model.TriggerBasic{
					TriggerID:       0, // filled by session
					InfoTriggerType: 28,
					LimitAt:         nil,
					Description:     nil,
					ParamInt:        (tier+1)*1000 + panel.MemberID})
			}
		}
	}

	signBody := session.Finalize(GetData("userModel.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
