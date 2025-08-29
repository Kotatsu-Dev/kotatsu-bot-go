// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

//
// ----------------------------------------------------------------------------------
//
// 								Users (Пути)
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
	"github.com/go-telegram/bot/models"

	//Системные пакеты
	"context"
	"os"
	"os/signal"
)

// Получить всех зарегистрированных пользователей
func Handler_API_Users_GetList(c *gin.Context) {

	list_users := db.DB_GET_Users()
	answer := GetList_Users_Answer{
		ListUsers: list_users,
	}

	Answer_SendObject(c, answer)
}

// Обновить данные пользователя
func Handler_API_Users_UpdateObject(c *gin.Context) {

	var update_json map[string]interface{}

	err := c.ShouldBindJSON(&update_json)
	if err != nil {
		rr_debug.PrintLOG("api_users.go", "Handler_API_Users_UpdateObject", "c.ShouldBindJSON", "Неверные данные в запросе", err.Error())
		if config.GetConfig().CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	}

	db_answer_code, _ := db.DB_UPDATE_User(update_json)
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

// Обновить существование пользователя в клубе
func Handler_API_Users_UpdateObject_ClubMember(c *gin.Context) {

	var update_json map[string]interface{}

	err := c.ShouldBindJSON(&update_json)
	if err != nil {
		rr_debug.PrintLOG("api_users.go", "Handler_API_Users_UpdateObject_ClubMember", "c.ShouldBindJSON", "Неверные данные в запросе", err.Error())
		if config.GetConfig().CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	} else {
		v, ok := update_json["user_tg_id"].(float64)
		if !ok {
			rr_debug.PrintLOG("api_users.go", "Handler_API_Users_UpdateObject_ClubMember", "c.ShouldBindJSON", "Неверные данные в запросе", "ok")
			if config.GetConfig().CONFIG_IS_DEBUG {
				Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: tg_user_id not expected")
			} else {
				Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
			}
			return
		} else {
			update_json["user_tg_id"] = int64(v)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{}

	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN, opts...)
	if err != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_SendMessageUser", "gotgbot.NewBot", "Ошибка инициализации бота", err.Error())
	}

	db_answer_code, user := db.DB_UPDATE_User(update_json)
	switch db_answer_code {

	case db.DB_ANSWER_SUCCESS:

		params := &bot.SendMessageParams{
			ChatID:    user.UserTgID,
			ParseMode: models.ParseModeHTML,
		}

		if user.IsClubMember {
			params.Text = "<b>Руководство добавило вас в клуб</b>" + "\n" +
				"Если у вас остались какие-либо вопросы, свяжитесь с руководством клуба по кнопке в Меню снизу"
		} else {
			params.Text = "<b>Руководство исключило вас из клуба</b>" + "\n" +
				"Если у вас остались какие-либо вопросы, свяжитесь с руководством клуба по кнопке в Меню снизу"
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(user.IsSubscribeNewsletter)
		}

		_, err_send := b.SendMessage(ctx, params)
		if err_send != nil {
			rr_debug.PrintLOG("api_static.go", "Handler_API_Users_UpdateObject_ClubMember", "bot.Send", "Ошибка отправки сообщения", err_send.Error())
			Answer_BadRequest(c, ANSWER_BOT_SEND_MESSAGE_ERROR(params.Text).Code, ANSWER_BOT_SEND_MESSAGE_ERROR(params.Text).Message+" Error: "+err_send.Error())
		}
		Answer_OK(c)
		return

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
		return

	default:
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}
}

// Удалить всех пользователей
func Handler_API_Users_DeleteObject_ALL(c *gin.Context) {

	db_answer_code := db.DB_DELETE_Users()

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
