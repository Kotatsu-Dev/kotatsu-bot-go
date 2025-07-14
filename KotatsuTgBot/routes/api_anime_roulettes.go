// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

// ----------------------------------------------------------------------------------
//
//	AnimeRoulettes (Пути)
//
// ----------------------------------------------------------------------------------
package routes

import (

	//Внутренние пакеты проекта
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"github.com/gin-gonic/gin"

	//Системные пакеты
	"time"
)

// Получить все рулетки
func Handler_API_AnimeRoulettes_GetList(c *gin.Context) {

	list_anime_roulette := db.DB_GET_AnimeRoulettes()
	answer := GetList_AnimeRoulettes_Answer{
		ListAnimeRoulettes: list_anime_roulette,
	}

	Answer_SendObject(c, answer)
	return
}

// Получить активную рулетку
func Handler_API_AnimeRoulettes_GetActive(c *gin.Context) {

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		Answer_SendObject(c, current_anime_roulette)
		return

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
		return

	default:
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}
}

// Создать рулетку
func Handler_API_AnimeRoulettes_CreateObject(c *gin.Context) {

	json_data := new(Create_AnimeRoulettes)
	err := c.ShouldBindJSON(&json_data)

	//Проверка, JSON пришел или шляпа
	if err != nil {
		rr_debug.PrintLOG("api_anime_roulettes.go", "Handler_API_AnimeRoulettes_CreateObject", "c.ShouldBindJSON", "Неверные данные в запросе", err.Error())
		if config.GetConfig().CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	}

	if len(json_data.Stages) < 4 {
		rr_debug.PrintLOG("api_anime_roulettes.go", "Handler_API_AnimeRoulettes_CreateObject", "json_data.Stages == ''", "Пустые данные в запросе", "")
		Answer_BadRequest(c, ANSWER_EMPTY_FIELDS().Code, ANSWER_EMPTY_FIELDS().Message)
		return
	}

	// Формат строки даты и времени
	layout := "2006-01-02 15:04"

	var start_date, announce_date, distribution_date, end_date time.Time

	for index, stage := range json_data.Stages {
		// Парсим строку в time.Time
		end_date_stage, err_time := time.Parse(layout, stage.EndDate)
		if err_time != nil {
			rr_debug.PrintLOG("api_anime_roulettes.go", "Handler_API_AnimeRoulettes_CreateObject", "DateMeeting Parse", "Ошибка при парсинге времени", err_time.Error())
			return
		}

		switch index {
		case 0:
			start_date = end_date_stage
		case 1:
			announce_date = end_date_stage
		case 2:
			distribution_date = end_date_stage
		case 3:
			end_date = end_date_stage
		}
	}

	current_anime_roulette := new(db.AnimeRoulette_CreateJSON)
	current_anime_roulette.StartDate = start_date
	current_anime_roulette.AnnounceDate = announce_date
	current_anime_roulette.DistributionDate = distribution_date
	current_anime_roulette.EndDate = end_date

	db_error_code := db.DB_CREATE_AnimeRoulette(current_anime_roulette)

	switch db_error_code {
	case db.DB_ANSWER_SUCCESS:
		Answer_OK(c)
		return

	case db.DB_ANSWER_OBJECT_EXISTS:
		Answer_BadRequest(c, ANSWER_OBJECT_EXISTS().Code, ANSWER_OBJECT_EXISTS().Message)
		return

	default:
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}
}

// Обновить данные рулетки
func Handler_API_AnimeRoulettes_UpdateObject(c *gin.Context) {
	var update_json map[string]interface{}

	err := c.ShouldBindJSON(&update_json)
	if err != nil {
		rr_debug.PrintLOG("api_anime_roulettes.go", "Handler_API_AnimeRoulettes_UpdateObject", "c.ShouldBindJSON", "Неверные данные в запросе", err.Error())
		if config.GetConfig().CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	}

	db_answer_code := db.DB_UPDATE_AnimeRoulette(update_json)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		Answer_OK(c)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
		return

	default:
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}
}

// Удалить рулетки
func Handler_API_AnimeRoulettes_DeleteObject_ALL(c *gin.Context) {

	db_answer_code := db.DB_DELETE_AnimeRoulettes()
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		Answer_OK(c)
		return

	case db.DB_ANSWER_DELETE_ERROR:
		Answer_BadRequest(c, ANSWER_DB_DELETE_OBJECT_FAILED().Code, ANSWER_DB_DELETE_OBJECT_FAILED().Message)
		return

	default:
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}
}
