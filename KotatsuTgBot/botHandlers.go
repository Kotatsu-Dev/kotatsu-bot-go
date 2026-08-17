// ------------------------------------
// RR IT 2024
//
// ------------------------------------

//
// ----------------------------------------------------------------------------------
//
// 								Обработчики сообщений боту
//
// ----------------------------------------------------------------------------------
//

package main

import (
	"rr/kotatsutgbot/cb"
	"rr/kotatsutgbot/cb/roulette"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"
	"time"

	//Сторонние библиотеки
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	//Системные пакеты
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Удалить элемент массива
func RemoveIndex(s []int64, index int) []int64 {
	return append(s[:index], s[index+1:]...)
}

//
// Главные процессы
//

func BotHandler_Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil {
		// Обработка случая, когда update пуст
		return
	} else {
		if update.Message == nil {
			if update.CallbackQuery == nil {
				return
			} else {

				db_answer_code, user := db.DB_GET_User_BY_UserTgID(update.CallbackQuery.From.ID)
				switch db_answer_code {
				case db.DB_ANSWER_SUCCESS:
					BotHandler_CallbackQuery(ctx, b, update, user)

				case db.DB_ANSWER_OBJECT_NOT_FOUND:
					proccessRegistrationCallback(ctx, b, update)
				}

			}
		} else {
			if update.Message.From == nil {
				// Обработка случая, когда Chat или From равны nil
				return
			} else {
				if update.Message.Chat.ID == config.GetConfig().CONFIG_ID_CHAT_SUPPORT {
					return
				} else {
					db_answer_code, user := db.DB_GET_User_BY_UserTgID(update.Message.From.ID)
					switch db_answer_code {
					case db.DB_ANSWER_SUCCESS:
						switch update.Message.Text {
						case config.T("keyboard.gender_male"):
							cb.SetGender(ctx, b, update, user, "male")
						case config.T("keyboard.gender_female"):
							cb.SetGender(ctx, b, update, user, "female")

						case config.T("keyboard.visited_enough"):
							cb.WasAtEvents(ctx, b, update, user, true)
						case config.T("keyboard.not_visited_enough"):
							cb.WasAtEvents(ctx, b, update, user, false)
						case config.T("keyboard.fill_back_later"):
							cb.WasntAtEvents(ctx, b, update, user, false)
						case config.T("keyboard.fill_now"):
							cb.WasntAtEvents(ctx, b, update, user, true)
						case config.T("keyboard.join_club"):
							cb.JoinClub(ctx, b, update, user)

						case config.T("keyboard.event_registration"):
							cb.SigningUpForActivity(ctx, b, update)

						case config.T("keyboard.to_main_menu"):
							cb.BackMainMenu(ctx, b, update, user)

						case config.T("keyboard.leave_club"):
							cb.LeaveClub(ctx, b, update, user)

						case config.T("keyboard.to_roulette_menu"):
							roulette.Main(ctx, b, update, user)

						case config.T("keyboard.participate_roulette"):
							roulette.Participate(ctx, b, update, user)

						case config.T("keyboard.leave_roulette"):
							roulette.CancelParticipate(ctx, b, update, user)

						case config.T("keyboard.send_title"):
							roulette.AnimeWish(ctx, b, update, user)

						case config.T("keyboard.roulette_rules"):
							roulette.Rules(ctx, b, update)

						case config.T("keyboard.roulette_theme"):
							roulette.MainTheme(ctx, b, update)

						case config.T("keyboard.roulette_list"):
							roulette.LinkMyList(ctx, b, update, user)

						case config.T("keyboard.my_events"):
							cb.MyActivities(ctx, b, update, user)

						case config.T("keyboard.to_main_menu"):
							cb.BackMainMenu(ctx, b, update, user)

						case config.T("keyboard.not_my_number"):
							cb.NoPhoneNumber(ctx, b, update, user)

						default:

							switch user.Step {
							case config.STEP_ITMO_ENTER_ISU:
								cb.ITMO_EnterISU(ctx, b, update, user, "join_club")

							case config.STEP_APPOINTMENT_ITMO_ENTER_ISU:
								cb.ITMO_EnterISU(ctx, b, update, user, "activity")

							case config.STEP_ITMO_ENTER_FULLNAME:
								cb.ITMO_EnterFullName(ctx, b, update, user, "join_club")

							case config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME:
								cb.ITMO_EnterFullName(ctx, b, update, user, "activity")

							case config.STEP_NOITMO_ENTER_FULLNAME:
								cb.NoITMO_EnterFullName(ctx, b, update, user, "join_club")

							case config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME:
								cb.NoITMO_EnterFullName(ctx, b, update, user, "activity")

							case config.STEP_NOITMO_ENTER_PHONE:
								cb.NoITMO_EnterPhoneNumber(ctx, b, update, user, "join_club")

							case config.STEP_CHANGING_PHONE:
								cb.ChangePhoneNumber(ctx, b, update, user)

							case config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE:
								cb.NoITMO_EnterPhoneNumber(ctx, b, update, user, "activity")

							case config.STEP_USER_LEAVES_CLUB:
								cb.LeavesClub(ctx, b, update, user)

							case config.STEP_ANIME_RULETTE_ENTER_ENIGMATIC_TITLE:
								roulette.EnterEnigmaticTitle(ctx, b, update, user)

							case config.STEP_ANIME_RULETTE_ENTER_LINK_MY_ANIME_LIST:
								roulette.EnterLinkMyAnimeList(ctx, b, update, user)

							default:
								proccessText_Unknown(ctx, b, update)
							}

						}
					case db.DB_ANSWER_OBJECT_NOT_FOUND:
						proccessRegistrationMessage(ctx, b, update)
					}
				}
			}
		}
	}
}

