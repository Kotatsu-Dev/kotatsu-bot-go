package db

import (
	//Внутренние пакеты проекта

	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"github.com/lib/pq"
	"gorm.io/gorm"

	//Системные пакеты
	"time"
)

type AnimeRoulette struct {
	gorm.Model
	StartDate        time.Time      `json:"start_date"`                                                                   // Дата начала рулетки
	AnnounceDate     time.Time      `json:"announce_date"`                                                                // Дата объявления темы
	DistributionDate time.Time      `json:"distribution_date"`                                                            // Дата распределения тайтлов
	EndDate          time.Time      `json:"end_date"`                                                                     // Дата окончания
	Theme            string         `json:"theme"`                                                                        // Тема рулетки
	Participants     []User         `json:"participants" gorm:"foreignKey:AnimeRouletteID;constraint:OnDelete:SET NULL;"` // Участники рулетки
	Distribution     *pq.Int32Array `json:"distribution" gorm:"type:integer[]"`                                           // Распределение участников: ID пользователей по кругу, каждый получает тайтл следующего
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
	StartDate        time.Time `json:"start_date"`
	AnnounceDate     time.Time `json:"announce_date"`
	DistributionDate time.Time `json:"distribution_date"`
	EndDate          time.Time `json:"end_date"`
	Theme            *string   `json:"theme"`
}

type AnimeRoulette_ReadJSON struct {
	ID               uint            `json:"id"`
	CreatedAt        time.Time       `json:"created_at"`
	StartDate        time.Time       `json:"start_date"`
	AnnounceDate     time.Time       `json:"announce_date"`
	DistributionDate time.Time       `json:"distribution_date"`
	EndDate          time.Time       `json:"end_date"`
	Theme            string          `json:"theme"`
	Participants     []User_ReadJSON `json:"participants"`
	Distribution     *pq.Int32Array  `json:"distribution"`
}

func (roulette *AnimeRoulette) ToRead() *AnimeRoulette_ReadJSON {
	return &AnimeRoulette_ReadJSON{
		ID:               roulette.ID,
		CreatedAt:        roulette.CreatedAt,
		Theme:            roulette.Theme,
		StartDate:        roulette.StartDate,
		AnnounceDate:     roulette.AnnounceDate,
		DistributionDate: roulette.DistributionDate,
		EndDate:          roulette.EndDate,
		Participants:     UserToReadSlice(roulette.Participants),
		Distribution:     roulette.Distribution,
	}
}

func RouletteToReadSlice(roulettes []AnimeRoulette) []AnimeRoulette_ReadJSON {
	res := make([]AnimeRoulette_ReadJSON, len(roulettes))
	for i, roulette := range roulettes {
		res[i] = *roulette.ToRead()
	}
	return res
}

// Добавить аниме рулетку
func DB_CREATE_AnimeRoulette(anime_roulette_to_add *AnimeRoulette_CreateJSON) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette

	var theme string

	if anime_roulette_to_add.Theme != nil {
		theme = *anime_roulette_to_add.Theme
	}

	anime_roulette = AnimeRoulette{
		StartDate:        anime_roulette_to_add.StartDate,
		AnnounceDate:     anime_roulette_to_add.AnnounceDate,
		DistributionDate: anime_roulette_to_add.DistributionDate,
		EndDate:          anime_roulette_to_add.EndDate,
		Theme:            theme,
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

	return DB_ANSWER_SUCCESS, anime_roulette.ToRead()
}

func DB_GET_AnimeRoulette_BY_Status(status bool) (int, *AnimeRoulette_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	anime_roulette := new(AnimeRoulette)
	now := time.Now()
	if status {
		db.Preload("Participants").Where("start_date < ?", now).Where("end_date > ?", now).First(&anime_roulette)
	} else {
		// Последняя завершившаяся рулетка (будущие рулетки сюда не попадают)
		db.Preload("Participants").Where("end_date < ?", now).Order("end_date desc").First(&anime_roulette)
	}
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	return DB_ANSWER_SUCCESS, anime_roulette.ToRead()
}

// Получить список аниме рулеток
func DB_GET_AnimeRoulettes() []AnimeRoulette_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulettes []AnimeRoulette

	db.Preload("Participants").Find(&anime_roulettes)

	return RouletteToReadSlice(anime_roulettes)
}

func parseDate(value any, field *time.Time) {
	if v, ok := value.(string); ok {
		v_date, err_time := time.Parse(time.RFC3339, v)
		if err_time != nil {
			rr_debug.PrintLOG("AnimeRoulettes.go", "parseDate", "DateMeeting Parse", "Ошибка при парсинге времени", err_time.Error())
		} else {
			*field = v_date
		}
	}
}

func DB_UPDATE_AnimeRoulette(update_json map[string]any) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette

	if roulette_id, ok := update_json["id"].(float64); ok {
		db.First(&anime_roulette, int64(roulette_id))
	} else {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	for key, value := range update_json {
		switch key {
		case "status":
			return DB_ANSWER_NOT_SUPPORTED

		case "current_stage":
			return DB_ANSWER_NOT_SUPPORTED

		case "theme":
			if v, ok := value.(string); ok && v != anime_roulette.Theme {
				anime_roulette.Theme = v
			}

		case "stage_new_date":
			return DB_ANSWER_NOT_SUPPORTED

		case "start_date":
			parseDate(value, &anime_roulette.StartDate)

		case "announce_date":
			parseDate(value, &anime_roulette.AnnounceDate)

		case "distribution_date":
			parseDate(value, &anime_roulette.DistributionDate)

		case "end_date":
			parseDate(value, &anime_roulette.EndDate)
		}
	}

	db.Save(&anime_roulette)
	return DB_ANSWER_SUCCESS
}

func DB_UPDATE_AnimeRoulette_SET_Distribution(roulette_id uint, distr []int32) int {
	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette

	db.First(&anime_roulette, roulette_id)
	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	anime_roulette.Distribution = (*pq.Int32Array)(&distr)

	db.Save(&anime_roulette)
	return DB_ANSWER_SUCCESS
}

// Добавляем пользователей в аниме рулетку
func DB_UPDATE_AnimeRoulette_ADD_Participants(user_id uint) int {
	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette
	now := time.Now()
	db.Preload("Participants").Where("start_date < ?", now).Where("end_date > ?", now).First(&anime_roulette)

	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	var user User
	db.First(&user, user_id)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	res := db.Model(&anime_roulette).Association("Participants").Append(&user)
	fmt.Println(res)
	return DB_ANSWER_SUCCESS
}

// Удаляем пользователя из аниме рулетки
func DB_UPDATE_AnimeRoulette_REMOVE_Participants(user_id uint) int {
	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var anime_roulette AnimeRoulette

	now := time.Now()
	db.Preload("Participants").Where("start_date < ?", now).Where("end_date > ?", now).First(&anime_roulette)

	if anime_roulette.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	var user User
	db.First(&user, user_id)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	db.Model(&anime_roulette).Association("Participants").Delete(&user)

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
