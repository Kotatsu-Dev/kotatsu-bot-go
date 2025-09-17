package db

import (
	//Внутренние пакеты проекта

	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	//Системные пакеты
	"time"
)

type Gender string

type User struct {
	gorm.Model
	Step                  int         `json:"step"`                                                                        // Текущий шаг
	UserTgID              int64       `json:"user_tg_id"`                                                                  // ID пользователя в Телеграм
	LastMessageID         int         `json:"last_message_id"`                                                             // ID последнего сообщения от бота
	UserName              string      `json:"user_name"`                                                                   // Имя пользователя в Телеграм
	FullTgName            string      `json:"full_tg_name"`                                                                // Полное имя пользователя в Телеграм
	Gender                Gender      `json:"gender"`                                                                      // Пол пользователя
	IsVisitedEvents       bool        `json:"is_visited_events"`                                                           // Посетил ли пользователь достаточное количество мероприятий
	ISU                   string      `json:"isu"`                                                                         // ИСУ для ИТМО
	FullName              string      `json:"full_name"`                                                                   // Имя пользователя
	PhoneNumber           string      `json:"phone_number"`                                                                // Номер телефона пользователя
	SecretCode            string      `json:"secret_code"`                                                                 // Секретный код пользователя
	IsITMO                bool        `json:"is_itmo"`                                                                     // Студент ИТМО
	IsClubMember          bool        `json:"is_club_member"`                                                              // Член клуба
	IsSubscribeNewsletter bool        `json:"is_subscribe_newsletter"`                                                     // Подписка на рассылку
	IsSentRequest         bool        `json:"is_sent_request"`                                                             // Отправлена ли заявка
	IsFilledData          bool        `json:"is_filled_data"`                                                              // Заполнены ли данные?
	TempActivityID        int         `json:"temp_activity_id"`                                                            // Временное хранение при записи на мероприятие
	MyActivities          []*Activity `json:"my_activities" gorm:"many2many:user_activities;constraint:OnDelete:CASCADE;"` // Простой список моих мероприятий
	LinkMyAnimeList       string      `json:"link_my_anime_list"`                                                          // Мой список аниме
	MyRequest             *Request    `json:"my_request" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AnimeRouletteID       *uint       `json:"anime_roulette_id"`
	EnigmaticTitle        string      `json:"enigmatic_title"` // Загаданная тема для аниме рулетки
}

type User_CreateJSON struct {
	UserTgID   int64  `json:"user_tg_id"`
	UserName   string `json:"user_name"`
	FullTgName string `json:"full_tg_name"`
}

type User_ReadJSON struct {
	ID                    uint              `json:"id"`
	CreatedAt             time.Time         `json:"created_at"`
	Step                  int               `json:"step"`
	UserTgID              int64             `json:"user_tg_id"`
	LastMessageID         int               `json:"last_message_id"`
	UserName              string            `json:"user_name"`
	FullTgName            string            `json:"full_tg_name"`
	Gender                Gender            `json:"gender"`
	IsVisitedEvents       bool              `json:"is_visited_events"`
	ISU                   string            `json:"isu"`
	FullName              string            `json:"full_name"`
	PhoneNumber           string            `json:"phone_number"`
	SecretCode            string            `json:"secret_code"`
	IsITMO                bool              `json:"is_itmo"`
	IsClubMember          bool              `json:"is_club_member"`
	IsSubscribeNewsletter bool              `json:"is_subscribe_newsletter"`
	IsSentRequest         bool              `json:"is_sent_request"`
	IsFilledData          bool              `json:"is_filled_data"`
	TempActivityID        int               `json:"temp_activity_id"`
	MyActivities          []*Activity       `json:"my_activities"`
	LinkMyAnimeList       string            `json:"link_my_anime_list"`
	MyRequest             *Request_ReadJSON `json:"my_request"`
	EnigmaticTitle        string            `json:"enigmatic_title"`
}

// Добавить пользователя
func DB_CREATE_User(user_to_add *User_CreateJSON) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user User
	db.Where("user_tg_id = ?", user_to_add.UserTgID).First(&user)
	if user.ID != 0 {
		return DB_ANSWER_OBJECT_EXISTS
	}

	user = User{
		UserTgID:              user_to_add.UserTgID,
		UserName:              user_to_add.UserName,
		FullTgName:            user.FullTgName,
		Step:                  config.STEP_DEFAULT,
		IsClubMember:          false,
		IsSubscribeNewsletter: false,
		IsFilledData:          false,
	}

	db.Save(&user)
	return DB_ANSWER_SUCCESS
}

// Получить пользователя по TgID
func DB_GET_User_BY_UserTgID(user_tg_id int64) (int, *User_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	user := new(User)
	db.Preload("MyRequest").Preload("MyActivities").Where("user_tg_id = ?", user_tg_id).First(&user)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	user_read := User_ReadJSON{
		ID:                    user.ID,
		CreatedAt:             user.CreatedAt,
		Step:                  user.Step,
		UserTgID:              user.UserTgID,
		Gender:                user.Gender,
		IsVisitedEvents:       user.IsVisitedEvents,
		LastMessageID:         user.LastMessageID,
		UserName:              user.UserName,
		FullTgName:            user.FullTgName,
		ISU:                   user.ISU,
		FullName:              user.FullName,
		PhoneNumber:           user.PhoneNumber,
		SecretCode:            user.SecretCode,
		IsITMO:                user.IsITMO,
		IsClubMember:          user.IsClubMember,
		IsSubscribeNewsletter: user.IsSubscribeNewsletter,
		IsSentRequest:         user.IsSentRequest,
		IsFilledData:          user.IsFilledData,
		TempActivityID:        user.TempActivityID,
		MyActivities:          user.MyActivities,
		LinkMyAnimeList:       user.LinkMyAnimeList,
		EnigmaticTitle:        user.EnigmaticTitle,
	}

	return DB_ANSWER_SUCCESS, &user_read
}

func DB_GET_Users_BY_Step(step int) []User_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var users []User
	db.Preload("MyRequest").Preload("MyActivities").Where("step = ?", step).Find(&users)

	users_read := make([]User_ReadJSON, len(users))

	for i, user := range users {
		users_read[i] = User_ReadJSON{
			ID:                    user.ID,
			CreatedAt:             user.CreatedAt,
			Step:                  user.Step,
			UserTgID:              user.UserTgID,
			Gender:                user.Gender,
			IsVisitedEvents:       user.IsVisitedEvents,
			LastMessageID:         user.LastMessageID,
			UserName:              user.UserName,
			FullTgName:            user.FullTgName,
			ISU:                   user.ISU,
			FullName:              user.FullName,
			PhoneNumber:           user.PhoneNumber,
			SecretCode:            user.SecretCode,
			IsITMO:                user.IsITMO,
			IsClubMember:          user.IsClubMember,
			IsSubscribeNewsletter: user.IsSubscribeNewsletter,
			IsSentRequest:         user.IsSentRequest,
			IsFilledData:          user.IsFilledData,
			TempActivityID:        user.TempActivityID,
			MyActivities:          user.MyActivities,
			LinkMyAnimeList:       user.LinkMyAnimeList,
			EnigmaticTitle:        user.EnigmaticTitle,
		}
	}

	return users_read
}

// Получить список всех пользователей
func DB_GET_Users() []User_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var users []User

	// Загружаем связанные сущности MyActivities
	db.Preload(clause.Associations).Find(&users)

	users_list := make([]User_ReadJSON, 0)
	if len(users) <= 0 {
		return users_list
	}

	for _, user := range users {
		var request *Request_ReadJSON
		if user.MyRequest != nil {
			request = &Request_ReadJSON{
				ID:        user.MyRequest.ID,
				CreatedAt: user.MyRequest.CreatedAt,
				Type:      user.MyRequest.Type,
				Status:    user.MyRequest.Status,
				UserID:    user.MyRequest.UserID,
			}
		}

		current_user := User_ReadJSON{
			ID:                    user.ID,
			CreatedAt:             user.CreatedAt,
			Step:                  user.Step,
			UserTgID:              user.UserTgID,
			LastMessageID:         user.LastMessageID,
			UserName:              user.UserName,
			FullTgName:            user.FullTgName,
			ISU:                   user.ISU,
			FullName:              user.FullName,
			PhoneNumber:           user.PhoneNumber,
			SecretCode:            user.SecretCode,
			IsITMO:                user.IsITMO,
			IsClubMember:          user.IsClubMember,
			IsSubscribeNewsletter: user.IsSubscribeNewsletter,
			IsSentRequest:         user.IsSentRequest,
			IsFilledData:          user.IsFilledData,
			TempActivityID:        user.TempActivityID,
			MyActivities:          user.MyActivities,
			MyRequest:             request,
			LinkMyAnimeList:       user.LinkMyAnimeList,
			EnigmaticTitle:        user.EnigmaticTitle,
		}
		users_list = append(users_list, current_user)
	}

	return users_list
}

// Обновляем пользователя
func DB_UPDATE_User(update_json map[string]interface{}) (int, *User) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user User
	user_tg_id, ok := update_json["user_tg_id"].(int64)
	if ok {
		db.Where("user_tg_id = ?", user_tg_id).First(&user)
		if user.ID == 0 {
			return DB_ANSWER_OBJECT_NOT_FOUND, nil
		}
	} else {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	// Обновляем поля, если они присутствуют в карте
	for key, value := range update_json {
		switch key {
		case "step":
			if v, ok := value.(int); ok && v != user.Step {
				user.Step = v
			}
		case "last_message_id":
			if v, ok := value.(int); ok && v != user.LastMessageID {
				user.LastMessageID = v
			}
		case "isu":
			if v, ok := value.(string); ok && v != user.ISU {
				user.ISU = v
			}
		case "full_name":
			if v, ok := value.(string); ok && v != user.FullName {
				user.FullName = v
			}
		case "gender":
			if v, ok := value.(Gender); ok && v != user.Gender {
				user.Gender = v
			} else if v, ok := value.(string); ok && v != string(user.Gender) {
				user.Gender = Gender(v)
			}
		case "is_visited_events":
			if v, ok := value.(bool); ok && v != user.IsVisitedEvents {
				user.IsVisitedEvents = v
			}
		case "phone_number":
			if v, ok := value.(string); ok && v != user.PhoneNumber {
				user.PhoneNumber = v
			}
		case "secret_code":
			if v, ok := value.(string); ok && v != user.SecretCode {
				user.SecretCode = v
			}

		case "link_my_anime_list":
			if v, ok := value.(string); ok && v != user.LinkMyAnimeList {
				user.LinkMyAnimeList = v
			}

		case "enigmatic_title":
			if v, ok := value.(string); ok && v != user.EnigmaticTitle {
				user.EnigmaticTitle = v
			}

		case "is_itmo":
			if v, ok := value.(bool); ok && v != user.IsITMO {
				user.IsITMO = v
			}

		case "is_club_member":
			if v, ok := value.(bool); ok && v != user.IsClubMember {
				user.IsClubMember = v
			}

		case "is_subscribe_newsletter":
			if v, ok := value.(bool); ok && v != user.IsSubscribeNewsletter {
				user.IsSubscribeNewsletter = v
			}

		case "is_sent_request":
			if v, ok := value.(bool); ok && v != user.IsSentRequest {
				user.IsSentRequest = v
			}

		case "is_filled_data":
			if v, ok := value.(bool); ok && v != user.IsFilledData {
				user.IsFilledData = v
			}

		case "temp_activity_id":
			if v, ok := value.(int); ok && v != user.TempActivityID {
				user.TempActivityID = v
			}
		}
	}

	db.Save(&user)
	return DB_ANSWER_SUCCESS, &user
}

// Удаление пользователя по tg_id
func DB_DELETE_User_BY_UserTgID(user_tg_id int64) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user User
	db.Where("user_tg_id = ?", user_tg_id).First(&user)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	result := db.Unscoped().Delete(&user)
	if result.Error != nil {
		rr_debug.PrintLOG("Users.go", "DB_DELETE_User_BY_UserTgID", "Error deleting user:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}

// Удалить всех пользователей
func DB_DELETE_Users() int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Exec("DELETE FROM users")
	if result.Error != nil {
		rr_debug.PrintLOG("Users.go", "DB_DELETE_Users", "Error deleting users:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}
