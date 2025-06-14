package db

import (
	//Внутренние пакеты проекта

	"database/sql/driver"
	"encoding/json"
	"errors"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"gorm.io/gorm"

	//Системные пакеты
	"time"
)

type AnimeRoulette struct {
	gorm.Model
	Status       bool           `json:"status"`                   // Прошла или не прошла в целом
	Stages       RouletteStages `gorm:"type:jsonb" json:"stages"` // Этапы рулетки
	CurrentStage int            `json:"current_stage"`            // Текущий этап рулетки
	Theme        string         `json:"theme"`                    // Тема рулетки
	Participants []User         `json:"participants"`             // Участники рулетки
}

type RouletteStages []RouletteStage

func (a RouletteStages) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *RouletteStages) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type RouletteStage struct {
	Stage     int       `json:"stage"`      // Этап
	StartDate time.Time `json:"start_date"` // Дата начала этапа рулетки
	EndDate   time.Time `json:"end_date"`   // Дата окончания этапа рулетки
}

type AnimeRoulette_CreateJSON struct {
	Stages []RouletteStage `json:"stages"`
}

type AnimeRoulette_ReadJSON struct {
	ID           uint            `json:"id"`
	CreatedAt    time.Time       `json:"created_at"`
	Status       bool            `json:"status"`
	Stages       []RouletteStage `json:"stages"`
	CurrentStage int             `json:"current_stage"`
	Theme        string          `json:"theme"`
	Participants []User          `json:"participants"`
}

// Добавить аниме рулетку
func DB_CREATE_AnimeRoulette(anime_roulette_to_add *AnimeRoulette_CreateJSON) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette
	db.First(&anime_roulette)
	if anime_roulette.ID != 0 {
		return DB_ANSWER_OBJECT_EXISTS
	}

	anime_roulette = AnimeRoulette{
		Stages:       anime_roulette_to_add.Stages,
		Status:       true,
		CurrentStage: config.ANIME_RUOLETTE_STAGE_START_REGISTRATION,
	}

	db.Save(&anime_roulette)
	return DB_ANSWER_SUCCESS
}

// Получить аниме рулетку по Theme
func DB_GET_AnimeRoulette_BY_Theme(theme string) (int, *AnimeRoulette_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	anime_roulette := new(AnimeRoulette)
	db.Preload("Participants").Where("theme = ?", theme).First(&anime_roulette)
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	anime_roulette_read := AnimeRoulette_ReadJSON{
		ID:           anime_roulette.ID,
		CreatedAt:    anime_roulette.CreatedAt,
		Theme:        anime_roulette.Theme,
		Status:       anime_roulette.Status,
		CurrentStage: anime_roulette.CurrentStage,
		Participants: anime_roulette.Participants,
	}

	return DB_ANSWER_SUCCESS, &anime_roulette_read
}

// Получить аниме рулетку по Status
func DB_GET_AnimeRoulette_BY_Status(status bool) (int, *AnimeRoulette_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	anime_roulette := new(AnimeRoulette)
	db.Preload("Participants").Where("status = ?", status).First(&anime_roulette)
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	anime_roulette_read := AnimeRoulette_ReadJSON{
		ID:           anime_roulette.ID,
		CreatedAt:    anime_roulette.CreatedAt,
		Theme:        anime_roulette.Theme,
		Status:       anime_roulette.Status,
		CurrentStage: anime_roulette.CurrentStage,
		Stages:       anime_roulette.Stages,
		Participants: anime_roulette.Participants,
	}

	return DB_ANSWER_SUCCESS, &anime_roulette_read
}

// Получить список аниме рулеток
func DB_GET_AnimeRoulettes() []AnimeRoulette_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulettes []AnimeRoulette

	db.Preload("Participants").Find(&anime_roulettes)

	anime_roulettes_list := make([]AnimeRoulette_ReadJSON, 0)
	if len(anime_roulettes) <= 0 {
		return anime_roulettes_list
	}

	for _, anime_roulette := range anime_roulettes {

		current_anime_roulette := AnimeRoulette_ReadJSON{
			ID:           anime_roulette.ID,
			CreatedAt:    anime_roulette.CreatedAt,
			Theme:        anime_roulette.Theme,
			Status:       anime_roulette.Status,
			CurrentStage: anime_roulette.CurrentStage,
			Participants: anime_roulette.Participants,
		}
		anime_roulettes_list = append(anime_roulettes_list, current_anime_roulette)
	}

	return anime_roulettes_list
}