// Процесс регистрации из сообщения
func proccessRegistrationMessage(ctx context.Context, b *bot.Bot, update *models.Update) {

	params := &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
	}

	if update.Message.Text == config.T("keyboard.continue") {
		full_tg_name := update.Message.From.FirstName + " " + update.Message.From.LastName
		db_answer_reg := regUser(update.Message.From.ID, full_tg_name, update.Message.From.Username)

		switch db_answer_reg {
		case db.DB_ANSWER_SUCCESS:
			params.Text = config.T("gender_select")
			params.ReplyMarkup = keyboards.Keyboard_GenderSelect

		case db.DB_ANSWER_OBJECT_EXISTS:
			params.Text = config.T("registered")

			_, old_user := db.DB_GET_User_BY_UserTgID(update.Message.From.ID)

			if old_user.IsClubMember {
				params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsClubMember
			} else {
				params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsDefault
			}

		default:
			params.Text = config.T("error.database")
			rr_debug.PrintLOG("main.go", "update.Message.Text", "activity_GetObjects()", "Ошибка работы с БД", "")
		}
	} else {
		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:   update.Message.Chat.ID,
			Document: &models.InputFileString{Data: "CAACAgIAAx0CbgUG4QACCWpostfAVRPNDHNAWu8vcIbjv0nuagACrXQAAl8iQUmAFQIjshq4bTYE"},
		})
		params.Text = config.T("personal_data")
		params.ReplyMarkup = keyboards.Registration
	}

	_, err := b.SendMessage(ctx, params)
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessRegistration", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

