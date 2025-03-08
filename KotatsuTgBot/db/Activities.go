package db

import (

	//Внутренние пакеты проекта
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"gorm.io/gorm"

	//Системные пакеты
	"time"
)

type Activity struct {
	gorm.Model
	Title        string    `json:"title"`                                          // Название мероприятия
	Participants []*User   `json:"participants" gorm:"many2many:user_activities;"` // Участники мероприятия
	DateMeeting  time.Time `json:"date_meeting"`                                   // Дата проведения мероприятия
	Description  string    `json:"description"`                                    // Описание мероприятия
	Location     string    `json:"location"`                                       // Место проведения мероприятия
	PathsImages  []string  `json:"paths_images" gorm:"type:text[]"`                                   // Пути к картинкам мероприятия
	Status       bool      `json:"status"`                                         // Статус мероприятия
}

type Activity_CreateJSON struct {
	Title       string    `json:"title"`
	DateMeeting time.Time `json:"date_meeting"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	PathsImages []string  `json:"paths_images"`
}

type Activity_ReadJSON struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Title        string    `json:"title"`
	Participants []*User   `json:"participants"`
	DateMeeting  time.Time `json:"date_meeting"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	PathsImages  []string  `json:"paths_images"`
	Status       bool      `json:"status"`
}

// Добавить мероприятие
func DB_CREATE_Activity(activity_to_add *Activity_CreateJSON) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var activity Activity
	db.Where("title = ?", activity_to_add.Title).First(&activity)
	if activity.ID != 0 {
		return DB_ANSWER_OBJECT_EXISTS
	}

	activity = Activity{
		Title:       activity_to_add.Title,
		DateMeeting: activity_to_add.DateMeeting,
		Description: activity_to_add.Description,
		Location:    activity_to_add.Location,
		PathsImages: activity_to_add.PathsImages,
		Status:      true,
	}

	db.Save(&activity)
	return DB_ANSWER_SUCCESS
}

// Получить мероприятие по названию
func DB_GET_Activity_BY_Title(title string) (int, *Activity_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	activity := new(Activity)
	db.Preload("Participants").Where("title = ?", title).First(&activity)
	if activity.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	activity_read := Activity_ReadJSON{
		ID:           activity.ID,
		CreatedAt:    activity.CreatedAt,
		Title:        activity.Title,
		Participants: activity.Participants,
		DateMeeting:  activity.DateMeeting,
		Description:  activity.Description,
		Location:     activity.Location,
		PathsImages:  activity.PathsImages,
		Status:       activity.Status,
	}

	return DB_ANSWER_SUCCESS, &activity_read
}

// Получить мероприятие по ID
func DB_GET_Activity_BY_ID(activity_id uint) (int, *Activity_ReadJSON) {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	activity := new(Activity)
	db.Preload("Participants").First(&activity, activity_id)
	if activity.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND, nil
	}

	activity_read := Activity_ReadJSON{
		ID:           activity.ID,
		CreatedAt:    activity.CreatedAt,
		Title:        activity.Title,
		Participants: activity.Participants,
		DateMeeting:  activity.DateMeeting,
		Description:  activity.Description,
		Location:     activity.Location,
		PathsImages:  activity.PathsImages,
		Status:       activity.Status,
	}

	return DB_ANSWER_SUCCESS, &activity_read
}

// Получить список всех мероприятий
func DB_GET_Activities() []Activity_ReadJSON {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var activities []Activity

	// Загружаем связанные сущности InvitedUsers
	db.Preload("Participants").Find(&activities)

	activities_list := make([]Activity_ReadJSON, 0)
	if len(activities) <= 0 {
		return activities_list
	}

	for _, activity := range activities {

		current_activity := Activity_ReadJSON{
			ID:           activity.ID,
			CreatedAt:    activity.CreatedAt,
			Title:        activity.Title,
			Participants: activity.Participants,
			DateMeeting:  activity.DateMeeting,
			Description:  activity.Description,
			Location:     activity.Location,
			PathsImages:  activity.PathsImages,
			Status:       activity.Status,
		}
		activities_list = append(activities_list, current_activity)
	}

	return activities_list
}

// Обновить данные мероприятия
func DB_UPDATE_Activity(update_json map[string]interface{}) int {
	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var activity Activity
	activity_id, ok := update_json["activity_id"].(uint)
	if ok {
		db.Preload("Participants").First(&activity, activity_id)
		if activity.ID == 0 {
			return DB_ANSWER_OBJECT_NOT_FOUND
		}
	} else {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	// Обновляем поля, если они присутствуют в карте
	for key, value := range update_json {
		switch key {
		case "status":
			if v, ok := value.(bool); ok && v != activity.Status {
				activity.Status = v
			}
		}
	}

	db.Save(&activity)
	return DB_ANSWER_SUCCESS
}

// Добавляем пользователей в мероприятие
func DB_UPDATE_Activity_ADD_Participants(activity_id uint, user_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var activity Activity

	db.First(&activity, activity_id)
	if activity.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	var user User
	db.First(&user, user_id)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	// Добавить пользователя в массив Participants
	activity.Participants = append(activity.Participants, &user)

	db.Save(&activity)
	return DB_ANSWER_SUCCESS
}

// Удаляем пользователя из мероприятия
func DB_UPDATE_Activity_REMOVE_Participant(activity_id uint, user_id uint) int {
	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var activity Activity

	db.First(&activity, activity_id)
	if activity.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	var user User
	db.First(&user, user_id)
	if user.ID == 0 {
		return DB_ANSWER_OBJECT_NOT_FOUND
	}

	// Находим индекс пользователя в массиве Participants
	userIndex := -1
	for i, participant := range activity.Participants {
		if participant.ID == user.ID {
			userIndex = i
			break
		}
	}

	// Если пользователь не найден в массиве Participants
	if userIndex == -1 {
		return DB_ANSWER_OBJECT_EXISTS
	}

	// Удаляем пользователя из массива Participants
	activity.Participants = append(activity.Participants[:userIndex], activity.Participants[userIndex+1:]...)

	// Сохраняем изменения в базе данных
	db.Save(&activity)

	return DB_ANSWER_SUCCESS
}

// Удаление мероприятие по id
func DB_DELETE_Activity_BY_ID(activity_id uint) int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Unscoped().Delete(&Activity{}, activity_id)
	if result.Error != nil {
		rr_debug.PrintLOG("Activities.go", "DB_DELETE_Activity_BY_ID", "Error deleting anime_roulette:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}

// Удалить все мероприятия
func DB_DELETE_Activities() int {

	db := DB_Database()

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := db.Exec("DELETE FROM activities")
	if result.Error != nil {
		rr_debug.PrintLOG("Activities.go", "DB_DELETE_Activities", "Error deleting activity:", "Ошибка удаления", result.Error.Error())
		return DB_ANSWER_DELETE_ERROR
	} else {
		return DB_ANSWER_SUCCESS
	}
}
