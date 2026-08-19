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
		cb.SendMessageMT(
			ctx, b, update,
			"welcome", update.Message.From,
			cb.ITE(
				user.IsClubMember,
				keyboards.Keyboard_MainMenuButtonsClubMember,
				keyboards.Keyboard_MainMenuButtonsDefault,
			),
		)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		proccessRegistrationMessage(ctx, b, update)
	}
}

// TODO
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
			formattedDate = cb.FormatDate(activity.DateMeeting.In(loc))

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

	cb.SendMessageM(
		ctx, b, update,
		"unknown_command", nil,
	)
}

func BotHandler_CallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	switch {
	case strings.HasPrefix(update.CallbackQuery.Data, "JOIN_CLUB"):
		cb.JoinClubQuery(ctx, b, update, current_user)

	case strings.HasPrefix(update.CallbackQuery.Data, "APPOINTMENT"):
		cb.AppointQuery(ctx, b, update, current_user)

	case strings.HasPrefix(update.CallbackQuery.Data, "ACTIVITIES"), strings.HasPrefix(update.CallbackQuery.Data, "MY_ACTIVITIES"):
		cb.ActivitiesQuery(ctx, b, update, current_user)

	case strings.HasPrefix(update.CallbackQuery.Data, "ACTIVITY_SUBSCRIBE"):
		cb.SubscribeQuery(ctx, b, update, current_user)

	case strings.HasPrefix(update.CallbackQuery.Data, "ACTIVITY_UNSUBSCRIBE"):
		cb.UnsubscribeQuery(ctx, b, update, current_user)

	case strings.HasPrefix(update.CallbackQuery.Data, "RELEVANC_PHONE"):
		cb.RelevancePhoneQuery(ctx, b, update, current_user)

	case strings.HasPrefix(update.CallbackQuery.Data, "ROULETTES"):
		roulette.MainQuery(ctx, b, update, current_user)
	}
}
