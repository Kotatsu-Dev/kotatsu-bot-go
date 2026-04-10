package routes

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/middleware"
)

func RunServer() {
	if !config.GetConfig().CONFIG_IS_DEBUG {
		gin.DisableConsoleColor()
		f, _ := os.Create("gin_server.log")
		gin.DefaultWriter = io.MultiWriter(f)
	}

	r := gin.Default()

	r.NoRoute(static.ServeRoot("/", config.ByUI("./static/dist/")))

	// TODO: Index page + robots.txt
	r.GET("/admin", Handler_NewAdminPanel)
	r.GET("/login", Handler_Login)

	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	}))

	api := r.Group("/api")
	if !config.GetConfig().IGNORE_AUTH {
		api.Use(middleware.AuthMiddleware())
	}
	{
		users := api.Group("/users")
		{
			users.GET("/", Handler_API_Users_GetList)
			users.PUT("/", Handler_API_Users_UpdateObject)
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

		broadcast := api.Group("/broadcast")
		{
			broadcast.POST("/", Handler_API_SendBroadcast)
		}

		calendar := api.Group("/calendar")
		{
			// TODO: Move to static file url
			calendar.GET("/", Handler_GetCalendarActivities_Image_File)
			calendar.POST("/", Handler_UploadFile_CalendarActivities)
		}
	}

	if config.GetConfig().CONFIG_IS_DEBUG_SERVERLESS {
		r.Run(":" + config.GetConfig().CONFIG_DEBUG_SERVERLESS_SERVER_PORT)
	} else {
		r.Run(":" + config.GetConfig().CONFIG_RELEASE_SERVER_PORT)
	}
}

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
