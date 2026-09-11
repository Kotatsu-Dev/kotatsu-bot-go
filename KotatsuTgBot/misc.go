package main

import (
	//Внутренние пакеты проекта
	"rr/kotatsutgbot/db"

	//Системные пакеты
	"fmt"
)

// ===================================================================
//
//	CREATE
//
// ===================================================================
// Зарегистрировать пользователя
func regUser(user_tg_id int64, full_tg_name string, user_name string) int {
	user_to_add := db.User_CreateJSON{
		UserTgID:   user_tg_id,
		UserName:   user_name,
		FullTgName: full_tg_name,
	}

	db_answer_code, _ := db.DB_CREATE_User(&user_to_add)
	return db_answer_code
}

// ===================================================================
//
//	MISC
//
// ===================================================================

// Получить ссылку на профиль
func GetProfileTgURL(username string) string {
	if username != "" {
		profileURL := fmt.Sprintf("https://t.me/%s", username)
		return profileURL
	} else {
		return ""
	}
}
