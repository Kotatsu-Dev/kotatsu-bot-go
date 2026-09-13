package db

import (
	//Внутренние пакеты проекта
	"rr/kotatsutgbot/rr_debug"
)

type UserActivity struct {
	UserID     uint `gorm:"primaryKey"`
	ActivityID uint `gorm:"primaryKey"`
	Visited    bool
}

type UserActivity_CreateJSON struct {
	UserID     uint `json:"user_id"`
	ActivityID uint `json:"activity_id"`
	Visited    bool `json:"visited"`
}

type UserActivity_ReadJSON struct {
	UserID     uint `json:"user_id" gorm:"primaryKey"`
	ActivityID uint `json:"activity_id" gorm:"primaryKey"`
	Visited    bool `json:"visited"`
}

func (ua *UserActivity) ToRead() *UserActivity_ReadJSON {
	return &UserActivity_ReadJSON{
		UserID:     ua.UserID,
		ActivityID: ua.ActivityID,
		Visited:    ua.Visited,
	}
}

func UserActivityToReadSlice(uas []UserActivity) []UserActivity_ReadJSON {
	res := make([]UserActivity_ReadJSON, len(uas))
	for i, ua := range uas {
		res[i] = *ua.ToRead()
	}
	return res
}

func DB_CREATE_UserActivity(user_activity_to_add *UserActivity_CreateJSON) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user User
	db.First(&user, user_activity_to_add.UserID)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	var activity Activity
	db.First(&activity, user_activity_to_add.ActivityID)
	if activity.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	var user_activity UserActivity
	result := db.Where("user_id = ? AND activity_id = ?", user_activity_to_add.UserID, user_activity_to_add.ActivityID).First(&user_activity)
	if result.RowsAffected != 0 {
		return DB_ANSWER_OBJECT_EXISTS
	}

	user_activity = UserActivity{
		UserID:     user_activity_to_add.UserID,
		ActivityID: user_activity_to_add.ActivityID,
		Visited:    user_activity_to_add.Visited,
	}

	db.Create(&user_activity)
	return DB_ANSWER_SUCCESS
}

func DB_GET_UserActivity(user_id uint, activity_id uint) (int, *UserActivity_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	user_activity := new(UserActivity)
	result := db.Where("user_id = ? AND activity_id = ?", user_id, activity_id).First(&user_activity)
	if result.RowsAffected == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	return DB_ANSWER_SUCCESS, user_activity.ToRead()
}

func DB_GET_UserActivities_BY_UserID(user_id uint) []UserActivity_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user_activities []UserActivity
	db.Where("user_id = ?", user_id).Find(&user_activities)

	return UserActivityToReadSlice(user_activities)
}

func DB_GET_UserActivities_BY_ActivityID(activity_id uint) []UserActivity_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user_activities []UserActivity
	db.Where("activity_id = ?", activity_id).Find(&user_activities)

	return UserActivityToReadSlice(user_activities)
}

func DB_GET_UserActivities() []UserActivity_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user_activities []UserActivity
	db.Find(&user_activities)

	return UserActivityToReadSlice(user_activities)
}

func DB_UPDATE_UserActivity(update_json map[string]interface{}) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	user_id_, ok := update_json["user_id"].(float64)
	if !ok {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}
	user_id := uint(user_id_)

	activity_id_, ok := update_json["activity_id"].(float64)
	if !ok {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}
	activity_id := uint(activity_id_)

	var user_activity UserActivity
	result := db.Where("user_id = ? AND activity_id = ?", user_id, activity_id).First(&user_activity)
	if result.RowsAffected == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	for key, value := range update_json {
		switch key {
		case "visited":
			if v, ok := value.(bool); ok && v != user_activity.Visited {
				user_activity.Visited = v
			}
		}
	}

	db.Save(&user_activity)
	return DB_ANSWER_SUCCESS
}

func DB_UPDATE_UserActivity_Visited(user_id uint, activity_id uint, visited bool) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user_activity UserActivity
	result := db.Where("user_id = ? AND activity_id = ?", user_id, activity_id).First(&user_activity)
	if result.RowsAffected == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	user_activity.Visited = visited

	db.Save(&user_activity)
	return DB_ANSWER_SUCCESS
}

func DB_DELETE_UserActivity(user_id uint, activity_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user_activity UserActivity
	result := db.Where("user_id = ? AND activity_id = ?", user_id, activity_id).First(&user_activity)
	if result.RowsAffected == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	result = db.Unscoped().Delete(&user_activity)
	if result.Error != nil {
		rr_debug.PrintLOG("UserActivities.go", "DB_DELETE_UserActivity", "Error deleting user_activity:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}

func DB_DELETE_UserActivities_BY_UserID(user_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Unscoped().Where("user_id = ?", user_id).Delete(&UserActivity{})
	if result.Error != nil {
		rr_debug.PrintLOG("UserActivities.go", "DB_DELETE_UserActivities_BY_UserID", "Error deleting user_activities:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}

func DB_DELETE_UserActivities_BY_ActivityID(activity_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Unscoped().Where("activity_id = ?", activity_id).Delete(&UserActivity{})
	if result.Error != nil {
		rr_debug.PrintLOG("UserActivities.go", "DB_DELETE_UserActivities_BY_ActivityID", "Error deleting user_activities:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}

func DB_DELETE_UserActivities() int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Exec("DELETE FROM user_activities")
	if result.Error != nil {
		rr_debug.PrintLOG("UserActivities.go", "DB_DELETE_UserActivities", "Error deleting user_activities:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}
