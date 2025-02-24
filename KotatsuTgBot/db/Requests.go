package db

import (
	//Внутренние пакеты проекта
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"gorm.io/gorm"

	//Системные пакеты
	"time"
)

type Request struct {
	gorm.Model
	Type   int  `json:"type"`    // Тип запроса: 1 = Вступить в клуб
	Status int  `json:"status"`  // Статус заявки: 0 = новая, 1 - Принята, 2 - Отклонена
	UserID uint `json:"user_id"` // Информация о пользователе
}

type Request_CreateJSON struct {
	UserID uint `json:"user_id"`
}

type Request_ReadJSON struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Type      int       `json:"type"`
	Status    int       `json:"status"`
	UserID    uint      `json:"user_id"`
}

// Добавить заявку
func DB_CREATE_Request(user_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var request Request
	db.Where("user_id = ?", request.UserID).First(&request)
	if request.ID != 0 {
		return DB_ANSWER_OBJECT_EXISTS
	}

	request = Request{
		Type:   1,
		Status: 0,
		UserID: user_id,
	}

	db.Save(&request)
	return DB_ANSWER_SUCCESS
}

// Получить заявку по ID
func DB_GET_Request_BY_ID(request_id uint) (int, *Request_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	request := new(Request)
	db.First(&request, request_id)
	if request.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	request_read := Request_ReadJSON{
		ID:        request.ID,
		CreatedAt: request.CreatedAt,
		Type:      request.Type,
		Status:    request.Status,
		UserID:    request.UserID,
	}

	return DB_ANSWER_SUCCESS, &request_read
}

// Получить список всех заявки
func DB_GET_Requests() []Request_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var requests []Request

	// Загружаем связанные сущности MyActivities
	db.Find(&requests)

	requests_list := make([]Request_ReadJSON, 0)
	if len(requests) <= 0 {
		return requests_list
	}

	for _, request := range requests {

		current_request := Request_ReadJSON{
			ID:        request.ID,
			CreatedAt: request.CreatedAt,
			Type:      request.Type,
			Status:    request.Status,
			UserID:    request.UserID,
		}
		requests_list = append(requests_list, current_request)
	}

	return requests_list
}

// Обновляем заявку
func DB_UPDATE_Request(update_json map[string]interface{}) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var request Request
	request_id, ok := update_json["request_id"].(uint)
	if ok {
		db.First(&request, request_id)
	} else {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	// Обновляем поля, если они присутствуют в карте
	for key, value := range update_json {
		switch key {
		case "type":
			if v, ok := value.(int); ok && v != request.Type {
				request.Type = v
			}
		case "status":
			if v, ok := value.(int); ok && v != request.Status {
				request.Status = v
			}
		}
	}

	db.Save(&request)
	return DB_ANSWER_SUCCESS
}

// Одобряем или отклоняем заявку
func DB_UPDATE_Choise_Request(update_json map[string]interface{}) (int, *User) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var request Request
	var user User

	request_id, ok := update_json["request_id"].(uint)
	if ok {
		db.First(&request, request_id)
	} else {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	// Загрузка пользователя, связанного с заявкой
	db.Model(&request).Association("User").Find(&user)

	if status, ok := update_json["status"].(int); ok {
		if status == 1 {
			user.IsClubMember = true
			user.IsSentRequest = false
			db.Save(&user)
		}
	}

	db.Save(&request)
	return DB_ANSWER_SUCCESS, &user
}

// Удаление заявки по id
func DB_DELETE_Request_BY_ID(request_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Unscoped().Delete(&Request{}, request_id)
	if result.Error != nil {
		rr_debug.PrintLOG("Requests.go", "DB_DELETE_Request_BY_ID", "Error deleting request:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}

// Удалить все заявок
func DB_DELETE_Requests() int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Exec("DELETE FROM requests")
	if result.Error != nil {
		rr_debug.PrintLOG("Requests.go", "DB_DELETE_Requests", "Error deleting request:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}
