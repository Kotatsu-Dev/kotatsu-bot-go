package cb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func check_is_participant(user *db.User_ReadJSON, activity *db.Activity_ReadJSON) bool {
	for _, participant := range activity.Participants {
		if user.UserTgID == participant.UserTgID {
			return true
		}
	}
	return false
}

func QueryData(update *models.Update) string {
	parts := strings.Split(update.CallbackQuery.Data, "::")
	data := parts[1]

	return data
}

func QueryDataUint(update *models.Update) (uint64, error) {
	data := QueryData(update)
	return strconv.ParseUint(data, 10, 64)
}

func JoinClubQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	AnswerQuery(ctx, b, update)

	switch QueryData(update) {
	case "from_ITMO_student":
		UpdateCurrentUser(current_user, map[string]any{
			"step":        config.STEP_ITMO_ENTER_ISU,
			"itmo_status": db.Student,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_graduate":
		UpdateCurrentUser(current_user, map[string]any{
			"step":        config.STEP_ITMO_ENTER_ISU,
			"itmo_status": db.Graduate,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_employee":
		UpdateCurrentUser(current_user, map[string]any{
			"step":        config.STEP_ITMO_ENTER_ISU,
			"itmo_status": db.Employee,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_student_employee":
		UpdateCurrentUser(current_user, map[string]any{
			"step":        config.STEP_ITMO_ENTER_ISU,
			"itmo_status": db.StudentEmployee,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_graduate_employee":
		UpdateCurrentUser(current_user, map[string]any{
			"step":        config.STEP_ITMO_ENTER_ISU,
			"itmo_status": db.GraduateEmployee,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	default:
		UpdateCurrentUser(current_user, map[string]any{
			"step":        config.STEP_NOITMO_ENTER_FULLNAME,
			"itmo_status": db.Guest,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_full_name", nil,
		)
	}
}

func RelevancePhoneQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	AnswerQuery(ctx, b, update)

	db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
	if db_answer_code == db.DB_ANSWER_SUCCESS {
		if QueryData(update) == "yes" {
			if activity.Status {
				if activity.GuestRegistrationUntil != nil &&
					!current_user.IsITMO &&
					activity.GuestRegistrationUntil.Before(time.Now()) {
					UpdateStep(current_user, config.STEP_DEFAULT)
					SendMessageM(
						ctx, b, update,
						"events.registration_closed",
						keyboards.ListEvents,
					)
				} else {
					db.DB_UPDATE_Activity_ADD_Participants(uint(activity.ID), current_user.ID)
					UpdateStep(current_user, config.STEP_DEFAULT)
					SendMessageMT(
						ctx, b, update,
						"events.registered", activity,
						keyboards.ListEvents,
					)
				}
			} else {
				SendMessageMT(
					ctx, b, update,
					"events.non_existent", activity,
					keyboards.ListEvents,
				)
			}

		} else {
			UpdateStep(current_user, config.STEP_CHANGING_PHONE)
			SendMessageM(
				ctx, b, update,
				"request.send_phone", keyboards.Keyboard_RequestContact,
			)
		}
	}
}

func UnsubscribeQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	AnswerQuery(ctx, b, update)

	activity_id, err := QueryDataUint(update)
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITY_SUBSCRIBE", "strconv.ParseUint", "Ошибка конвертации строки в uint", err.Error())
		return
	}

	db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(activity_id))
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		db_answer_code_remove := db.DB_UPDATE_Activity_REMOVE_Participant(uint(activity_id), current_user.ID)
		switch db_answer_code_remove {
		case db.DB_ANSWER_SUCCESS:
			SendMessageMT(
				ctx, b, update,
				"events.unregistered", activity,
				keyboards.ListEvents,
			)

		case db.DB_ANSWER_OBJECT_NOT_FOUND:
			SendMessageM(
				ctx, b, update,
				"events.non_existent",
				keyboards.ListEvents,
			)

		case db.DB_ANSWER_OBJECT_EXISTS:
			SendMessageM(
				ctx, b, update,
				"events.not_registered",
				keyboards.ListEvents,
			)
		}
	}
}

func SubscribeQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	AnswerQuery(ctx, b, update)

	activity_id, err := QueryDataUint(update)
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITY_SUBSCRIBE", "strconv.ParseUint", "Ошибка конвертации строки в uint", err.Error())
		return
	}

	if current_user.IsFilledData {
		if current_user.IsITMO {
			db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(activity_id))
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:
				if activity.Status {
					// Is ITMO = true
					db.DB_UPDATE_Activity_ADD_Participants(uint(activity_id), current_user.ID)
					SendMessageMT(
						ctx, b, update,
						"events.registered", activity,
						keyboards.ListEvents,
					)
				} else {
					SendMessageM(
						ctx, b, update,
						"events.non_existent",
						keyboards.ListEvents,
					)
				}
				UpdateStep(current_user, config.STEP_DEFAULT)
			}
		} else {
			SendMessageMT(
				ctx, b, update,
				"events.phone_number", current_user.PhoneNumber,
				keyboards.InlineKbd_RelevancePhoneNumber,
			)

			UpdateCurrentUser(current_user, map[string]any{
				"step":             config.STEP_DEFAULT,
				"temp_activity_id": int(activity_id),
			})
		}
	} else {
		SendMessageM(
			ctx, b, update,
			"request.unknown", keyboards.InlineKbd_Appointment,
		)

		UpdateCurrentUser(current_user, map[string]any{
			"temp_activity_id": int(activity_id),
		})
	}
}

func ActivitiesQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	AnswerQuery(ctx, b, update)

	activity_id, err := QueryDataUint(update)
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "strconv.ParseUint", "Ошибка конвертации строки в uint", err.Error())
		return
	}

	db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(activity_id))
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		var formattedTime, formattedDate string
		is_participant := check_is_participant(current_user, activity)

		loc, _ := time.LoadLocation("Europe/Moscow")
		formattedTime = activity.DateMeeting.In(loc).Format("15:04")
		formattedDate = FormatDate(activity.DateMeeting.In(loc))

		if len(activity.PathsImages) != 0 {
			files := make([]io.Reader, len(activity.PathsImages))
			for i, path := range activity.PathsImages {
				fileData, err := os.ReadFile(path)
				if err != nil {
					rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "os.Open(output_image_path)", "Ошибка открытия файла", err.Error())
					return
				}

				files[i] = bytes.NewReader(fileData)
			}
			SendPhotosM(
				ctx, b, update,
				files, "",
			)
		}

		SendMessageMT(
			ctx, b, update,
			"events.format", &map[string]any{
				"activity":      activity,
				"formattedDate": formattedDate,
				"formattedTime": formattedTime,
			},
			ITE(is_participant,
				keyboards.CreateInlineKbd_UnsubscribeActivity(int(activity.ID)),
				keyboards.CreateInlineKbd_SubscribeActivity(int(activity.ID)),
			),
		)

		UpdateStep(current_user, config.STEP_ACTIVITY)
	}
}

func AppointQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	AnswerQuery(ctx, b, update)

	switch QueryData(update) {
	case "from_ITMO_student":
		UpdateCurrentUser(current_user, map[string]any{
			"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
			"itmo_status": db.Student,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_graduate":
		UpdateCurrentUser(current_user, map[string]any{
			"user_tg_id":  update.CallbackQuery.From.ID,
			"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
			"itmo_status": db.Graduate,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_employee":
		UpdateCurrentUser(current_user, map[string]any{
			"user_tg_id":  update.CallbackQuery.From.ID,
			"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
			"itmo_status": db.Employee,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_student_employee":
		UpdateCurrentUser(current_user, map[string]any{
			"user_tg_id":  update.CallbackQuery.From.ID,
			"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
			"itmo_status": db.StudentEmployee,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	case "from_ITMO_graduate_employee":
		UpdateCurrentUser(current_user, map[string]any{
			"user_tg_id":  update.CallbackQuery.From.ID,
			"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
			"itmo_status": db.GraduateEmployee,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_isu_number", nil,
		)
	default:
		UpdateCurrentUser(current_user, map[string]any{
			"user_tg_id":  update.CallbackQuery.From.ID,
			"step":        config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME,
			"itmo_status": db.Guest,
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_full_name", nil,
		)
	}
}

func FormatDate(t time.Time) string {
	var weekday, month string
	switch t.Weekday() {
	case time.Monday:
		weekday = "понедельник"
	case time.Tuesday:
		weekday = "вторник"
	case time.Wednesday:
		weekday = "среда"
	case time.Thursday:
		weekday = "четверг"
	case time.Friday:
		weekday = "пятница"
	case time.Saturday:
		weekday = "суббота"
	case time.Sunday:
		weekday = "воскресенье"
	}

	switch t.Month() {
	case time.January:
		month = "января"
	case time.February:
		month = "февраля"
	case time.April:
		month = "апреля"
	case time.March:
		month = "марта"
	case time.May:
		month = "мая"
	case time.June:
		month = "июня"
	case time.July:
		month = "июля"
	case time.August:
		month = "августа"
	case time.September:
		month = "сентября"
	case time.October:
		month = "октября"
	case time.November:
		month = "ноября"
	case time.December:
		month = "декабря"
	}

	return fmt.Sprintf("%d %s (%s)", t.Day(), month, weekday)
}
