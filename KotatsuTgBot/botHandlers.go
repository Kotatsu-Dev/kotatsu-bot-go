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
	//Внутренние пакеты проекта
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/gen_certs"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"

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
	"regexp"
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
						case "⛩ Вступить в клуб":
							proccessText_JoinClub(ctx, b, update, user)

						case "📝 Запись на мероприятия":
							proccessText_SigningUpForActivity(ctx, b, update)

						case "📰 Подписаться на рассылку":
							proccessText_SubscribeNewsletter(ctx, b, update, user)

						case "❌ Отписаться от рассылки":
							proccessText_UnsubscribeNewsletter(ctx, b, update, user)

						case "📟 Связаться с клубом":
							proccessText_ContactClubManager(ctx, b, update, user)

						case "📟 Связь с клубом":
							proccessText_ContactClubManager(ctx, b, update, user)

						case "☎️ Связь с руководителем клуба":
							proccessText_ContactClubManager(ctx, b, update, user)

						case "⬅ Вернуться в меню":
							proccessText_BackMeinMenu(ctx, b, update, user)

						case "🚪 Покинуть клуб":
							proccessText_LeaveClub(ctx, b, update, user)

						case "📅 Мероприятия":
							proccessText_SigningUpForActivity(ctx, b, update)

						case "🤝 Акции и партнёры":
							proccessText_Partners(ctx, b, update)

						case "🟡 Аниме рулетка":
							processText_AnimeRoulette(ctx, b, update, user)

						case "⬅ Вернуться в меню рулетки":
							processText_AnimeRoulette(ctx, b, update, user)

						case "⬅️Вернуться в главное меню":
							proccessText_BackMeinMenu(ctx, b, update, user)

						case "✅ Участвовать в рулетке":
							processText_AnimeRoulette_Participate(ctx, b, update, user)

						case "🚪 Покинуть рулетку":
							processText_AnimeRoulette_CancelParticipate(ctx, b, update, user)

						case "❔ Загадать аниме":
							processText_AnimeRoulette_AnimeWish(ctx, b, update, user)

						case "🗞 Рассылка":
							proccessText_InDevelopment(ctx, b, update)

						case "📋 Правила":
							proccessText_AnimeRoulette_Rules(ctx, b, update)

						case "📔 Тема":
							proccessText_AnimeRoulette_MainTheme(ctx, b, update)

						case "📚 Мой список":
							proccessText_AnimeRoulette_LinkMyList(ctx, b, update, user)

						case "📂 Мои мероприятия":
							proccessText_MyActivities(ctx, b, update, user)

						case "⬅ Вернуться в главное меню":
							proccessText_BackMeinMenu(ctx, b, update, user)

						default:

							switch user.Step {
							case config.STEP_MESSAGE_SUPPORT:
								proccessStep_ContactClubManager(ctx, b, update, user)

							case config.STEP_ITMO_ENTER_ISU:
								proccessStep_ITMO_EnterISU(ctx, b, update, user, "join_club")

							case config.STEP_APPOINTMENT_ITMO_ENTER_ISU:
								proccessStep_ITMO_EnterISU(ctx, b, update, user, "activity")

							case config.STEP_ITMO_ENTER_FULLNAME:
								proccessStep_ITMO_EnterFullName(ctx, b, update, user, "join_club")

							case config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME:
								proccessStep_ITMO_EnterFullName(ctx, b, update, user, "activity")

							case config.STEP_ITMO_ENTER_SECRET_CODE:
								proccessStep_EnterSecretCode(ctx, b, update, user, "itmo")

							case config.STEP_NOITMO_ENTER_FULLNAME:
								proccessStep_NoITMO_EnterFullName(ctx, b, update, user, "join_club")

							case config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME:
								proccessStep_NoITMO_EnterFullName(ctx, b, update, user, "activity")

							case config.STEP_NOITMO_ENTER_PHONE:
								proccessStep_NoITMO_EnterPhoneNumber(ctx, b, update, user, "join_club")

							case config.STEP_CHANGING_PHONE:
								proccessStep_ChangePhoneNumber(ctx, b, update, user)

							case config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE:
								proccessStep_NoITMO_EnterPhoneNumber(ctx, b, update, user, "activity")

							case config.STEP_NOITMO_ENTER_SECRET_CODE:
								proccessStep_EnterSecretCode(ctx, b, update, user, "no_itmo")

							case config.STEP_USER_LEAVES_CLUB:
								proccessStep_LeavesClub(ctx, b, update, user)

							case config.STEP_ANIME_RUOLETTE_ENTER_ENIGMATIC_TITLE:
								proccessStep_AnimeRoulette_EnterEnigmaticTitle(ctx, b, update, user)

							case config.STEP_ANIME_RUOLETTE_ENTER_LINK_MY_ANIME_LIST:
								proccessStep_AnimeRoulette_EnterLinkMyAnimeList(ctx, b, update, user)

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

	if update.Message.Text == "🗃 Продолжить" {
		full_tg_name := update.Message.From.FirstName + " " + update.Message.From.LastName
		db_answer_reg := regUser(update.Message.From.ID, full_tg_name, update.Message.From.Username)

		switch db_answer_reg {
		case db.DB_ANSWER_SUCCESS:
			params.Text = "Главное меню"
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(false)

		case db.DB_ANSWER_OBJECT_EXISTS:
			params.Text = "Ты уже проходил(а) регистрацию в нашей системе" + "\n" +
				"Выбери интересующий тебя раздел:"

			_, old_user := db.DB_GET_User_BY_UserTgID(update.Message.From.ID)

			if old_user.IsClubMember {
				params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(old_user.IsSubscribeNewsletter)
			} else {
				params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(old_user.IsSubscribeNewsletter)
			}

		default:
			params.Text = "Произошла ошибка работы с БД"
			rr_debug.PrintLOG("main.go", "update.Message.Text", "activity_GetObjects()", "Ошибка работы с БД", "")
		}
	} else {
		params.Text = "Вы не зарегистрированны в системе" + "\n" +
			"Продолжая использование чат-бота, вы соглашаетесь на обработку персональных данных в соответствии с 152-ФЗ «О персональных данных»."
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

	if update.Message.Text == "🗃 Продолжить" {
		full_tg_name := update.CallbackQuery.From.FirstName + " " + update.CallbackQuery.From.LastName
		db_answer_reg := regUser(update.CallbackQuery.From.ID, full_tg_name, update.CallbackQuery.From.Username)

		switch db_answer_reg {
		case db.DB_ANSWER_SUCCESS:
			params.Text = "Главное меню"
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(false)

		case db.DB_ANSWER_OBJECT_EXISTS:
			params.Text = "Ты уже проходил(а) регистрацию в нашей системе" + "\n" +
				"Выбери интересующий тебя раздел:"

			_, old_user := db.DB_GET_User_BY_UserTgID(update.Message.From.ID)

			if old_user.IsClubMember {
				params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(old_user.IsSubscribeNewsletter)
			} else {
				params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(old_user.IsSubscribeNewsletter)
			}

		default:
			params.Text = "Произошла ошибка работы с БД"
			rr_debug.PrintLOG("main.go", "update.Message.Text", "activity_GetObjects()", "Ошибка работы с БД", "")
		}
	} else {
		params.Text = "Вы не зарегистрированны в системе" + "\n" +
			"Продолжая использование чат-бота, вы соглашаетесь на обработку персональных данных в соответствии с 152-ФЗ «О персональных данных»."
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

		var full_tg_name string

		full_tg_name = update.Message.From.FirstName + " " + update.Message.From.LastName

		params.Text = "Добро пожаловать: " + full_tg_name + "\n" +
			"Главное меню"

		if user.IsClubMember {
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(user.IsSubscribeNewsletter)
		} else {
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(user.IsSubscribeNewsletter)
		}

		_, err := b.SendMessage(ctx, params)
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Start", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
		}
	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		proccessRegistrationMessage(ctx, b, update)
	}
}

