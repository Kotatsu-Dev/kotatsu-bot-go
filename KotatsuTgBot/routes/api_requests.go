// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

//
// ----------------------------------------------------------------------------------
//
// 								Requests (Пути)
//
// ----------------------------------------------------------------------------------
//

package routes

import (

	//Внутренние пакеты проекта

	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"

	//Системные пакеты
	"context"
	"os"
	"os/signal"
)

// Получить все заявки
func Handler_API_Requests_GetList(c *gin.Context) {
	list_requests := db.DB_GET_Requests()
	answer := GetList_Requests_Answer{
		ListRequests: list_requests,
	}

	Answer_SendObject(c, answer)
	return
}

// Обновить данные в заявке
func Handler_API_Requests_UpdateObject(c *gin.Context) {

	var update_json map[string]interface{}

	err := c.ShouldBindJSON(&update_json)
	if err != nil {
		rr_debug.PrintLOG("api_requests.go", "Handler_API_Requests_UpdateObject", "c.ShouldBindJSON", "Неверные данные в запросе", err.Error())
		if config.GetConfig().CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	}

	db_answer_code := db.DB_UPDATE_Request(update_json)
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

// Удалить все заявки
func Handler_API_Requests_DeleteObject_ALL(c *gin.Context) {
	db_answer_code := db.DB_DELETE_Requests()
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

// Одобрить или отклонить заявку
func Handler_API_Requests_UpdateObject_Choise(c *gin.Context) {

	var update_json map[string]interface{}
	err := c.ShouldBindJSON(&update_json)
	if err != nil {
		rr_debug.PrintLOG("api_requests.go", "Handler_API_Requests_UpdateObject", "c.ShouldBindJSON", "Неверные данные в запросе", err.Error())
		if config.GetConfig().CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	}

	status_, ok := update_json["status"].(float64)
	if !ok {
		Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		return
	}
	status := int(status_)

	request_id_, ok := update_json["request_id"].(float64)
	if !ok {
		Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		return
	}
	request_id := uint(request_id_)

	db_answer_code, user := db.DB_UPDATE_Choise_Request(map[string]interface{}{"status": status, "request_id": request_id})

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		opts := []bot.Option{}

		b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN, opts...)
		if err != nil {
			rr_debug.PrintLOG("api_requests.go", "Handler_API_Requests_UpdateObject_Choise", "gotgbot.NewBot", "Ошибка инициализации бота", err.Error())
		}

		params := &bot.SendMessageParams{
			ChatID: user.UserTgID,
		}

		if status == 1 {
			params.Text = "Добро пожаловать в аниме-клуб «Котацу» — твоя заявка на вступление успешно обработана!"
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(user.IsSubscribeNewsletter)
		} else {
			params.Text = "К сожалению, клуб отклонил твою заявку на вступление" + "\n" +
				"Возможно, ты указал неверные данные. Чтобы узнать причину, свяжись с клубом через кнопку в меню бота"
		}

		// Удаляем заявку по ID
		db_answer_code_request_del := db.DB_DELETE_Request_BY_ID(request_id)
		switch db_answer_code_request_del {
		case db.DB_ANSWER_SUCCESS:
			_, err := b.SendMessage(ctx, params)
			if err != nil {
				rr_debug.PrintLOG("api_requests.go", "Handler_API_Requests_UpdateObject_Choise", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
				Answer_BadRequest(c, ANSWER_BOT_SEND_MESSAGE_ERROR("").Code, ANSWER_BOT_SEND_MESSAGE_ERROR("").Message+" Error: "+err.Error())
				return
			} else {
				Answer_OK(c)
				return
			}

		case db.DB_ANSWER_DELETE_ERROR:
			Answer_BadRequest(c, ANSWER_DB_DELETE_OBJECT_FAILED().Code, ANSWER_DB_DELETE_OBJECT_FAILED().Message)
			return

		default:
			Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
			return
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
		return

	default:
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}
}
