package routes

import (
	"rr/kotatsutgbot/db"
)

// ----------------------------------------------
//
// 				Служебные структуры
//
// ----------------------------------------------

// ----------------------------------------------
//
//	API ACTIVITIES
//
// ----------------------------------------------
// =========================================================================
//
// Handler_API_Activities_CreateObject
//
// =========================================================================
type Create_Activities struct {
	Title       string `json:"title"`
	DateMeeting string `json:"date_meeting"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// =========================================================================
//
//	Handler_API_Activities_GetList
//
// =========================================================================
type GetList_Activities_Answer struct {
	ListActivities []db.Activity_ReadJSON `json:"list_activities"`
}

// =========================================================================
//
//	Handler_API_Activities_UpdateObject
//
// =========================================================================
type UpdateObject_Activity_Data struct {
	ActivityID string `json:"activity_id"`
	Status     int    `json:"status"`
}

// ----------------------------------------------
//
//	API ANIME_ROULETTES
//
// ----------------------------------------------
// =========================================================================
//
// Handler_API_AnimeRoulettes_CreateObject
//
// =========================================================================
type Create_AnimeRoulettes struct {
	Stages []Create_AnimeRoulettes_Stages `json:"stages"`
	Theme  *string                        `json:"theme"`
}

type Create_AnimeRoulettes_Stages struct {
	Stage     int    `json:"stage"`      // Этап
	StartDate string `json:"start_date"` // Дата начала этапа рулетки
	EndDate   string `json:"end_date"`   // Дата окончания этапа рулетки
}

// =========================================================================
//
// Handler_API_AnimeRoulettes_GetList
//
// =========================================================================
type GetList_AnimeRoulettes_Answer struct {
	ListAnimeRoulettes []db.AnimeRoulette_ReadJSON `json:"list_anime_roulettes"`
}

// =========================================================================
//
//	Handler_API_AnimeRoulettes_UpdateObject
//
// =========================================================================
type UpdateObject_AnimeRoulette_Data struct {
	Status       int                          `json:"status"`
	CurrentStage int                          `json:"current_stage"`
	Theme        string                       `json:"theme"`
	StageDate    Create_AnimeRoulettes_Stages `json:"stage_date"`
	StageNewDate int                          `json:"stage_new_date"`
}

// ----------------------------------------------
//
//	API USERS
//
// ----------------------------------------------
// =========================================================================
//
//	GetList
//
// =========================================================================
// Ответ на получение списка пользователей
type GetList_Users_Answer struct {
	ListUsers []db.User_ReadJSON `json:"list_users"`
}

type SendMessageUser_Request struct {
	Message string `json:"message"`
}

// =========================================================================
//
//	UPDATE
//
// =========================================================================
type UpdateObject_User_Data struct {
	UserTgID     string `json:"user_tg_id"`
	IsClubMember int    `json:"is_club_member"`
}

type UpdateObject_User_ClubMember struct {
	UserTgID     string `json:"user_tg_id"`
	IsClubMember int    `json:"is_club_member"`
}

// ----------------------------------------------
//
//	API REQUESTS
//
// ----------------------------------------------
// =========================================================================
//
//	GetList
//
// =========================================================================
type GetList_Requests_Answer struct {
	ListRequests []db.Request_ReadJSON `json:"list_requests"`
}

// =========================================================================
//
//	UPDATE
//
// =========================================================================
type UpdateObject_Request_Data struct {
	RequestID string `json:"request_id"`
	Status    int    `json:"status"`
}

// =========================================================================
//
//	SendMessageUserFromSupport
//
// =========================================================================
type SendMessageUserFromSupport_Request struct {
	UserTgID        int64  `json:"user_tg_id"`
	ReferenceNumber string `json:"reference_number"`
	Message         string `json:"message"`
}