//
// Сообщения
//

// Вступление в клуб
func proccessText_JoinClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {

	params_load := &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
	}

	params := &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
	}

	if current_user.IsSentRequest {
		params.Text = "Ты уже отправил(а) заявку на вступление в клуб. Ожидай ответа от бота"

	} else {
		params_load.Text = "Инициализация..."
		params.Text = "Перед вступлением в клуб, пожалуйста, ознакомься с правилами:" + "\n" + "\n" +
			"1. У клуба открытый тип членства — достаточно проживать в Санкт-Петербурге и интересоваться аниме, мангой, ранобэ, JRPG, визуальными новеллами, косплеем или другими произведениями отаку-культуры." + "\n" + "\n" +
			"Уважай интересы и взгляды других участников. За разжигание ненависти и оскорбления можем исключить из клуба." + "\n" + "\n" +
			"2. Посещать все мероприятия клуба не обязательно — выбирай те, что приходятся тебе по душе. Но если мы не видели и не слышали тебя более 4 месяцев, членство в клубе может быть прекращено, но мы обязательно заранее свяжемся и предупредим. Вернуться можно в любой момент — это не бан и не наказание, а просто наш способ держать список участников актуальным, чтобы в нём не оставалось студентов, которые потеряли интерес к клубу или отчислились из ИТМО." + "\n" + "\n" +
			"Если участию в мероприятиях мешала учёба или работа — из клуба не исключаем, достаточно ответить на наше предупреждение. Мы тоже студенты, всё понимаем."
		params_load.ReplyMarkup = keyboards.CommunicationManager
		params.ReplyMarkup = keyboards.CreateInlineKbd_JoinClub()

		_, err_msg_load := b.SendMessage(ctx, params_load)
		if err_msg_load != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage(params_load)", "Ошибка отправки сообщения", err_msg_load.Error())
		}
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Запись на мероприятия
func proccessText_SigningUpForActivity(ctx context.Context, b *bot.Bot, update *models.Update) {

	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params_load := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params_photo := &bot.SendPhotoParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	status_one := false
	var active_activities_list []db.Activity_ReadJSON

	activities_list := db.DB_GET_Activities()

	if len(activities_list) == 0 {
		params.Text = "Никаких мероприятий на данный момент не запланированно"
		params.ReplyMarkup = keyboards.ListEvents

		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}

	} else {

		for _, activity_load := range activities_list {
			if activity_load.Status {
				status_one = true
				break
			}
		}

		if !status_one {
			params.Text = "Никаких мероприятий на данный момент не проводится"
			params.ReplyMarkup = keyboards.ListEvents

			_, err_msg := b.SendMessage(ctx, params)
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
			}
			return
		} else {
			params_load.Text = "Загрузка мероприятий..."
			params_load.ReplyMarkup = keyboards.ListEvents

			_, err_msg := b.SendMessage(ctx, params_load)
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
			}

			directory := "./img/calendar_activities"
			// Получите список файлов в каталоге
			files, err_dir := os.ReadDir(directory)
			if err_dir != nil {
				rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "os.ReadDir", "Ошибка поиска файла календаря", err_dir.Error())
			}

			fileInfo := files[0]
			filePath := filepath.Join(directory, fileInfo.Name())

			for _, activity := range activities_list {
				if activity.Status {
					active_activities_list = append(active_activities_list, activity)
				}
			}

			// Проверить наличие файла - Календарь мероприятий
			calendar_activities_path := filePath
			_, err := os.Stat(calendar_activities_path)
			if err == nil {

				// Открываем файл
				file, err := os.Open(filePath)
				if err != nil {
					rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "os.Open(filePath)", "Ошибка открытия файла календаря", err.Error())
					return
				}
				defer file.Close()

				// Создаем экземпляр InputFileUpload
				inputFile := &models.InputFileUpload{
					Filename: filepath.Base(filePath),
					Data:     file,
				}

				params_photo.Photo = inputFile
				params_photo.Caption = "Список текущих мероприятий:"
				params_photo.ReplyMarkup = keyboards.CreateInlineKbd_ActivitiesList(active_activities_list)

				// Отправляем фото
				_, err = b.SendPhoto(ctx, params_photo)
				if err != nil {
					rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "b.SendPhoto(ctx, params_photo)", "Ошибка отправки фото файла календаря", err.Error())
					return
				}

			} else if os.IsNotExist(err) {
				params.Text = "Список текущих мероприятий:"
				params.ReplyMarkup = keyboards.CreateInlineKbd_ActivitiesList(active_activities_list)

				_, err_msg := b.SendMessage(ctx, params_load)
				if err_msg != nil {
					rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
				}
			} else {
				rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "os.Stat", "Ошибка проверки наличия изображения мероприятий", err.Error())
				params.Text = "Список текущих мероприятий:"
				params.ReplyMarkup = keyboards.CreateInlineKbd_ActivitiesList(active_activities_list)

				_, err_msg := b.SendMessage(ctx, params_load)
				if err_msg != nil {
					rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
				}
			}
		}
	}
}

