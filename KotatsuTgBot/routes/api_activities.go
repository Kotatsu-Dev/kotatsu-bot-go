// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

//
// ----------------------------------------------------------------------------------
//
// 								Activities (Пути)
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

	//Системные пакеты
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Создать мероприятие
func Handler_API_Activities_CreateObject(c *gin.Context) {

	// Получение текстовых полей из формы
	title := c.PostForm("title")
	dateMeeting := c.PostForm("date_meeting")
	description := c.PostForm("description")
	location := c.PostForm("location")

	files := c.Request.MultipartForm.File["send_images"]

	x := 0
	x_str := ""

	var uploadDir string
	var filePath string
	var err_file error
	var images_path []string

	if len(files) != 0 {
		for i, file := range files {
			// Используем filepath.Ext для получения расширения
			extension := filepath.Ext(file.Filename)
			x = i + 1

			x_str = strconv.Itoa(x)

			// // Создание пути для сохранения файла
			uploadDir = "./img/activities/" + title + "/"
			os.MkdirAll(uploadDir, os.ModePerm)
			filePath = filepath.Join(uploadDir, x_str+extension) // Замените на желаемое имя файла и расширение

			images_path = append(images_path, uploadDir+x_str+extension)

			// // Сохранение файла
			if err_file = c.SaveUploadedFile(file, filePath); err_file != nil {
				Answer_BadRequest(c, ANSWER_INVALID_FILE_UPLOAD().Code, ANSWER_INVALID_FILE_UPLOAD().Message)
				return
			}
		}
	}

	// Проверка на необходимые поля
	if title == "" {
		rr_debug.PrintLOG("api_activities.go", "Handler_API_Activities_CreateObject", "Empty Fields", "Пустые данные в запросе", "")
		Answer_BadRequest(c, ANSWER_EMPTY_FIELDS().Code, ANSWER_EMPTY_FIELDS().Message)
		return
	} else {

		// Формат строки даты и времени
		layout := "2006-01-02 15:04"

		// Парсим строку в time.Time
		date_meeting_time, err_time := time.Parse(layout, dateMeeting)
		if err_time != nil {
			rr_debug.PrintLOG("api_activities.go", "Handler_API_Activities_CreateObject", "DateMeeting Parse", "Ошибка при парсинге времени", err_time.Error())
			return
		}

		activity_to_add := db.Activity_CreateJSON{
			Title:       title,
			DateMeeting: date_meeting_time,
			Description: description,
			Location:    location,
			PathsImages: images_path,
		}

		db_answer_code := db.DB_CREATE_Activity(&activity_to_add)

		switch db_answer_code {
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
}

// Получить все мероприятия
func Handler_API_Activities_GetList(c *gin.Context) {

	list_activities := db.DB_GET_Activities()
	answer := GetList_Activities_Answer{
		ListActivities: list_activities,
	}

	Answer_SendObject(c, answer)
	return
}

// Обновить данные мероприятия
func Handler_API_Activities_UpdateObject(c *gin.Context) {

	var update_json map[string]interface{}

	err := c.ShouldBindJSON(&update_json)
	if err != nil {
		rr_debug.PrintLOG("api_requests.go", "Handler_API_Activities_UpdateObject", "c.ShouldBindJSON", "Неверные данные в запросе", err.Error())
		if config.CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	}

	db_answer_code := db.DB_UPDATE_Activity(update_json)
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

// Удалить все мероприятия
func Handler_API_Activities_DeleteObject_ALL(c *gin.Context) {

	db_answer_code := db.DB_DELETE_Activities()
	switch db_answer_code {

	case db.DB_ANSWER_SUCCESS:
		err := removeAllContents("./img/activities")
		if err != nil {
			rr_debug.PrintLOG("api_activities.go", "Handler_API_Activities_DeleteObject_ALL", "removeAllContents('./img/activities')", "Ошибка при удалении каталога для картинок мероприятий", err.Error())
		}

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

func removeAllContents(directory string) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Пропускаем саму папку
		if path == directory {
			return nil
		}
		if info.IsDir() {
			// Удаляем подкаталоги
			return os.RemoveAll(path)
		}
		// Удаляем файлы
		return os.Remove(path)
	})

	if err != nil {
		return err
	}

	// Создаем пустую папку после удаления
	return os.MkdirAll(directory, os.ModePerm)
}
