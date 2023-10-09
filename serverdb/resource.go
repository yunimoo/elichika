package serverdb

// this file contain all the logic for actually changing what an user own
// more precisely, it handle things found in m_content_route_guide
// things like player bonds / exp are handled elsewhere

import (
	"elichika/model"
	"elichika/utils"

	"fmt"

	"github.com/tidwall/sjson"
	"xorm.io/xorm"
)

type UserResource struct {
	UserID        int   `xorm:"pk 'user_id'"`
	ContentType   int   `xorm:"pk 'content_type'"`
	ContentID     int   `xorm:"pk 'content_id'"`
	ContentAmount int64 `xorm:"'content_amount'"`
}

var (
	ResourceHandler   map[int]func(*Session, int, int, int64)
	ResourceFinalizer map[int]func(*Session, *xorm.Session, int, map[int]UserResource) (string, string)
	ResourceGroupKey  map[int]string
	ResourceItemKey   map[int]string
)

func (session *Session) AddGameMoney(money int) {
	fmt.Println("TODO: Add money:", money)
}

func (session *Session) RemoveGameMoney(money int) {
	fmt.Println("TODO: Remove money:", money)
}

func (session *Session) GetUserResource(contentType, contentID int) UserResource {
	_, exists := session.UserResourceDiffs[contentType]
	if !exists {
		session.UserResourceDiffs[contentType] = make(map[int]UserResource)
	}
	resource, exists := session.UserResourceDiffs[contentType][contentID]
	if exists {
		return resource
	}
	// load from db
	exists, err := Engine.Table("s_user_resource").Where("user_id = ? AND content_type = ? AND content_id = ?",
		session.UserStatus.UserID, contentType, contentID).Get(&resource)
	utils.CheckErr(err)
	if !exists {
		resource = UserResource{
			UserID:        session.UserStatus.UserID,
			ContentType:   contentType,
			ContentID:     contentID,
			ContentAmount: 100000,
		}
	}
	return resource
}

func (session *Session) UpdateUserResource(resource UserResource) {
	_, exists := session.UserResourceDiffs[resource.ContentType]
	if !exists {
		session.UserResourceDiffs[resource.ContentType] = make(map[int]UserResource)
	}
	session.UserResourceDiffs[resource.ContentType][resource.ContentID] = resource
}

func (session *Session) AddResource(resource model.Content) {
	// fmt.Println(resource)
	handler, exists := ResourceHandler[resource.ContentType]
	if !exists {
		fmt.Println("TODO: Add handler for content type ", resource.ContentType)
		return
	}
	handler(session, resource.ContentType, resource.ContentID, resource.ContentAmount)
}

func (session *Session) RemoveResource(resource model.Content) {
	// fmt.Println(resource)
	handler, exists := ResourceHandler[resource.ContentType]
	if !exists {
		fmt.Println("TODO: Add handler for content type ", resource.ContentType)
		return
	}
	handler(session, resource.ContentType, resource.ContentID, -resource.ContentAmount)
}

// return an array of pair of key / content to set that key to
// this key start from user_model / user_model_diff (so add that to the front)
// the value inside []any are PREMARSHALED json, so use sjson.SetRaw
func (session *Session) FinalizeUserResources(dbSession *xorm.Session) ([]string, []string) {
	keys := []string{}
	values := []string{}
	for contentType, resourceDiffByContentID := range session.UserResourceDiffs {
		key, value := ResourceFinalizer[contentType](session, dbSession, contentType, resourceDiffByContentID)
		keys = append(keys, key)
		values = append(values, value)
	}
	return keys, values
}

func init() {
	ResourceHandler = make(map[int]func(*Session, int, int, int64))
	ResourceHandler[1] = SnsCoinResourceHandler
	ResourceHandler[7] = SuitResourceHandler
	ResourceHandler[12] = GenericResourceHandler // training items (macarons, memorials)
	ResourceHandler[24] = GenericResourceHandler // accessory stickers
	ResourceHandler[25] = GenericResourceHandler // accessory rarity up items

	ResourceFinalizer = make(map[int]func(*Session, *xorm.Session, int, map[int]UserResource) (string, string))
	ResourceGroupKey = make(map[int]string)
	ResourceItemKey = make(map[int]string)
	ResourceFinalizer[12] = GenericResourceFinalizer
	ResourceGroupKey[12] = "user_training_material_by_item_id"
	ResourceItemKey[12] = "training_material_master_id"
	ResourceFinalizer[24] = GenericResourceFinalizer
	ResourceGroupKey[24] = "user_accessory_level_up_item_by_id"
	ResourceItemKey[24] = "accessory_level_up_item_master_id"
	ResourceFinalizer[25] = GenericResourceFinalizer
	ResourceGroupKey[25] = "user_accessory_rarity_up_item_by_id"
	ResourceItemKey[25] = "accessory_rarity_up_item_master_id"
}

func GenericResourceHandler(session *Session, contentType, contentID int, contentAmount int64) {
	resource := session.GetUserResource(contentType, contentID)
	resource.ContentAmount += contentAmount
	session.UpdateUserResource(resource)
}

func GenericResourceFinalizer(session *Session, dbSession *xorm.Session, contentType int,
	resourceDiffByContentID map[int]UserResource) (string, string) {
	groupKey := ResourceGroupKey[contentType]
	itemKey := ResourceItemKey[contentType]
	result := "[]"
	index := 0
	for contentID, resource := range resourceDiffByContentID {
		// update or insert the resource
		affected, err := dbSession.Table("s_user_resource").Where("user_id = ? AND content_type = ? AND content_id = ?",
			session.UserStatus.UserID, contentType, contentID).Update(resource)
		utils.CheckErr(err)
		if affected == 0 { // doesn't exists, insert
			dbSession.Table("s_user_resource").Insert(resource)
		}

		// set the json
		result, _ = sjson.Set(result, fmt.Sprintf("%d", index), contentID)
		index++
		result, _ = sjson.Set(result, fmt.Sprintf("%d.%s", index, itemKey), contentID)
		result, _ = sjson.Set(result, fmt.Sprintf("%d.amount", index), resource.ContentAmount)
		index++
	}
	return groupKey, result
}

func SuitResourceHandler(session *Session, suitContentType, suitMasterID int, amount int64) {
	session.InsertUserSuit(model.UserSuit{
		UserID:       session.UserStatus.UserID,
		SuitMasterID: suitMasterID,
		IsNew:        true})
}

func SnsCoinResourceHandler(session *Session, snsCoinContentType, snsCoinType int, amount int64) {
	fmt.Println("TODO: Add economy")
}

func TrainingMaterialResourceHandler(session *Session, trainingMaterialContentType, trainingMaterialType int, amount int64) {
	fmt.Println("TODO: Add training items")
}
