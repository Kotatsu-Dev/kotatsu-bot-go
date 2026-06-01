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
	"fmt"
	"mime/multipart"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	guestRegistrationUntil := c.PostForm("guest_registration_until")
	description := c.PostForm("description")
	location := c.PostForm("location")

	files := c.Request.MultipartForm.File["send_images"]
	if len(files) <= 0 {
		files = c.Request.MultipartForm.File["send_images[]"]
	}

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

			uploadDir = config.ByUI(filepath.Join(uploadDir, title))
			fileName := x_str + "." + uuid.NewString() + extension
			os.MkdirAll(uploadDir, os.ModePerm)
			filePath = filepath.Join(uploadDir, fileName)

			images_path = append(images_path, filePath)

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
		date_meeting_time, err_time := time.Parse(time.RFC3339, dateMeeting)
		if err_time != nil {
			rr_debug.PrintLOG("api_activities.go", "Handler_API_Activities_CreateObject", "DateMeeting Parse", "Ошибка при парсинге времени", err_time.Error())
			return
		}

		var guest_registration_until *time.Time = nil
		if guestRegistrationUntil != "" {
			_guest_registration_until, err_time := time.Parse(time.RFC3339, guestRegistrationUntil)
			if err_time != nil {
				rr_debug.PrintLOG("api_activities.go", "Handler_API_Activities_CreateObject", "GuestRegistrationUntil Parse", "Ошибка при парсинге времени", err_time.Error())
				Answer_BadRequest(c, ANSWER_EMPTY_FIELDS().Code, ANSWER_EMPTY_FIELDS().Message)
				return
			}
			guest_registration_until = &_guest_registration_until
		}

		activity_to_add := db.Activity_CreateJSON{
			Title:                  title,
			DateMeeting:            date_meeting_time,
			GuestRegistrationUntil: guest_registration_until,
			Description:            description,
			Location:               location,
			PathsImages:            images_path,
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
}

// Обновить данные мероприятия
func Handler_API_Activities_UpdateObject(c *gin.Context) {

	update_json := make(map[string]interface{})

	if err := c.Request.ParseMultipartForm(2 << 32); err != nil {
		if config.GetConfig().CONFIG_IS_DEBUG {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message+" Error: "+err.Error())
		} else {
			Answer_BadRequest(c, ANSWER_INVALID_JSON().Code, ANSWER_INVALID_JSON().Message)
		}
		return
	}

	for key, values := range c.Request.MultipartForm.Value {
		if len(values) > 0 {
			switch key {
			case "activity_id":
				update_json[key], _ = strconv.ParseFloat(values[0], 64)
			case "title":
				update_json[key] = values[0]
			case "date_meeting":
				update_json[key] = values[0]
			case "guest_registration_until":
				update_json[key] = values[0]
			case "description":
				update_json[key] = values[0]
			case "location":
				update_json[key] = values[0]
			case "status":
				update_json[key], _ = strconv.ParseBool(values[0])
			}
		}
	}

	var files []*multipart.FileHeader
	if c.Request.MultipartForm != nil {
		files = c.Request.MultipartForm.File["send_images"]
		if len(files) <= 0 {
			files = c.Request.MultipartForm.File["send_images[]"]
		}
	}

	var uploadDir string
	var filePath string
	var err_file error
	var images_path []string

	if len(files) != 0 {
		for _, file := range files {
			extension := filepath.Ext(file.Filename)
			name := strings.TrimSuffix(filepath.Base(file.Filename), extension)

			// TODO: Fix subfolder
			uploadDir = config.ByUI(filepath.Join(uploadDir, "new_files"))
			fileName := name + "." + uuid.NewString() + extension
			os.MkdirAll(uploadDir, os.ModePerm)
			filePath = filepath.Join(uploadDir, fileName)

			images_path = append(images_path, filePath)

			if err_file = c.SaveUploadedFile(file, filePath); err_file != nil {
				Answer_BadRequest(c, ANSWER_INVALID_FILE_UPLOAD().Code, ANSWER_INVALID_FILE_UPLOAD().Message)
				return
			}
		}

		update_json["path_images"] = images_path
	}

	fmt.Println(c.Request.MultipartForm)

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
