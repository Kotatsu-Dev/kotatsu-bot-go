package main

import (
	//Системные пакеты
	"fmt"
)

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