// Обновляем аниме рулетку
func DB_UPDATE_AnimeRoulette(update_json map[string]interface{}) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette

	db.First(&anime_roulette)
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	// Обновляем поля, если они присутствуют в карте
	for key, value := range update_json {
		switch key {
		case "status":
			if v, ok := value.(bool); ok && v != anime_roulette.Status {
				anime_roulette.Status = v
			}

		case "current_stage":
			if v, ok := value.(int); ok && v != anime_roulette.CurrentStage {
				anime_roulette.CurrentStage = v
			}

		case "theme":
			if v, ok := value.(string); ok && v != anime_roulette.Theme {
				anime_roulette.Theme = v
			}

		case "stage_new_date":
			if stage_new_date, ok := update_json["stage_new_date"].(map[string]interface{}); ok {
				if stage_date, ok := stage_new_date["stage_date"].(map[string]interface{}); ok {
					if stage_v, ok := stage_date["stage"].(int); ok {
						for i, stage := range anime_roulette.Stages {
							if stage.Stage == stage_v {

								if end_date_v, ok := stage_date["end_date"].(string); ok {
									// Формат строки даты и времени
									layout := "2006-01-02 15:04"

									// Парсим строку в time.Time
									end_date_stage, err_time := time.Parse(layout, end_date_v)
									if err_time != nil {
										rr_debug.PrintLOG("api_anime_roulettes.go", "DB_UPDATE_AnimeRoulette", "DateMeeting Parse", "Ошибка при парсинге времени", err_time.Error())
									}

									anime_roulette.Stages[i].EndDate = end_date_stage
								}
							}
						}
					}
				}
			}
		}
	}

	db.Save(&anime_roulette)
	return DB_ANSWER_SUCCESS
}

// Добавляем пользователей в аниме рулетку
func DB_UPDATE_AnimeRoulette_ADD_Participants(user_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette_ReadJSON

	db.First(&anime_roulette)
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	var user User
	db.First(&user, user_id)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	// Добавить пользователя в массив Participants
	anime_roulette.Participants = append(anime_roulette.Participants, user)

	db.Save(&anime_roulette)
	return DB_ANSWER_SUCCESS
}

// Удаляем пользователя из аниме рулетки
func DB_UPDATE_AnimeRoulette_REMOVE_Participants(user_id uint) int {
	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette_ReadJSON

	db.First(&anime_roulette)
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	// Ищем индекс пользователя в массиве Participants
	userIndex := -1
	for i, participant := range anime_roulette.Participants {
		if participant.ID == user_id {
			userIndex = i
			break
		}
	}

	if userIndex == -1 {
		return DB_ANSWER_OBJECT_NOT_FOUND // Пользователь не найден в рулетке
	}

	// Удаляем пользователя из массива Participants
	anime_roulette.Participants = append(anime_roulette.Participants[:userIndex], anime_roulette.Participants[userIndex+1:]...)

	db.Save(&anime_roulette)
	return DB_ANSWER_SUCCESS
}

// Удаление аниме рулетку по theme
func DB_DELETE_AnimeRoulette_BY_Theme(theme string) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette
	db.Where("theme = ?", theme).First(&anime_roulette)
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	result := db.Unscoped().Delete(&anime_roulette)
	if result.Error != nil {
		rr_debug.PrintLOG("AnimeRoulettes.go", "DB_DELETE_AnimeRoulette_BY_Theme", "Error deleting anime_roulette:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}

// Удалить все аниме рулетки
func DB_DELETE_AnimeRoulettes() int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Exec("DELETE FROM anime_roulettes")
	if result.Error != nil {
		rr_debug.PrintLOG("AnimeRoulettes.go", "DB_DELETE_AnimeRoulette_BY_Theme", "Error deleting anime_roulette:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}
