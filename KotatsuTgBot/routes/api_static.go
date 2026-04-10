// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

//
// ----------------------------------------------------------------------------------
//
// 								Static (Пути)
//
// ----------------------------------------------------------------------------------
//

package routes

import (
	//Внутренние пакеты проекта

	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	//Системные пакеты

	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

// ----------------------------------------------
//
// 				Root requests
//
// ----------------------------------------------

func Handler_NewAdminPanel(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.File(config.ByUI("./static/dist/index.html"))
}

func Handler_Login(c *gin.Context) {
	c.SetCookie("session_token", c.Request.URL.RawQuery, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin")
}

func Handler_GetCalendarActivities_Image_File(c *gin.Context) {
	if _, err := os.Stat(config.ByUI("./img/calendar_activities/calendar_activities.png")); err == nil {
		Answer_File(c, "/img/calendar_activities/calendar_activities.png")
	} else if _, err := os.Stat(config.ByUI("./img/calendar_activities/calendar_activities.jpg")); err == nil {
		Answer_File(c, "/img/calendar_activities/calendar_activities.jpg")
	} else {
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
	}
}

func Handler_UploadFile_CalendarActivities(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_UploadFile_CalendarActivities", "c.FormFile", "Неверные данные в запросе", err.Error())
		Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		return
	}

	err_rem := removeAllContents(config.ByUI("./img/calendar_activities"))
	if err_rem != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_UploadFile_CalendarActivities", "removeAllContents('./img/calendar_activities')", "Ошибка при удалении каталога для календаря мероприятий", err_rem.Error())
	}

	extension := filepath.Ext(file.Filename)

	uploadDir := config.ByUI("./img/calendar_activities")
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, "calendar_activities"+extension)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Answer_BadRequest(c, ANSWER_INVALID_FILE_UPLOAD().Code, ANSWER_INVALID_FILE_UPLOAD().Message)
		return
	}

	Answer_OK(c)
}

type SendBroadcast struct {
	Events           []int64  `json:"events"`
	Users            []int64  `json:"users"`
	Roulettes        []int64  `json:"roulettes"`
	ClubMemberStatus *bool    `json:"club_member_status"`
	ItmoStatus       []string `json:"itmo_status"`
	Message          string   `json:"message"`
}

type BroadcastResult struct {
	User         db.User_ReadJSON `json:"user"`
	Success      bool             `json:"success"`
	ErrorMessage string           `json:"error_message"`
}

func Handler_API_SendBroadcast(c *gin.Context) {
	json_data := new(SendBroadcast)
	err_json_bin := c.ShouldBindJSON(&json_data)

	if err_json_bin != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_API_SendBroadcast", "c.ShouldBindJSON", "Неверные данные в запросе", err_json_bin.Error())
		Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err_json_bin.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{}

	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN, opts...)
	if err != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_API_SendBroadcast", "bot.New", "Ошибка инициализации бота", err.Error())
		Answer_BadRequest(c, ANSWER_BOT_CONNECT_ERROR("").Code, ANSWER_BOT_CONNECT_ERROR("").Message+" Error: "+err.Error())
		return
	}

	search_params := db.UserSearchParams{
		Events:           json_data.Events,
		Users:            json_data.Users,
		Roulettes:        json_data.Roulettes,
		ClubMemberStatus: json_data.ClubMemberStatus,
		ItmoStatus:       json_data.ItmoStatus,
	}

	filtered_users := db.DB_User_Search(search_params)
	if len(filtered_users) == 0 {
		Answer_NotFound(c, ANSWER_NO_USERS_FOUND().Code, ANSWER_NO_USERS_FOUND().Message)
		return
	}

	results := make([]BroadcastResult, 0, len(filtered_users))

	for _, current_user := range filtered_users {
		params := &bot.SendMessageParams{
			ChatID:    current_user.UserTgID,
			ParseMode: models.ParseModeHTML,
			Text:      config.TT("broadcast", json_data.Message),
		}

		_, err_send := b.SendMessage(ctx, params)
		if err_send != nil {
			rr_debug.PrintLOG("api_static.go", "Handler_API_SendBroadcast", "bot.SendMessage", "Ошибка отправки сообщения пользователю", err_send.Error())
			// Record failure for this user
			results = append(results, BroadcastResult{
				User:         current_user,
				Success:      false,
				ErrorMessage: err_send.Error(),
			})
		} else {
			// Record success for this user
			results = append(results, BroadcastResult{
				User:         current_user,
				Success:      true,
				ErrorMessage: "",
			})
		}
	}

	// Return detailed results for all users
	Answer_SendObject(c, gin.H{
		"results": results,
	})
}
