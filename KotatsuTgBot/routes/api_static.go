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
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

// ----------------------------------------------
//
// 				Root requests
//
// ----------------------------------------------
//HTML-пути

// /Выдача фронта для получения лога
func Handler_Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// /Выдача фронта для панели Админа
func Handler_AdminPanel(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_panel.html", gin.H{})
}

func Handler_NewAdminPanel(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.File("./static/dist/index.html")
}

// /Выдача фронта для ответа от техподдержки
func Handler_SupportResponse(c *gin.Context) {
	c.HTML(http.StatusOK, "support_response.html", gin.H{})
}

func Handler_Login(c *gin.Context) {
	c.SetCookie("session_token", c.Request.URL.RawQuery, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/new-admin-panel")
}

// Получить изображение текущего календаря мероприятий
func Handler_GetCalendarActivities_Image_File(c *gin.Context) {
	// Проверка существования файла
	if _, err := os.Stat("./img/calendar_activities/calendar_activities.png"); err == nil {
		// Answer_OK(c)
		Answer_File(c, "/img/calendar_activities/calendar_activities.png")
	} else if _, err := os.Stat("./img/calendar_activities/calendar_activities.jpg"); err == nil {
		Answer_File(c, "/img/calendar_activities/calendar_activities.jpg")
	} else {
		// Файл не найден, возвращаем ошибку 404
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
	}
}

// SendMessageUser
func Handler_SendMessageUser(c *gin.Context) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{}

	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN, opts...)
	if err != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_SendMessageUser", "gotgbot.NewBot", "Ошибка инициализации бота", err.Error())
	}

	is_send_newsletter := false
	var media_group []models.InputMedia

	// Получение текстовых полей из формы
	message_text := c.PostForm("message")

	// Проверяем наличие файлов
	form, err := c.MultipartForm()
	if err != nil {
		Answer_BadRequest(c, ANSWER_DB_GENERAL_ERROR().Code, ANSWER_DB_GENERAL_ERROR().Message)
		return
	}

	user_list := db.DB_GET_Users()

	files := form.File["files[]"]

	if len(files) == 0 {
		if message_text == "" {
			Answer_BadRequest(c, ANSWER_EMPTY_FIELDS().Code, ANSWER_EMPTY_FIELDS().Message)
			return
		}

		if len(user_list) != 0 {
			for _, current_user := range user_list {
				if current_user.IsSubscribeNewsletter {
					params := &bot.SendMessageParams{
						ChatID: current_user.UserTgID,
						Text:   "[РАССЫЛКА]" + "\n" + message_text,
					}

					b.SendMessage(ctx, params)
					is_send_newsletter = true
				}
			}
		}
	} else {
		// Обрабатываем файлы
		for i, file := range files {
			// Открываем файл
			src, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer src.Close()

			// Читаем файл в байтовый массив
			fileData, err := io.ReadAll(src)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Добавляем файл в группу медиа
			media := &models.InputMediaPhoto{
				Media:           "attach://" + file.Filename,
				ParseMode:       models.ParseModeHTML,
				MediaAttachment: bytes.NewReader(fileData),
			}

			if i == 0 {
				media.Caption = "[РАССЫЛКА]" + "\n" + message_text
			}

			media_group = append(media_group, media)
		}

		if len(user_list) != 0 {
			for _, current_user := range user_list {
				if current_user.IsSubscribeNewsletter {
					params := &bot.SendMediaGroupParams{
						ChatID: current_user.UserTgID,
						Media:  media_group,
					}

					b.SendMediaGroup(context.Background(), params)
					is_send_newsletter = true
				}
			}
		}
	}

	if is_send_newsletter {
		Answer_OK(c)
	} else {
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
	}
}

// SendMessageUserFromSupport
func Handler_SendMessageUserFromSupport(c *gin.Context) {
	json_data := new(SendMessageUserFromSupport_Request)
	err_json_bin := c.ShouldBindJSON(&json_data)

	if err_json_bin != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_SendMessageUserFromSupport", "c.ShouldBindJSON", "Неверные данные в запросе", err_json_bin.Error())
		Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err_json_bin.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{}

	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN, opts...)
	if err != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_SendMessageUser", "gotgbot.NewBot", "Ошибка инициализации бота", err.Error())
	}

	params := &bot.SendMessageParams{
		ChatID:    json_data.UserTgID,
		ParseMode: models.ParseModeHTML,
	}

	// Получаем текущую дату и время.
	currentTime := time.Now()

	// Форматируем дату и время в "дд.мм.гггг чч:мм".
	formattedDateTime := currentTime.Format("02.01.2006 15:04")

	params.Text = formattedDateTime + "\n" +
		"<b>Вам поступил ответ от техподдержки</b>" + "\n" +
		"Ответ на обращение №: <b>" + json_data.ReferenceNumber + "</b>" + "\n" +
		"------------------------------" + "\n" + "\n" +
		json_data.Message

	_, err_send := b.SendMessage(ctx, params)
	if err_send != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_SendMessageUserFromSupport", "bot.Send", "Ошибка отправки сообщения", err_send.Error())
		Answer_BadRequest(c, ANSWER_BOT_SEND_MESSAGE_ERROR(err.Error()).Code, ANSWER_BOT_SEND_MESSAGE_ERROR(err_send.Error()).Message+" Error: "+err_send.Error())
	} else {
		Answer_OK(c)
	}
}

// Загрузка файла для календаря мероприятий
func Handler_UploadFile_CalendarActivities(c *gin.Context) {
	// Получение файла из запроса
	file, err := c.FormFile("image")
	if err != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_UploadFile_CalendarActivities", "c.FormFile", "Неверные данные в запросе", err.Error())
		Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		return
	}

	err_rem := removeAllContents("./img/calendar_activities")
	if err_rem != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_UploadFile_CalendarActivities", "removeAllContents('./img/calendar_activities')", "Ошибка при удалении каталога для календаря мероприятий", err_rem.Error())
	}

	extension := filepath.Ext(file.Filename)

	// Создание пути для сохранения файла
	uploadDir := "./img/calendar_activities" // Замените на путь к вашей папке img
	os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, "calendar_activities"+extension) // Замените на желаемое имя файла и расширение

	// Сохранение файла
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Answer_BadRequest(c, ANSWER_INVALID_FILE_UPLOAD().Code, ANSWER_INVALID_FILE_UPLOAD().Message)
		return
	}

	Answer_OK(c)
}

// Удалить всё в БД
func Handler_DeleteObjects_All(c *gin.Context) {

	db.DB_DELETE_Users()
	db.DB_DELETE_Requests()
	db.DB_DELETE_Activities()
	db.DB_DELETE_AnimeRoulettes()

	err := removeAllContents("./img/activities")
	if err != nil {
		rr_debug.PrintLOG("api_static.go", "Handler_DeleteObjects_All", "removeAllContents('./img/activities')", "Ошибка при удалении каталога для картинок мероприятий", err.Error())
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

// BroadcastResult tracks the success/failure status for each user in a broadcast
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
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
		return
	}

	if len(filtered_users) > 0 {
		// Create a slice to store results for each user
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
	} else {
		Answer_NotFound(c, ANSWER_OBJECT_NOT_FOUND().Code, ANSWER_OBJECT_NOT_FOUND().Message)
	}
}
