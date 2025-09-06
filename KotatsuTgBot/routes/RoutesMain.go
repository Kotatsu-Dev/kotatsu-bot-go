// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

//
// ----------------------------------------------------------------------------------
//
// 								Routes (Пути)
//
// ----------------------------------------------------------------------------------
//

package routes

import (
	//Внутренние пакеты проекта

	"io"
	"os"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/middleware"

	// "../modules/rr_randstr"

	//Сторонние библиотеки

	//Системные пакеты
	"fmt"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// ----------------------------------------------
//
//	MAIN
//
// ----------------------------------------------
func RunServer() {

	// ЛОГ ФАЙЛ, если у нас не отладка
	if !config.GetConfig().CONFIG_IS_DEBUG {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()

		// Logging to a file.
		f, _ := os.Create("gin_server.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}

	//Создаем роутер для обработки запросов
	r := gin.Default()

	//Раздача статики для дебаг-версии
	if config.GetConfig().CONFIG_IS_DEBUG_SERVERLESS {
		r.NoRoute(static.ServeRoot("/", "./static/dist/"))
		r.Static("/assets", "./assets") //Для статики в режиме отладки
		//Загружаем HTML
		r.LoadHTMLGlob("assets/html/*")
		// r.LoadHTMLFiles("static/dist/index.html")
	} else {
		//Загружаем HTML
		r.LoadHTMLGlob("static/assets/html/*")
	}

	//CORS
	r.Use(middleware.CORSMiddleware())

	//
	// 	   --------- Пути ---------
	// 	Реализацию путей см. routes.go
	//

	//
	//Пути общие
	//

	r.GET("/", Handler_Index)
	r.GET("/admin-panel", Handler_AdminPanel)
	r.GET("/new-admin-panel", Handler_NewAdminPanel)
	r.GET("/get-calendar-file", Handler_GetCalendarActivities_Image_File)
	r.GET("/support-response", Handler_SupportResponse)
	r.GET("/login", Handler_Login)

	// Основные пути
	r.POST("/send-message-user", Handler_SendMessageUser)
	r.POST("/send-message-user-from-support", Handler_SendMessageUserFromSupport)
	r.POST("/upload-file-calendar-activities", Handler_UploadFile_CalendarActivities)
	r.DELETE("/all-db", Handler_DeleteObjects_All)
	// r.GET("/get-result", Handler_GetWorkerStatus)

	// Группа API
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	// api.Use(middleware.AuthRequired())
	{

		users := api.Group("/users")
		{
			users.GET("/", Handler_API_Users_GetList)
			users.PUT("/", Handler_API_Users_UpdateObject)
			users.PUT("/club-member", Handler_API_Users_UpdateObject_ClubMember)
			users.DELETE("/", Handler_API_Users_DeleteObject_ALL)
		}

		activities := api.Group("/activities")
		{
			activities.POST("/", Handler_API_Activities_CreateObject)
			activities.GET("/", Handler_API_Activities_GetList)
			activities.PUT("/", Handler_API_Activities_UpdateObject)
			activities.DELETE("/", Handler_API_Activities_DeleteObject_ALL)
		}

		requests := api.Group("/requests")
		{
			requests.GET("/", Handler_API_Requests_GetList)
			requests.PUT("/", Handler_API_Requests_UpdateObject)
			requests.PUT("/choice", Handler_API_Requests_UpdateObject_Choice)
			requests.DELETE("/", Handler_API_Requests_DeleteObject_ALL)
		}

		anime_roulettes := api.Group("/roulettes")
		{
			anime_roulettes.GET("/", Handler_API_AnimeRoulettes_GetList)
			anime_roulettes.GET("/active", Handler_API_AnimeRoulettes_GetActive)
			anime_roulettes.POST("/", Handler_API_AnimeRoulettes_CreateObject)
			anime_roulettes.PUT("/", Handler_API_AnimeRoulettes_UpdateObject)
			anime_roulettes.DELETE("/", Handler_API_AnimeRoulettes_DeleteObject_ALL)
		}
	}

	if config.GetConfig().CONFIG_IS_DEBUG_SERVERLESS {
		//Запуск сервера
		r.Run(":" + config.GetConfig().CONFIG_DEBUG_SERVERLESS_SERVER_PORT) // listen and serve on 0.0.0.0:PORT
	} else {
		//Запуск сервера
		r.Run(":" + config.GetConfig().CONFIG_RELEASE_SERVER_PORT) // listen and serve on 0.0.0.0:PORT
	}
}

// ----------------------------------------------
//
// 				Структуры
//
// ----------------------------------------------

// ----------------------------------------------
//
// 				Root requests
//
// ----------------------------------------------

// Вывод отладочного сообщения В КОНСОЛЬ, если у нас отладка
func LOG(message string) {
	if config.GetConfig().CONFIG_IS_DEBUG {
		fmt.Println("[DEBUG]: " + message)
	}
}