// Процесс регистрации из кулбека
func proccessRegistrationCallback(ctx context.Context, b *bot.Bot, update *models.Update) {

	params := &bot.SendMessageParams{
		ChatID:    update.CallbackQuery.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	if update.Message.Text == config.T("keyboard.continue") {
		full_tg_name := update.CallbackQuery.From.FirstName + " " + update.CallbackQuery.From.LastName
		db_answer_reg := regUser(update.CallbackQuery.From.ID, full_tg_name, update.CallbackQuery.From.Username)

		switch db_answer_reg {
		case db.DB_ANSWER_SUCCESS:
			params.Text = config.T("main_menu")
			params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsDefault

		case db.DB_ANSWER_OBJECT_EXISTS:
			params.Text = config.T("registered")

			_, old_user := db.DB_GET_User_BY_UserTgID(update.Message.From.ID)

			if old_user.IsClubMember {
				params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsClubMember
			} else {
				params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsDefault
			}

		default:
			params.Text = config.T("error.database")
			rr_debug.PrintLOG("main.go", "update.Message.Text", "activity_GetObjects()", "Ошибка работы с БД", "")
		}
	} else {
		params.Text = config.T("hello") + "\n" + config.T("personal_data")
		params.ReplyMarkup = keyboards.Registration
	}

	_, err := b.SendMessage(ctx, params)
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessRegistration", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

//
//	Команды
//

// Главное меню
func BotHandler_Command_Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	db_answer_code, user := db.DB_GET_User_BY_UserTgID(update.Message.From.ID)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		params := &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			ParseMode: models.ParseModeHTML,
		}

		params.Text = config.TT("welcome", update.Message.From)

		if user.IsClubMember {
			params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsClubMember
		} else {
			params.ReplyMarkup = keyboards.Keyboard_MainMenuButtonsDefault
		}

		_, err := b.SendMessage(ctx, params)
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Start", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
		}
	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		proccessRegistrationMessage(ctx, b, update)
	}
}

func BotHandler_Command_Start_Link(ctx context.Context, b *bot.Bot, update *models.Update) {
	db_answer_code, _ := db.DB_GET_User_BY_UserTgID(update.Message.From.ID)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		params := &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			ParseMode: models.ParseModeHTML,
		}
		params_photos := &bot.SendMediaGroupParams{
			ChatID: update.Message.Chat.ID,
		}

		data := strings.TrimPrefix(update.Message.Text, "/start ")
		activity_id, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_Command_Start_Link", "strconv.ParseUint", "Ошибка конвертации строки в uint", err.Error())
			return
		}

		var media_group []models.InputMedia

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(activity_id))
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			var formattedTime, formattedDate string
			is_participant := false

			for _, participant := range activity.Participants {
				if participant.UserTgID == update.CallbackQuery.From.ID {
					is_participant = true
					break
				}
			}

			// Определите желаемый формат дд.мм чч:мм
			loc, _ := time.LoadLocation("Europe/Moscow")

			// Используйте метод Format для форматирования времени
			formattedTime = activity.DateMeeting.In(loc).Format("15:04")
			formattedDate = formatDate(activity.DateMeeting.In(loc))

			if len(activity.PathsImages) != 0 {
				for _, output_image_path := range activity.PathsImages {
					// Открываем файл
					file, err := os.Open(output_image_path)
					if err != nil {
						rr_debug.PrintLOG("botHandlers.go", "BotHandler_Command_Start_Link", "os.Open(output_image_path)", "Ошибка открытия файла", err.Error())
						return
					}
					defer file.Close()

					// Читаем файл в байтовый массив
					fileData, err := io.ReadAll(file)
					if err != nil {
						rr_debug.PrintLOG("botHandlers.go", "BotHandler_Command_Start_Link", "io.ReadAll", "Ошибка перевода файла в массив байт", err.Error())
						return
					}

					// Добавляем файл в группу медиа
					media := &models.InputMediaPhoto{
						Media:           "attach://" + filepath.Base(output_image_path),
						ParseMode:       models.ParseModeHTML,
						MediaAttachment: bytes.NewReader(fileData),
					}

					media_group = append(media_group, media)
				}

				params_photos.Media = media_group

				params.Text = config.TT("events.format", &map[string]any{
					"activity":      activity,
					"formattedDate": formattedDate,
					"formattedTime": formattedTime,
				})

				if is_participant {
					params.ReplyMarkup = keyboards.CreateInlineKbd_UnsubscribeActivity(int(activity.ID))
				} else {
					params.ReplyMarkup = keyboards.CreateInlineKbd_SubscribeActivity(int(activity.ID))
				}

				_, err_media := b.SendMediaGroup(ctx, params_photos)
				if err_media != nil {
					rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "b.SendMessage", "Ошибка отправки сообщения", err_media.Error())
				}

				_, err_msg := b.SendMessage(ctx, params)
				if err_msg != nil {
					rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
				}
			} else {

				params.Text = config.TT("events.format", &map[string]any{
					"activity":      activity,
					"formattedDate": formattedDate,
					"formattedTime": formattedTime,
				})

				if is_participant {
					params.ReplyMarkup = keyboards.CreateInlineKbd_UnsubscribeActivity(int(activity.ID))
				} else {
					params.ReplyMarkup = keyboards.CreateInlineKbd_SubscribeActivity(int(activity.ID))
				}

				_, err_msg := b.SendMessage(ctx, params)
				if err_msg != nil {
					rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
				}
			}
			db.DB_UPDATE_User(map[string]interface{}{
				"user_tg_id": update.Message.From.ID,
				"step":       config.STEP_ACTIVITY,
			})
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		proccessRegistrationMessage(ctx, b, update)
	}
}

