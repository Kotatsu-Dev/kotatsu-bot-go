package main

import (
	//Внутренние пакеты проекта
	"rr/kotatsutgbot/db"

	//Системные пакеты
	"fmt"
	"math/rand"
	"strings"
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

	db_answer_code := db.DB_CREATE_User(&user_to_add)
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

func splitName(input string) (user1, user2, user3 string) {
	words := strings.Fields(input) // Разбиваем строку на слова

	if len(words) >= 3 {
		user1 = words[0]
		user2 = words[1]
		user3 = words[2]
	} else if len(words) == 2 {
		user1 = words[0]
		user2 = words[1]
		user3 = ""
	} else if len(words) == 1 {
		user1 = words[0]
		user2 = ""
		user3 = ""
	}

	return user1, user2, user3
}

// func UpdateBotApply(bot *tgbotapi.BotAPI) {
// 	var list_users []db.User
// 	msg := tgbotapi.NewMessage(0, "")
// 	var message_text, message_text_club string

// 	user := new(db.User)
// 	list_users, error_code, _ := user.GetManyRecordsFromDB()

// 	message_text = "🔄 Вышло обновление бота! 🔄" + "\n" + "\n" +
// 		"Оптимизация основных функций бота" + "\n" +
// 		"Улучшена стабильность работы"

// 	message_text_club = "🔄 Вышло обновление бота! 🔄" + "\n" + "\n" +
// 		"Оптимизация основных функций бота" + "\n" +
// 		"Улучшена стабильность работы" + "\n" +
// 		"В главном меню появился новый раздел '🤝 Акции и партнёры'"

// 	switch error_code {
// 	case db.ERROR_OK:
// 		for _, user := range list_users {
// 			msg.ChatID = user.UserTgID

// 			if user.IsClubMember {
// 				msg.Text = message_text_club
// 				msg.ReplyMarkup = keyboards.MainMenuButtonsReg
// 			} else {
// 				msg.Text = message_text
// 			}

// 			_, err := bot.Send(msg)
// 			if err != nil {
// 				rr_debug.PrintLOG("misc.go", "UpdateBotApply", "bot.Send", "Ошибка отправки сообщения", err.Error())
// 			}
// 		}
// 	}
// 	fmt.Println("UPDATE SUCCESS!!!")
// }

// Генератор случайных чисел
func generateRandomNumber(digits int) int {
	min := int64(1)
	max := int64(1)

	// Вычисляем минимальное и максимальное значение для указанного количества цифр
	for i := 0; i < digits-1; i++ {
		min *= 10
		max *= 10
	}

	max = max*10 - 1

	// Генерируем случайное число в заданном диапазоне
	randomValue := rand.Int63n(max-min+1) + min

	return int(randomValue)
}
