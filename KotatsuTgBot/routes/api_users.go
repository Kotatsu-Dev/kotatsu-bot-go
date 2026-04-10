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

	db_answer_code, user, has_changed_membership := db.DB_UPDATE_User(update_json)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		if has_changed_membership {
			err := SendMembershipUpdate(user.UserTgID, user.IsClubMember)
			if err != nil {
				rr_debug.PrintLOG("api_static.go", "Handler_API_Users_UpdateObject_ClubMember", "bot.Send", "Ошибка отправки сообщения", err.Error())
				Answer_BadRequest(c, ANSWER_BOT_SEND_MESSAGE_ERROR("").Code, ANSWER_BOT_SEND_MESSAGE_ERROR("").Message+" Error: "+err.Error())
			}
		}
		Answer_OK(c)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
		return

	default:
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}
}

func SendMembershipUpdate(user_tg_id int64, is_member bool) error {
	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN)
	if err != nil {
		return err
	}

	params := &bot.SendMessageParams{
		ChatID:    user_tg_id,
		ParseMode: models.ParseModeHTML,
	}

	if is_member {
		params.Text = config.T("member_welcome")
		params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsClubMember
	} else {
		params.Text = config.T("member_kick")
		params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsDefault
	}

	_, err_send := b.SendMessage(context.TODO(), params)
	return err_send
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
