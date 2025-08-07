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

// /Выдача фронта для ответа от техподдержки
func Handler_SupportResponse(c *gin.Context) {
	c.HTML(http.StatusOK, "support_response.html", gin.H{})
}

// Получить изображение текущего календаря мероприятий
func Handler_GetCalendarActivities_Image_File(c *gin.Context) {
	// Проверка существования файла
	if _, err := os.Stat("./img/calendar_activities/calendar_activities.png"); err == nil {
		// Answer_OK(c)
		Answer_File(c, "/img/calendar_activities/calendar_activities.png")
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