// Неизвестное сообщение или шаг
func proccessText_Unknown(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !(update.Message != nil && update.Message.Chat.Type == models.ChatTypePrivate) {
		return
	}

	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = config.T("unknown_command")

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_Unknown", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

//
// Inline - клавиатура
//

func formatDate(t time.Time) string {
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

// Вступление в клуб - клавиши "из ИТМО", "не из ИТМО"
func BotHandler_CallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {

	var (
		parts []string
		data  string
	)

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.CallbackQuery.From.ID

	switch {

	// Вступление в клуб
	case strings.HasPrefix(update.CallbackQuery.Data, "JOIN_CLUB"):

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		params := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		switch data {
		case "from_ITMO_student":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_ITMO_ENTER_ISU,
				"itmo_status": db.Student,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_graduate":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_ITMO_ENTER_ISU,
				"itmo_status": db.Graduate,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_employee":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_ITMO_ENTER_ISU,
				"itmo_status": db.Employee,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_student_employee":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_ITMO_ENTER_ISU,
				"itmo_status": db.StudentEmployee,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_graduate_employee":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_ITMO_ENTER_ISU,
				"itmo_status": db.GraduateEmployee,
			})

			params.Text = config.T("request.enter_isu_number")
		default:
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_NOITMO_ENTER_FULLNAME,
				"itmo_status": db.Guest,
			})

			params.Text = config.T("request.enter_full_name")
		}

		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_JOIN_CLUB", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}

		// Запись на мероприятие (для не участников клуба)
	case strings.HasPrefix(update.CallbackQuery.Data, "APPOINTMENT"):
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		params := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		switch data {
		case "from_ITMO_student":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
				"itmo_status": db.Student,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_graduate":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
				"itmo_status": db.Graduate,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_employee":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
				"itmo_status": db.Employee,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_student_employee":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
				"itmo_status": db.StudentEmployee,
			})

			params.Text = config.T("request.enter_isu_number")
		case "from_ITMO_graduate_employee":
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_APPOINTMENT_ITMO_ENTER_ISU,
				"itmo_status": db.GraduateEmployee,
			})

			params.Text = config.T("request.enter_isu_number")
		default:
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":  update.CallbackQuery.From.ID,
				"step":        config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME,
				"itmo_status": db.Guest,
			})

			params.Text = config.T("request.enter_full_name")
		}

		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_APPOINTMENT", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}

	// Список мероприятий
	case strings.HasPrefix(update.CallbackQuery.Data, "ACTIVITIES"), strings.HasPrefix(update.CallbackQuery.Data, "MY_ACTIVITIES"):

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		var media_group []models.InputMedia

		params := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
			LinkPreviewOptions: &models.LinkPreviewOptions{
				IsDisabled: func(b bool) *bool { return &b }(true),
			},
		}

		params_photos := &bot.SendMediaGroupParams{
			ChatID: update.CallbackQuery.From.ID,
		}

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		// Преобразуем строку в uint
		activity_id, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "strconv.ParseUint", "Ошибка конвертации строки в uint", err.Error())
			return
		}

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(activity_id))
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			var formattedTime, formattedDate string
			is_participant := false

			for _, participant := range activity.Participants {
				if participant.UserTgID == update.CallbackQuery.From.ID {
					is_participant = true
					break
				}
			}

			// Определите желаемый формат дд.мм чч:мм
			loc, _ := time.LoadLocation("Europe/Moscow")

			// Используйте метод Format для форматирования времени
			formattedTime = activity.DateMeeting.In(loc).Format("15:04")
			formattedDate = formatDate(activity.DateMeeting.In(loc))

			if len(activity.PathsImages) != 0 {
				for _, output_image_path := range activity.PathsImages {
					// Открываем файл
					file, err := os.Open(output_image_path)
					if err != nil {
						rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "os.Open(output_image_path)", "Ошибка открытия файла", err.Error())
						return
					}
					defer file.Close()

					// Читаем файл в байтовый массив
					fileData, err := io.ReadAll(file)
					if err != nil {
						rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "io.ReadAll", "Ошибка перевода файла в массив байт", err.Error())
						return
					}

					// Добавляем файл в группу медиа
					media := &models.InputMediaPhoto{
						Media:           "attach://" + filepath.Base(output_image_path),
						ParseMode:       models.ParseModeHTML,
						MediaAttachment: bytes.NewReader(fileData),
					}

					media_group = append(media_group, media)
				}

				params_photos.Media = media_group

				params.Text = config.TT("events.format", &map[string]any{
					"activity":      activity,
					"formattedDate": formattedDate,
					"formattedTime": formattedTime,
				})

				if is_participant {
					params.ReplyMarkup = keyboards.CreateInlineKbd_UnsubscribeActivity(int(activity.ID))
				} else {
					params.ReplyMarkup = keyboards.CreateInlineKbd_SubscribeActivity(int(activity.ID))
				}

				_, err_media := b.SendMediaGroup(ctx, params_photos)
				if err_media != nil {
					rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "b.SendMessage", "Ошибка отправки сообщения", err_media.Error())
				}

				_, err_msg := b.SendMessage(ctx, params)
				if err_msg != nil {
					rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
				}
			} else {

				params.Text = config.TT("events.format", &map[string]any{
					"activity":      activity,
					"formattedDate": formattedDate,
					"formattedTime": formattedTime,
				})

				if is_participant {
					params.ReplyMarkup = keyboards.CreateInlineKbd_UnsubscribeActivity(int(activity.ID))
				} else {
					params.ReplyMarkup = keyboards.CreateInlineKbd_SubscribeActivity(int(activity.ID))
				}

				_, err_msg := b.SendMessage(ctx, params)
				if err_msg != nil {
					rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITIES", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
				}
			}
			db.DB_UPDATE_User(map[string]interface{}{
				"user_tg_id": update.CallbackQuery.From.ID,
				"step":       config.STEP_ACTIVITY,
			})
		}

	// Подписаться на мероприятие
	case strings.HasPrefix(update.CallbackQuery.Data, "ACTIVITY_SUBSCRIBE"):
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		params := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		params_load := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		// Преобразуем строку в uint
		activity_id, err := strconv.ParseUint(data, 10, 64)
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
						params.Text = config.TT("events.registered", activity)
						params.ReplyMarkup = keyboards.ListEvents
					} else {
						params.Text = config.TT("events.non_existent", activity)
						params.ReplyMarkup = keyboards.ListEvents
					}
					db.DB_UPDATE_User(map[string]interface{}{
						"user_tg_id": current_user.UserTgID,
						"step":       config.STEP_DEFAULT,
					})
				}
			} else {
				params.Text = config.TT("events.phone_number", current_user.PhoneNumber)
				params_load.ReplyMarkup = keyboards.Keyboard_ToMainMenu
				params.ReplyMarkup = keyboards.InlineKbd_RelevancePhoneNumber

				fmt.Println(activity_id)

				db.DB_UPDATE_User(map[string]interface{}{
					"user_tg_id":       current_user.UserTgID,
					"step":             config.STEP_DEFAULT,
					"temp_activity_id": int(activity_id),
				})
			}
		} else {
			params.Text = config.T("request.unknown")

			params_load.ReplyMarkup = keyboards.Keyboard_ToMainMenu
			params.ReplyMarkup = keyboards.InlineKbd_Appointment

			db.DB_UPDATE_User(map[string]interface{}{
				"user_tg_id":       current_user.UserTgID,
				"temp_activity_id": int(activity_id),
			})
		}

		_, err_msg_load := b.SendMessage(ctx, params_load)
		if err_msg_load != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITY_SUBSCRIBE", "b.SendMessage(ctx, params_load)", "Ошибка отправки сообщения", err_msg_load.Error())
		}

		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITY_SUBSCRIBE", "b.SendMessage(ctx, params)", "Ошибка отправки сообщения", err_msg.Error())
		}

	// Отписаться от мероприятия
	case strings.HasPrefix(update.CallbackQuery.Data, "ACTIVITY_UNSUBSCRIBE"):
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		params := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		// Преобразуем строку в uint
		activity_id, err := strconv.ParseUint(data, 10, 64)
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
				params.Text = config.TT("events.unregistered", activity)
				params.ReplyMarkup = keyboards.ListEvents

			case db.DB_ANSWER_OBJECT_NOT_FOUND:
				params.Text = config.T("events.non_existent")
				params.ReplyMarkup = keyboards.ListEvents

			case db.DB_ANSWER_OBJECT_EXISTS:
				params.Text = config.T("events.not_registered")
				params.ReplyMarkup = keyboards.ListEvents

			}
		}

		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_ACTIVITY_UNSUBSCRIBE", "b.SendMessage(ctx, params)", "Ошибка отправки сообщения", err_msg.Error())
		}

	// Проверка актуальности номера телефона пользователя
	case strings.HasPrefix(update.CallbackQuery.Data, "RELEVANC_PHONE"):
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		params := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		update_user_data := make(map[string]interface{})
		update_user_data["user_tg_id"] = update.CallbackQuery.From.ID

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
		if db_answer_code == db.DB_ANSWER_SUCCESS {
			if data == "yes" {
				if activity.Status {
					if activity.GuestRegistrationUntil != nil &&
						!current_user.IsITMO &&
						activity.GuestRegistrationUntil.Before(time.Now()) {
						params.Text = config.T("events.registration_closed")
						params.ReplyMarkup = keyboards.ListEvents
						update_user_data["step"] = config.STEP_DEFAULT
					} else {
						db.DB_UPDATE_Activity_ADD_Participants(uint(activity.ID), current_user.ID)
						params.Text = config.TT("events.registered", activity)
						params.ReplyMarkup = keyboards.ListEvents
						update_user_data["step"] = config.STEP_DEFAULT
					}
				} else {
					params.Text = config.TT("events.non_existent", activity)
					params.ReplyMarkup = keyboards.ListEvents
				}

			} else {
				update_user_data["step"] = config.STEP_CHANGING_PHONE
				db.DB_UPDATE_User(update_user_data)

				params.Text = config.T("request.send_phone")
				params.ReplyMarkup = keyboards.Keyboard_RequestContact
			}

			_, err_msg := b.SendMessage(ctx, params)
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_RELEVANC_PHONE", "b.SendMessage(ctx, params)", "Ошибка отправки сообщения", err_msg.Error())
			}

		}

	case strings.HasPrefix(update.CallbackQuery.Data, "ROULETTES"):
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		roulette.Main(ctx, b, update, current_user)
	}
}
