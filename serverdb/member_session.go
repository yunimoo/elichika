package serverdb

import (
	"elichika/model"

	"fmt"
)

func (session* Session) GetMember(memberMasterID int) model.UserMemberInfo {
	member := model.UserMemberInfo{}
	exists, err := Engine.Table("s_user_member").
	Where("user_id = ? AND member_master_id = ?", session.UserInfo.UserID, memberMasterID).Get(&member)
	// inserted at login if not exist
	if (err != nil) {
		panic(err)
	}
	if (!exists) {
		panic("member not found")
	}
	return member
}

func (session* Session) GetAllMembers() []model.UserMemberInfo {
	members := []model.UserMemberInfo{}
	err := Engine.Table("s_user_member").Where("user_id = ?", session.UserInfo.UserID).Find(&members)
	if (err != nil) {
		panic(err)
	}
	return members
}

func (session* Session) UpdateMember(member model.UserMemberInfo) {
	session.MemberDiffs[member.MemberMasterID] = member
}

func (session* Session) InsertMembers(members []model.UserMemberInfo) {
	affected, err := Engine.Table("s_user_member").Insert(&members)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", affected, " members")
}

func (session* Session) FinalizeMemberDiffs() []any {
	userMemberByMemberID := []any{}
	for memberMasterID, member := range session.MemberDiffs {
		userMemberByMemberID = append(userMemberByMemberID, memberMasterID)
		userMemberByMemberID = append(userMemberByMemberID, member)
		affected, err := Engine.Table("s_user_member").
			Where("user_id = ? AND member_master_id = ?", session.UserInfo.UserID, memberMasterID).AllCols().Update(member)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userMemberByMemberID
}

func (session* Session) GetLovePanelCellIDs(memberID int) []int {
	cells := []int{}
	err := Engine.Table("s_user_member_love_panel").
	Where("user_id = ? AND member_id = ?", session.UserInfo.UserID, memberID).Cols("member_love_panel_cell_id").
	Find(&cells)
	if err != nil {
		panic(err)
	}
	return cells
}