// Партнёры
func proccessText_Partners(ctx context.Context, b *bot.Bot, update *models.Update) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = "Список наших акций и партнёров"
	params.ReplyMarkup = keyboards.CreateInlineKbd_PartnersList()

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_Partners", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Мои мероприятия
func proccessText_MyActivities(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	var active_activities_list []*db.Activity

	if len(current_user.MyActivities) == 0 {
		params.Text = "Ты не записан(а) ни на одно мероприятие"
	} else {

		for _, activity := range current_user.MyActivities {
			if activity.Status {
				active_activities_list = append(active_activities_list, activity)
			}
		}

		params.Text = "Список мероприятий, на которые ты записан(а)"
		params.ReplyMarkup = keyboards.CreateInlineKbd_MyActivitiesList(active_activities_list)
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_MyActivities", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Подписаться на рассылку
func proccessText_SubscribeNewsletter(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = current_user.UserTgID
	update_user_data["is_subscribe_newsletter"] = true

	db.DB_UPDATE_User(update_user_data)

	params.Text = "Ты успешно был(а) подписан(а) на нашу рассылку!"
	params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(true)
	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_SubscribeNewsletter", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Отписаться от рассылки
func proccessText_UnsubscribeNewsletter(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = current_user.UserTgID
	update_user_data["is_subscribe_newsletter"] = false
	db.DB_UPDATE_User(update_user_data)

	params.Text = "Ты успешно был(а) подписан(а) на нашу рассылку!"
	params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(false)
	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_UnsubscribeNewsletter", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Связь с руководителем клуба
func proccessText_ContactClubManager(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = current_user.UserTgID
	update_user_data["step"] = config.STEP_MESSAGE_SUPPORT
	db.DB_UPDATE_User(update_user_data)

	params.Text = "Напиши своё обращение здесь, затем отправь его и оно будет направлено руководству клуба"
	params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("back")

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_ContactClubManager", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Вернуться в меню
func proccessText_BackMeinMenu(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = "Главное меню"

	if current_user.IsClubMember {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(current_user.IsSubscribeNewsletter)
	} else {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID
	update_user_data["step"] = config.STEP_DEFAULT
	db.DB_UPDATE_User(update_user_data)

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_BackMeinMenu", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Покинуть клуб
func proccessText_LeaveClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = current_user.UserTgID
	update_user_data["step"] = config.STEP_USER_LEAVES_CLUB
	db.DB_UPDATE_User(update_user_data)

	params.Text = "Введи причину выхода из клуба или нажми на кнопку 'Пропустить'"
	params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("skip")

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_LeaveClub", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Аниме рулетка
func processText_AnimeRoulette(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	is_participant := false

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		params.Text = "Меню рулетки"

		switch current_anime_roulette.CurrentStage {
		case config.ANIME_RUOLETTE_STAGE_START_REGISTRATION:
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

		default:
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

		}
		params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteMenu(is_participant)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Участвовать в рулетке
func processText_AnimeRoulette_Participate(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	is_participant := false

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		switch current_anime_roulette.CurrentStage {
		case config.ANIME_RUOLETTE_STAGE_START_REGISTRATION:
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

			if is_participant {
				params.Text = "Ты уже являешься участником аниме рулетки"

			} else {
				db.DB_UPDATE_AnimeRoulette_ADD_Participants(current_user.ID)
				params.Text = "Добро пожаловать в нашу аниме рулетку! Жди, когда появится тема аниме рулетки"
			}

			params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)

		default:
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

			if is_participant {
				params.Text = "Ты уже являешься участником аниме рулетки"
				params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)
			} else {
				params.Text = "К сожалению, набор учасников закончился. Возвращайтесь в аниме рулетку в следующий раз."
			}
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette_Participate", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Покинуть рулетку
func processText_AnimeRoulette_CancelParticipate(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	is_participant := false
	indexToRemove := -1

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		for i, participant := range current_anime_roulette.Participants {
			if current_user.UserTgID == participant.UserTgID {
				is_participant = true
				indexToRemove = i
				break
			}
		}

		if !is_participant {
			params.Text = "Ты уже покинул аниме рулетку"
		} else {
			if indexToRemove != -1 {
				db.DB_UPDATE_AnimeRoulette_REMOVE_Participants(current_user.ID)
				params.Text = "Вы покинули нашу аниме рулетку."
			}
		}

		params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette_CancelParticipate", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Загадать аниме
func processText_AnimeRoulette_AnimeWish(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	is_participant := false

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		switch current_anime_roulette.CurrentStage {
		case config.ANIME_RUOLETTE_STAGE_START_REGISTRATION:
			params.Text = "Тема пока не выдана. Ждите."

		case config.ANIME_RUOLETTE_STAGE_ANIME_GATHERING:

			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

			if is_participant {
				update_user_data := make(map[string]interface{})
				update_user_data["user_tg_id"] = update.Message.From.ID
				update_user_data["step"] = config.STEP_ANIME_RUOLETTE_ENTER_ENIGMATIC_TITLE
				db.DB_UPDATE_User(update_user_data)

				params.Text = "Введите и отправьте название аниме, которое вы хотите загадать"
				params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("anime_roulette")
			} else {
				params.Text = "Вы не являетесь участником рулетки."
			}

		case config.ANIME_RUOLETTE_STAGE_DATA_PROCESSING:
			params.Text = "Сбор названий аниме закончен."

		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette_AnimeWish", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Мой список
func proccessText_AnimeRoulette_LinkMyList(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	is_participant := false

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:

		for _, participant := range current_anime_roulette.Participants {
			if current_user.UserTgID == participant.UserTgID {
				is_participant = true
				break
			}
		}

		if is_participant {
			update_user_data := make(map[string]interface{})
			update_user_data["user_tg_id"] = update.Message.From.ID
			update_user_data["step"] = config.STEP_ANIME_RUOLETTE_ENTER_LINK_MY_ANIME_LIST
			db.DB_UPDATE_User(update_user_data)

			if current_user.LinkMyAnimeList == "" {
				params.Text = "Введите и отправьте ссылку на ваш список аниме, который вы хотите предложить"
			} else {
				params.Text = "Твой список аниме: " + current_user.LinkMyAnimeList + "\n" +
					"Хочешь изменить? Тогда укажи новую ссылку на свой список аниме"

				params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("anime_roulette")
			}

		} else {
			params.Text = "Вы не являетесь участником рулетки."
			params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("anime_roulette")
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"

	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_AnimeRoulette_LinkMyList", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Правила
func proccessText_AnimeRoulette_Rules(ctx context.Context, b *bot.Bot, update *models.Update) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: func(b bool) *bool { return &b }(true),
		},
	}

	params.Text = "Участники рулетки загадывают аниме по заданной теме и случайным образом получают для просмотра то, что загадал другой участник." + "\n" + "\n" +
		"Загадываемый тайтл должен иметь первый сезон не длиннее 30 серий, чтобы получивший его участник мог закончить просмотр в течение 3 недель." + "\n" +
		"Нельзя загадывать длинные франшизы (более 80 серий или 5 ТВ-сезонов), хентай и другие запрещённые в РФ тайтлы." + "\n" +
		"1 серия = 24 минуты." + "\n" + "\n" +
		"Если уже загаданный тайтл вы смотрели, то необходимо попросить замену." + "\n" + "\n" +
		"Цель рулетки: посмотреть загаданное аниме и написать отзыв в обсуждении: https://vk.com/topic-91030630_40877814." + "\n" +
		"Если вы решили бросить просмотр, то подробно опиши причину, иначе отзыв не засчитается." + "\n" +
		"За невыполнение цели следует наказание. И поверь, лучше до него не доводить: кто знает, что придётся выполнить в этот раз?"

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_AnimeRoulette_Rules", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Тема рулетки
func proccessText_AnimeRoulette_MainTheme(ctx context.Context, b *bot.Bot, update *models.Update) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		switch current_anime_roulette.CurrentStage {
		case config.ANIME_RUOLETTE_STAGE_START_REGISTRATION:
			params.Text = "Тема пока не выдана. Ждите."

		case config.ANIME_RUOLETTE_STAGE_ANIME_GATHERING:
			if current_anime_roulette.Theme == "" {
				params.Text = "Тему вот вот объявят"
			} else {
				params.Text = current_anime_roulette.Theme
			}

		case config.ANIME_RUOLETTE_STAGE_DATA_PROCESSING:
			params.Text = "Сбор названий аниме закончен."

		default:
			params.Text = "Аниме рулетка была проведена. Ждите следующий раз"
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_AnimeRoulette_MainTheme", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// В разработке
func proccessText_InDevelopment(ctx context.Context, b *bot.Bot, update *models.Update) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = "В разработке"

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_InDevelopment", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

//
// Шаги
//

// Шаг - Обращение к руководству клуба
func proccessStep_ContactClubManager(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {

	params_support := &bot.SendMessageParams{
		ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
		ParseMode: models.ParseModeHTML,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: func(b bool) *bool { return &b }(true),
		},
	}

	params_user := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = current_user.UserTgID
	update_user_data["step"] = config.STEP_DEFAULT
	db.DB_UPDATE_User(update_user_data)

	reference_number := generateRandomNumber(10)
	reference_number_str := strconv.Itoa(reference_number)

	params_user.Text = "Твоё сообщение успешно отправлено к руководству клуба." + "\n" +
		"Номер твоего обращения: " + reference_number_str

	params_user.ReplyMarkup = keyboards.CreateKeyboard_Cancel("back")

	user_name := update.Message.From.FirstName + " " + update.Message.From.LastName
	profileURL := fmt.Sprintf("https://t.me/%s", update.Message.From.Username)

	user_tg_id_str := strconv.FormatInt(update.Message.From.ID, 10)

	params_support.Text = "<b>Сообщение от пользователя</b>: " + user_name + "\n" +
		"<b>Ссылка на профиль пользователя</b>: " + profileURL + "\n" +
		"<b>Текст обращения</b>: " + update.Message.Text + "\n" +
		"<b>Ссылка для отправки ответа</b>: " + config.GetConfig().CONFIG_URL_BASE + "support-response/?user_tg_id=" + user_tg_id_str + "&reference_number=" + reference_number_str

	_, err_msg := b.SendMessage(ctx, params_support)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ContactClubManager", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
	}

	_, err_msg = b.SendMessage(ctx, params_user)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ContactClubManager", "b.SendMessage(ctx, params_user)", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - Человек из ИТМО вводит ИСУ
func proccessStep_ITMO_EnterISU(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {

	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	if _, err := strconv.Atoi(update.Message.Text); err == nil {

		update_user_data := make(map[string]interface{})
		update_user_data["user_tg_id"] = current_user.UserTgID
		update_user_data["isu"] = update.Message.Text

		if action == "join_club" {
			update_user_data["step"] = config.STEP_ITMO_ENTER_FULLNAME
		} else {
			update_user_data["step"] = config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME
		}

		db.DB_UPDATE_User(update_user_data)

		params.Text = "Введи свои ФИО"
	} else {
		params.Text = "Вы ввели номер ИСУ некорректно. Номер ИСУ не должен содержать буквы или иные символы." + "\n" +
			"Попробуйте ввести ещё раз или вернитесь в главное меню."
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterISU", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - Человек из ИТМО вводит ФИО
func proccessStep_ITMO_EnterFullName(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID
	update_user_data["full_name"] = update.Message.Text

	if action == "join_club" {
		params.Text = "Если у тебя есть код для вступления, отправь его" + "\n" +
			"Если кода нет, отправь цифру '0'"

		update_user_data["step"] = config.STEP_ITMO_ENTER_SECRET_CODE

	} else {
		update_user_data["step"] = config.STEP_DEFAULT
		update_user_data["is_itmo"] = true
		update_user_data["is_filled_data"] = true

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))

		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:

			db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

			params.Text = "Запись на мероприятие: " + activity.Title + " подтверждена"
			params.ReplyMarkup = keyboards.ListEvents
		}
	}

	db.DB_UPDATE_User(update_user_data)

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - Человек не из ИТМО вводит ФИО
func proccessStep_NoITMO_EnterFullName(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = "Введи свой номер мобильного телефона" + "\n" +
		"Он необходим для оформления пропуска на территорию Университета ИТМО, в котором проходят мероприятия клуба"

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = current_user.UserTgID
	update_user_data["full_name"] = update.Message.Text

	if action == "join_club" {
		update_user_data["step"] = config.STEP_NOITMO_ENTER_PHONE
	} else {
		update_user_data["step"] = config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE
	}

	db.DB_UPDATE_User(update_user_data)

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_NoITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - Человек не из ИТМО вводит номер мобильного телефона
func proccessStep_NoITMO_EnterPhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID

	// Регулярное выражение для валидации номера
	regex := regexp.MustCompile(`^(?:\+7|8)\d{10}$`)

	if regex.MatchString(update.Message.Text) {
		update_user_data["phone_number"] = update.Message.Text

		if action == "join_club" {
			params.Text = "Если у тебя есть код для вступления, отправь его" + "\n" +
				"Если кода нет, отправь цифру '0'"

			update_user_data["step"] = config.STEP_NOITMO_ENTER_SECRET_CODE
		} else {
			update_user_data["step"] = config.STEP_DEFAULT
			update_user_data["is_itmo"] = false
			update_user_data["is_filled_data"] = true

			db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:

				db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

				params.Text = "Запись на мероприятие: " + activity.Title + " подтверждена"
				params.ReplyMarkup = keyboards.ListEvents
			}
		}

		db.DB_UPDATE_User(update_user_data)

	} else {
		params.Text = "Номер введён некорректно" + "\n" +
			"Номер телефона должен иметь +7 или 8 в начале и 10 цифр после начала" + "\n" + "\n" +
			"Введи номер телефона ещё раз"
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_NoITMO_EnterPhoneNumber", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - человек вводит секретный код
func proccessStep_EnterSecretCode(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, status string) {

	params_user := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params_support := &bot.SendMessageParams{
		ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
		ParseMode: models.ParseModeHTML,
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID
	update_user_data["secret_code"] = update.Message.Text

	update_user_data["step"] = config.STEP_DEFAULT
	update_user_data["is_sent_request"] = true
	update_user_data["is_filled_data"] = true

	if status == "itmo" {
		update_user_data["is_itmo"] = true
	} else {
		update_user_data["is_itmo"] = false
	}

	db.DB_UPDATE_User(update_user_data)

	db_answer_code := db.DB_CREATE_Request(current_user.ID)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		params_user.Text = "Заявка на вступление в клуб была отправлена!" + "\n" +
			"Ожидай сообщение от бота — он уведомит о рассмотрении заявки"

		params_support.Text = "К нам поступила новая заявка на вступление в клуб от пользователя: " + current_user.FullName
		_, err_msg := b.SendMessage(ctx, params_support)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
		}

	default:
		params_user.Text = "Произошла системная ошибка при формировании заявки о вступлении в клуб." + "\n" +
			"Пожалуйста свяжиcь с нами и расскажи нам о данной проблеме."
	}

	params_user.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)

	_, err_msg := b.SendMessage(ctx, params_user)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_user)", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - человек меняет номер телефона
func proccessStep_ChangePhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	// Регулярное выражение для валидации номера
	regex := regexp.MustCompile(`^(?:\+7|8)\d{10}$`)

	if regex.MatchString(update.Message.Text) {

		update_user_data := make(map[string]interface{})
		update_user_data["user_tg_id"] = update.Message.From.ID
		update_user_data["phone_number"] = update.Message.Text

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:

			db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

			params.Text = "Твой новый номер был сохранён!" + "\n" +
				"Запись на мероприятие: " + activity.Title + " подтверждена"
			params.ReplyMarkup = keyboards.ListEvents
		}

		db.DB_UPDATE_User(update_user_data)

	} else {
		params.Text = "Номер введён некорректно" + "\n" +
			"Номер телефона должен иметь +7 или 8 в начале и 10 цифр после начала" + "\n" + "\n" +
			"Введи номер телефона ещё раз"
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ChangePhoneNumber", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - пользователь покидает клуб
func proccessStep_LeavesClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {

	params_user := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params_support := &bot.SendMessageParams{
		ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
		ParseMode: models.ParseModeHTML,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: func(b bool) *bool { return &b }(true),
		},
	}

	var user_isu_text string

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID

	switch update.Message.Text {
	case "Пропустить":
		update_user_data["is_club_member"] = false
		update_user_data["is_sent_request"] = false

		if current_user.ISU == "" {
			user_isu_text = "отсутствует, пользователь не из ИТМО"
		}

		params_support.Text = "Пользователь: " + current_user.FullName + " покинул наш клуб" + "\n" +
			"ИСУ пользователя: " + user_isu_text + "\n" +
			"Ссылка на Телеграмм: https://t.me/" + current_user.UserName + "\n" +
			"Причина: не указана"

		params_user.Text = "Ты покинул наш клуб"
		params_user.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)

		db.DB_UPDATE_User(update_user_data)

		_, err_msg := b.SendMessage(ctx, params_support)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_LeavesClub", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
		}

	case "⬅ Вернуться в главное меню":
		params_user.Text = "Главное меню"
		params_user.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)

	default:
		update_user_data["is_club_member"] = false
		update_user_data["is_sent_request"] = false

		if current_user.ISU == "" {
			user_isu_text = "отсутствует, пользователь не из ИТМО"
		}

		params_support.Text = "Пользователь: " + current_user.FullName + " покинул наш клуб" + "\n" +
			"ИСУ пользователя: " + user_isu_text + "\n" +
			"Ссылка на Телеграмм: https://t.me/" + current_user.UserName + "\n" +
			"Причина: " + update.Message.Text

		params_user.Text = "Ты покинул наш клуб"
		params_user.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)

		db.DB_UPDATE_User(update_user_data)

		_, err_msg := b.SendMessage(ctx, params_support)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_LeavesClub", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
		}
	}

	_, err_msg := b.SendMessage(ctx, params_user)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_LeavesClub", "b.SendMessage(ctx, params_user)", "Ошибка отправки сообщения", err_msg.Error())
	}
}

// Шаг - загадывание аниме
func proccessStep_AnimeRoulette_EnterEnigmaticTitle(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	is_participant := false

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		switch current_anime_roulette.CurrentStage {
		case config.ANIME_RUOLETTE_STAGE_START_REGISTRATION:
			params.Text = "Тема пока не выдана. Ждите."

		case config.ANIME_RUOLETTE_STAGE_ANIME_GATHERING:
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					update_user_data["enigmatic_title"] = update.Message.Text
					break
				}
			}

			if is_participant {
				db.DB_UPDATE_User(update_user_data)
				params.Text = "Аниме успешно загадано! Отлично!"

			} else {
				params.Text = "Вы не являетесь участником рулетки."
			}

		case config.ANIME_RUOLETTE_STAGE_DATA_PROCESSING:
			params.Text = "Сбор названий аниме закончен."

		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"
	}

	update_user_data["step"] = config.STEP_DEFAULT
	db.DB_UPDATE_User(update_user_data)

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_AnimeRoulette_EnterEnigmaticTitle", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}

	processText_AnimeRoulette(ctx, b, update, current_user)
}

// Шаг - предложить свой список аниме
func proccessStep_AnimeRoulette_EnterLinkMyAnimeList(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	is_participant := false

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		for _, participant := range current_anime_roulette.Participants {
			if current_user.UserTgID == participant.UserTgID {
				is_participant = true
				update_user_data["link_my_anime_list"] = update.Message.From.ID
				break
			}
		}

		if is_participant {
			db.DB_UPDATE_User(update_user_data)
			params.Text = "Ваша ссылка на список успешно принята! Отлично!"

		} else {
			params.Text = "Вы не являетесь участником рулетки."
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "На данный момент аниме рулетка не проводится"
	}

	update_user_data["step"] = config.STEP_DEFAULT
	db.DB_UPDATE_User(update_user_data)

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_AnimeRoulette_EnterLinkMyAnimeList", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}

	processText_AnimeRoulette(ctx, b, update, current_user)
}

// Неизвестное сообщение или шаг
func proccessText_Unknown(ctx context.Context, b *bot.Bot, update *models.Update) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = "Это команда которую я не знаю? Или сообщение админу, которое не понимаю?" + "\n" +
		"В любом случае используй команды из меню - я, бот, понимаю только их." + "\n" +
		"Для выхода в главное меню, нажми /start"

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_Unknown", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

//
// Inline - клавиатура
//

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

		if data == "from_ITMO" {
			update_user_data["step"] = config.STEP_ITMO_ENTER_ISU
			db.DB_UPDATE_User(update_user_data)

			params.Text = "Введи свой номер ИСУ"
		} else {
			update_user_data["step"] = config.STEP_NOITMO_ENTER_FULLNAME
			db.DB_UPDATE_User(update_user_data)

			params.Text = "Введи свои ФИО"
		}

		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_JOIN_CLUB", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}

	// Список акций и партнёров
	case strings.HasPrefix(update.CallbackQuery.Data, "PARTNERS"):
		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		params := &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		params_photo := &bot.SendPhotoParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
		}

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		switch data {
		case "cafeTaiyaki":
			last_name, first_name, paternity := splitName(current_user.FullName)
			user_tg_id_str := strconv.FormatInt(update.CallbackQuery.From.ID, 10)
			output_image_path := gen_certs.GenCerts(last_name, first_name, paternity, "./img/templates/taiyaki.png", user_tg_id_str)

			// Открываем файл
			file, err := os.Open(output_image_path)
			if err != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_PARTNERS", "os.Open(output_image_path)", "Ошибка открытия файла", err.Error())
				return
			}
			defer file.Close()

			// Создаем экземпляр InputFileUpload
			inputFile := &models.InputFileUpload{
				Filename: filepath.Base(output_image_path),
				Data:     file,
			}

			params_photo.Photo = inputFile

			// Отправляем фото
			_, err = b.SendPhoto(ctx, params_photo)
			if err != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_PARTNERS", "b.SendPhoto(ctx, params_photo)", "Ошибка отправки фото", err.Error())
				return
			}

		case "gemfest":
			output_image_path := config.FILE_PHOTO_GEMFEST_PATH
			// Открываем файл
			file, err := os.Open(output_image_path)
			if err != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_PARTNERS", "os.Open(output_image_path)", "Ошибка открытия файла", err.Error())
				return
			}
			defer file.Close()

			// Создаем экземпляр InputFileUpload
			inputFile := &models.InputFileUpload{
				Filename: filepath.Base(output_image_path),
				Data:     file,
			}

			params_photo.Photo = inputFile
			params_photo.Caption = "Приглашаем вас на новый мультифандомный аниме-фестиваль в Санкт-Петербурге https://vk.com/gemfestspb!" + "\n" +
				"Он будет посвящен Хэллоуину, а именно — теме Ковена." + "\n" + "\n" +
				"— 11 ноября с 12:00" + "\n" +
				"— Санкт-Петербург, Дом молодежи, Новоизмайловский пр. 48" + "\n" + "\n" +
				"Специально для нашего клуба — СКИДКА на любой из видов билетов по промокоду ITMOGEM23 до конца октября!" + "\n" + "\n" +
				"Пора достать из шкафов все самые жуткие наряды и отправиться навстречу приключениям!" + "\n" + "\n" +
				"🎫 Увидимся на Фестивале!: https://spb.qtickets.events/83613-gemfest-multifandomnyy-festival"

			// Отправляем фото
			_, err = b.SendPhoto(ctx, params_photo)
			if err != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_PARTNERS", "b.SendPhoto(ctx, params_photo)", "Ошибка отправки фото", err.Error())
				return
			}

		case "back":
			params.Text = "Ты вернулся в главное меню"
			_, err_msg := b.SendMessage(ctx, params)
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_PARTNERS", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
			}
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

		if data == "from_ITMO" {
			update_user_data["step"] = config.STEP_APPOINTMENT_ITMO_ENTER_ISU
			db.DB_UPDATE_User(update_user_data)

			params.Text = "Введи свой номер ИСУ"
		} else {
			update_user_data["step"] = config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME
			db.DB_UPDATE_User(update_user_data)

			params.Text = "Введи свои ФИО"
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
			var formattedTime string
			is_participant := false

			for _, participant := range activity.Participants {
				if participant.UserTgID == update.CallbackQuery.From.ID {
					is_participant = true
					break
				}
			}

			// Определите желаемый формат дд.мм чч:мм
			format := "02.01 15:04"

			// Используйте метод Format для форматирования времени
			formattedTime = activity.DateMeeting.Format(format)

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

				params.Text = "Подробнее о мероприятии: " + activity.Title + "\n" +
					"<b>Описание:</b> " + activity.Description + "\n" +
					"<b>Дата и время: </b>" + formattedTime + "\n" +
					"<b>Место проведения: </b>" + activity.Location

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

				params.Text = "Подробнее о мероприятии: " + activity.Title + "\n" +
					"<b>Описание:</b> " + activity.Description + "\n" +
					"<b>Дата и время: </b>" + formattedTime + "\n" +
					"<b>Место проведения: </b>" + activity.Location

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
					db.DB_UPDATE_Activity_ADD_Participants(uint(activity_id), current_user.ID)
					params_load.Text = "Загрузка..."
					params.Text = "Запись на мероприятие: " + activity.Title + " подтверждена"
					params.ReplyMarkup = keyboards.ListEvents
				}
			} else {
				params_load.Text = "Получение данных..."
				params.Text = "Твой номер телефона " + current_user.PhoneNumber + " является актуальным?"
				params_load.ReplyMarkup = keyboards.CreateKeyboard_Cancel("back")
				params.ReplyMarkup = keyboards.CreateInlineKbd_RelevancePhoneNumber()
			}
		} else {
			params_load.Text = "Инициализация..."
			params.Text = "Чтобы записаться на мероприятие, выбери один из вариантов ниже и предоставь нам нужные данные для записи"

			params_load.ReplyMarkup = keyboards.CreateKeyboard_Cancel("back")
			params.ReplyMarkup = keyboards.CreateInlineKbd_Appointment()
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
				params.Text = "Ты успешно отписался(ась) от мероприятия: " + activity.Title
				params.ReplyMarkup = keyboards.ListEvents

			case db.DB_ANSWER_OBJECT_NOT_FOUND:
				params.Text = "Мероприятие не найдено в базе данных!"
				params.ReplyMarkup = keyboards.ListEvents

			case db.DB_ANSWER_OBJECT_EXISTS:
				params.Text = "Ты изначально не был(а) записан на мероприятие: " + activity.Title
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
		update_user_data["user_tg_id"] = update.Message.From.ID

		parts = strings.Split(update.CallbackQuery.Data, "::")
		data = parts[1]

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
		if db_answer_code == db.DB_ANSWER_SUCCESS {
			if data == "yes" {
				db.DB_UPDATE_Activity_ADD_Participants(uint(activity.ID), current_user.ID)
				params.Text = "Запись на мероприятие: " + activity.Title + " подтверждена"
				params.ReplyMarkup = keyboards.ListEvents

			} else {
				update_user_data["step"] = config.STEP_CHANGING_PHONE
				db.DB_UPDATE_User(update_user_data)

				params.Text = "Укажи свой новый номер телефона в формате: +7 или 8 в начале, далее 10 цифр"
				params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("back")
			}

			_, err_msg := b.SendMessage(ctx, params)
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_RELEVANC_PHONE", "b.SendMessage(ctx, params)", "Ошибка отправки сообщения", err_msg.Error())
			}

		}

	}
}
