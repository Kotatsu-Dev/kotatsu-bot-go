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
	"rr/kotatsutgbot/middleware"
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
	"regexp"
	"strconv"
	"strings"
)

var fullNameRegexp = regexp.MustCompile(`([А-Яа-яЁё]+)\s([А-Яа-яЁё]+)\s([А-Яа-яЁё]+)`)

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
						case "Повелитель демонов":
							proccessText_SetGender(ctx, b, update, user, "male")
						case "Девочка волшебница":
							proccessText_SetGender(ctx, b, update, user, "female")

						case "Да, я уже мандаринка":
							proccessText_WasAtEvents(ctx, b, update, user, true)
						case "Ещё нет :(":
							proccessText_WasAtEvents(ctx, b, update, user, false)
						case "Хорошо, заполню позже":
							proccessText_WasntAtEvents(ctx, b, update, user, false)
						case "Хочу продолжить":
							proccessText_WasntAtEvents(ctx, b, update, user, true)
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

						case "Я не пользуюсь номером, к которому привязан Telegram":
							proccessText_NoPhoneNumber(ctx, b, update, user)

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
			params.Text = "Кто ты?"
			params.ReplyMarkup = keyboards.Keyboard_GenderSelect

		case db.DB_ANSWER_OBJECT_EXISTS:
			params.Text = "Привет! Мы уже знакомы, можешь выбирать нужный раздел."

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
		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:   update.Message.Chat.ID,
			Document: &models.InputFileString{Data: "CAACAgIAAx0CbgUG4QACCWpostfAVRPNDHNAWu8vcIbjv0nuagACrXQAAl8iQUmAFQIjshq4bTYE"},
		})
		params.Text = "Продолжая общение со мной, ты соглашаешься на обработку персональных данных в соответствии со 152-ФЗ «О персональных данных»."
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
			params.Text = "Привет! Мы уже знакомы, можешь выбирать нужный раздел."

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
		params.Text = "Привет!" + "\n" +
			"Продолжая общение со мной, ты соглашаешься на обработку персональных данных в соответствии со 152-ФЗ «О персональных данных»."
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

		params.Text = "Окаэринасай, " + full_tg_name

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

func BotHandler_Command_Login(ctx context.Context, b *bot.Bot, update *models.Update) {
	url := config.GetConfig().CONFIG_URL_BASE + "/login?" +
		middleware.CreateSessionCookie(strconv.FormatInt(update.Message.From.ID, 10), 24*time.Hour)

	middleware.CreateSessionCookie(strconv.FormatInt(update.Message.From.ID, 10), 24*time.Hour)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
		Text:      "Для входа нажмите кнопку",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{
						Text: "Войти", URL: url,
					},
				},
			},
		},
	})
	fmt.Println(err)
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
		params.Text = "Твою заявку ещё не обработали. Пожалуйста, подожди ответа руководителя или напиши сообщение в канал @anime_itmo (значок чата внизу канала)"

	} else {
		params.Text = "Перед вступлением в клуб немного о правилах:\n" +
			"0. Для посещения большинства мероприятий вступать в клуб не обязательно.\n" +
			"Если хочешь просто к нам прийти, перейди в меню «Запись на мероприятия»\n" +
			"1. Чтобы вступить в клуб, посети хотя бы 3 мероприятия. Онлайн-встречи тоже считаются :)\n" +
			"2. Относись ко всем участникам с уважением. Никого нельзя унижать за их интересы и вкусы\n" +
			"3. Наш клуб — официальная структура в ИТМО, поэтому не забывай о правилах Университета.\n\n" +
			"<a href=\"https://kotatsu.spb.ru/rules/current.pdf\">Полные правила</a> (там скучно и намного официальнее, но больше деталей)\n\n" +
			"Уже посетил(а) 3 наших мероприятия?"
		params.ParseMode = models.ParseModeHTML
		params_load.ReplyMarkup = keyboards.CommunicationManager
		// params.ReplyMarkup = keyboards.CreateInlineKbd_JoinClub()
		params.ReplyMarkup = keyboards.Keyboard_WasAtEvents

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

func proccessText_SetGender(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, gender db.Gender) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
	}

	db.DB_UPDATE_User(map[string]interface{}{
		"user_tg_id": current_user.UserTgID,
		"gender":     gender,
	})

	params.Text = "Главное меню"

	if current_user.IsClubMember {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(current_user.IsSubscribeNewsletter)
	} else {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)
	}

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

func proccessText_WasAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, actually bool) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
	}

	if actually {
		params.Text = "Подскажи, ты учишься или работаешь в ИТМО?"
		params.ReplyMarkup = keyboards.CreateInlineKbd_JoinClub()
	} else {
		params.Text = "К сожалению, вступить без посещения хотя бы 3 мероприятий не выйдет.\n" +
			"Пожалуйста, заполни заявку на вступление после того, как познакомишься с нами поближе.\n" +
			"Можешь продолжить заполнение заявки, тогда тебе напишет рук. клуба."

		params.ReplyMarkup = keyboards.Keyboard_WasntAtEvents
	}

	db.DB_UPDATE_User(map[string]interface{}{
		"user_tg_id":        current_user.UserTgID,
		"is_visited_events": actually,
	})

	_, err_msg := b.SendMessage(ctx, params)
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

func proccessText_WasntAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, cont bool) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
	}

	if cont {
		params.Text = "Подскажи, ты учишься или работаешь в ИТМО?"
		params.ReplyMarkup = keyboards.CreateInlineKbd_JoinClub()
	} else {
		params.Text = "Главное меню"

		if current_user.IsClubMember {
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(current_user.IsSubscribeNewsletter)
		} else {
			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)
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

	var active_activities_list []db.Activity_ReadJSON

	activities_list := db.DB_GET_Activities()

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
		if len(active_activities_list) > 0 {
			params_photo.Caption = "Список текущих мероприятий:"
			params_photo.ReplyMarkup = keyboards.CreateInlineKbd_ActivitiesList(active_activities_list, update.Message.From.ID)
		} else {
			params_photo.Caption = "Сейчас нет мероприятий, на которые я могу тебя записать." + "\n" +
				"Если в канале был анонс мероприятия, проверь, нет ли там ссылки на запись."
		}

		// Отправляем фото
		_, err = b.SendPhoto(ctx, params_photo)
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "b.SendPhoto(ctx, params_photo)", "Ошибка отправки фото файла календаря", err.Error())
			return
		}

	} else if os.IsNotExist(err) {
		if len(active_activities_list) > 0 {
			params.Text = "Список текущих мероприятий:"
			params.ReplyMarkup = keyboards.CreateInlineKbd_ActivitiesList(active_activities_list, update.Message.From.ID)
		} else {
			params.Text = "Сейчас нет мероприятий, на которые я могу тебя записать." + "\n" +
				"Если в канале был анонс мероприятия, проверь, нет ли там ссылки на запись."
		}

		_, err_msg := b.SendMessage(ctx, params_load)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}
	} else {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "os.Stat", "Ошибка проверки наличия изображения мероприятий", err.Error())
		if len(active_activities_list) > 0 {
			params.Text = "Список текущих мероприятий:"
			params.ReplyMarkup = keyboards.CreateInlineKbd_ActivitiesList(active_activities_list, update.Message.From.ID)
		} else {
			params.Text = "Сейчас нет мероприятий, на которые я могу тебя записать." + "\n" +
				"Если в канале был анонс мероприятия, проверь, нет ли там ссылки на запись."
		}

		_, err_msg := b.SendMessage(ctx, params_load)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "bot.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
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
		params.Text = "Сейчас ты не записан(а) ни на одно мероприятие"
	} else {

		for _, activity := range current_user.MyActivities {
			if activity.Status {
				active_activities_list = append(active_activities_list, activity)
			}
		}

		params.Text = "Я записывала тебя на эти мероприятия"
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

	_, user := db.DB_UPDATE_User(update_user_data)

	params.Text = "Теперь я буду присылать тебе важные сообщения от клуба прямо в этот чат"
	if user != nil && user.IsClubMember {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(true)
	} else {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(true)
	}
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
	_, user := db.DB_UPDATE_User(update_user_data)

	params.Text = "Хорошо-хорошо, больше не буду :("
	if user != nil && user.IsClubMember {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsClubMember(false)
	} else {
		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(false)
	}
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

func proccessText_NoPhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = "Я не могу записать тебя по номеру, не привязанному к аккаунту в Telegram :(\n" +
		"Пожалуйста, напиши в сообщения канала @anime_itmo (значок чата внизу канала), руководитель поможет тебе с записью и пропуском"

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

	params.Text = "Пожалуйста, напиши причину, по которой хочешь покинуть клуб, или нажми на кнопку «Пропустить»"
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

		for _, participant := range current_anime_roulette.Participants {
			if current_user.UserTgID == participant.UserTgID {
				is_participant = true
				break
			}
		}

		params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteMenu(is_participant)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"
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
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

			if is_participant {
				params.Text = "Ты уже участвуешь в рулетке"

			} else {
				db.DB_UPDATE_AnimeRoulette_ADD_Participants(current_user.ID)
				params.Text = "Теперь ты участник рулетки! Скоро я вышлю тему, на которую нужно будет загадать аниме."
			}

			params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)
		} else {
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

			if is_participant {
				params.Text = "Ты уже участвуешь в рулетке"
				params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)
			} else {
				params.Text = "К сожалению, набор участников закончился. Следи за анонсами в канале @anime_itmo, чтобы не пропустить следующую рулетку."
			}
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"
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
			params.Text = "Ты не участвуешь в рулетке :("
		} else {
			if indexToRemove != -1 {
				db.DB_UPDATE_AnimeRoulette_REMOVE_Participants(current_user.ID)
				params.Text = "Теперь ты не участвуешь в рулетке :("
			}
		}

		params.ReplyMarkup = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"
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
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			params.Text = "Ещё рано — я объявлю тему позже"
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
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

				params.Text = "Отправь мне название аниме, которое хочешь загадать"
				params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("anime_roulette")
			} else {
				params.Text = "Ты не участвуешь в рулетке :("
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			params.Text = "Сбор тайтлов уже закончился"
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"
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
				params.Text = "Отправь ссылку на свой список аниме"
			} else {
				params.Text = "Твой список аниме: " + current_user.LinkMyAnimeList + "\n" +
					"Хочешь изменить? Отправь новую ссылку."

				params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("anime_roulette")
			}

		} else {
			params.Text = "Ты не участвуешь в рулетке :("
			params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("anime_roulette")
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"

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
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			params.Text = "Ещё рано — я объявлю тему позже"
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
			if current_anime_roulette.Theme == "" {
				params.Text = "Ещё чуть-чуть — скоро объявлю тему"
			} else {
				params.Text = current_anime_roulette.Theme
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			params.Text = "Сбор тайтлов уже закончился"
		} else {
			params.Text = "К сожалению, набор участников закончился. Следи за анонсами в канале @anime_itmo, чтобы не пропустить следующую рулетку."
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"
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
		"<b>TG URL</b>: " + profileURL + "\n" +
		"<b>Текст обращения</b>: " + "\n" + update.Message.Text + "\n" +
		"<b>Ссылка для ответа</b>: " + config.GetConfig().CONFIG_URL_BASE + "support-response/?user_tg_id=" + user_tg_id_str + "&reference_number=" + reference_number_str

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
		params.Text = "Это не номер ИСУ!" + "\n" +
			"Попробуй ещё раз или напиши в сообщения канала @anime_itmo (значок чата внизу канала), руководитель поможет тебе."
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

	matched := fullNameRegexp.MatchString(update.Message.Text)

	if !matched {
		params.Text = "Неправильный формат ФИО, попробуй ещё раз в формате Фамилия Имя Отчество."
		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}
		return
	}

	update_user_data := make(map[string]interface{})
	update_user_data["user_tg_id"] = update.Message.From.ID
	update_user_data["full_name"] = update.Message.Text

	if action == "join_club" {
		params_support := &bot.SendMessageParams{
			ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
			ParseMode: models.ParseModeHTML,
		}

		update_user_data := make(map[string]interface{})
		update_user_data["user_tg_id"] = update.Message.From.ID
		update_user_data["secret_code"] = "0"

		update_user_data["step"] = config.STEP_DEFAULT
		update_user_data["is_sent_request"] = true
		update_user_data["is_filled_data"] = true

		update_user_data["is_itmo"] = true

		db.DB_UPDATE_User(update_user_data)

		db_answer_code := db.DB_CREATE_Request(current_user.ID)
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			params.Text = "Отправила твою заявку руководителю клуба." + "\n" +
				"Ожидай сообщение от меня, или если у нас появятся вопросы — от руководителя клуба."

			params_support.Text = "НОВАЯ ЗАЯВКА НА ВСТУПЛЕНИЕ" + "\n" + current_user.FullName
			_, err_msg := b.SendMessage(ctx, params_support)
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
			}

		default:
			params.Text = "Упс, кажется, у меня ошибка." + "\n" +
				"Напиши в сообщения канала @anime_itmo (значок чата внизу канала) и сообщи об ошибке."
		}

		params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)

	} else {
		update_user_data["step"] = config.STEP_DEFAULT
		update_user_data["is_itmo"] = true
		update_user_data["is_filled_data"] = true

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))

		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:

			db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

			params.Text = "Я записала тебя на мероприятие «" + activity.Title + "»"
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

	matched := fullNameRegexp.MatchString(update.Message.Text)

	if !matched {
		params.Text = "Неправильный формат ФИО, попробуй ещё раз в формате Фамилия Имя Отчество."
		_, err_msg := b.SendMessage(ctx, params)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}
		return
	}

	params.Text = "Введи свой номер мобильного телефона" + "\n" +
		"Он необходим для оформления пропуска на территорию Университета ИТМО, в котором проходят мероприятия клуба"
	params.ReplyMarkup = keyboards.CreateKeyboard_RequestContact()

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

	phone_number := ""

	if update.Message.Contact != nil {
		phone_number = update.Message.Contact.PhoneNumber
	}

	if phone_number != "" {
		update_user_data["phone_number"] = phone_number
		params.ReplyMarkup = keyboards.CreateKeyboard_Cancel("")

		if action == "join_club" {
			params_support := &bot.SendMessageParams{
				ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
				ParseMode: models.ParseModeHTML,
			}

			update_user_data := make(map[string]interface{})
			update_user_data["user_tg_id"] = update.Message.From.ID
			update_user_data["secret_code"] = "0"

			update_user_data["step"] = config.STEP_DEFAULT
			update_user_data["is_sent_request"] = true
			update_user_data["is_filled_data"] = true

			update_user_data["is_itmo"] = false

			db.DB_UPDATE_User(update_user_data)

			db_answer_code := db.DB_CREATE_Request(current_user.ID)
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:
				params.Text = "Отправила твою заявку руководителю клуба." + "\n" +
					"Ожидай сообщение от меня, или если у нас появятся вопросы — от руководителя клуба."

				params_support.Text = "НОВАЯ ЗАЯВКА НА ВСТУПЛЕНИЕ" + "\n" + current_user.FullName
				_, err_msg := b.SendMessage(ctx, params_support)
				if err_msg != nil {
					rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
				}

			default:
				params.Text = "Упс, кажется, у меня ошибка." + "\n" +
					"Напиши в сообщения канала @anime_itmo (значок чата внизу канала) и сообщи об ошибке."
			}

			params.ReplyMarkup = keyboards.CreateKeyboard_MainMenuButtonsDefault(current_user.IsSubscribeNewsletter)
		} else {
			update_user_data["step"] = config.STEP_DEFAULT
			update_user_data["is_itmo"] = false
			update_user_data["is_filled_data"] = true

			db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:

				db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

				params.Text = "Я записала тебя на мероприятие " + activity.Title
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
		params_user.Text = "Отправила твою заявку руководителю клуба." + "\n" +
			"Ожидай сообщение от меня, или если у нас появятся вопросы — от руководителя клуба."

		params_support.Text = "НОВАЯ ЗАЯВКА НА ВСТУПЛЕНИЕ" + "\n" + current_user.FullName
		_, err_msg := b.SendMessage(ctx, params_support)
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
		}

	default:
		params_user.Text = "Упс, кажется, у меня ошибка." + "\n" +
			"Напиши в сообщения канала @anime_itmo (значок чата внизу канала) и сообщи об ошибке."
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

	phone_number := ""

	if update.Message.Contact != nil {
		phone_number = update.Message.Contact.PhoneNumber
	}

	if phone_number != "" {

		update_user_data := make(map[string]interface{})
		update_user_data["user_tg_id"] = update.Message.From.ID
		update_user_data["phone_number"] = phone_number

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:

			db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

			params.Text = "Я сохранила твой номер и записала тебя на мероприятие «" + activity.Title + "»"
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
			user_isu_text = "Не из ИТМО"
		}

		params_support.Text = "Пользователь " + current_user.FullName + " покинул наш клуб" + "\n" +
			"ИСУ: " + user_isu_text + "\n" +
			"TG URL: https://t.me/" + current_user.UserName + "\n" +
			"Причина выхода не была указана"

		params_user.Text = "Жаль, что ты уходишь :(\n" +
			"Я передам запрос руководителю, он удалит запись в ИСУ в течение 3 дней.\n" +
			"Не забывай, что к нам можно приходить даже без членства в клубе — просто следи за анонсами встреч и не забывай на них записываться."
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
			user_isu_text = "Не из ИТМО"
		}

		params_support.Text = "Пользователь " + current_user.FullName + " покинул наш клуб" + "\n" +
			"ИСУ: " + user_isu_text + "\n" +
			"TG URL: https://t.me/" + current_user.UserName + "\n" +
			"Указанная причина: " + update.Message.Text

		params_user.Text = "Жаль, что ты уходишь :(\n" +
			"Я передам запрос руководителю, он удалит запись в ИСУ в течение 3 дней.\n" +
			"Не забывай, что к нам можно приходить даже без членства в клубе — просто следи за анонсами встреч и не забывай на них записываться."
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
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			params.Text = "Ещё рано — я объявлю тему позже"
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					update_user_data["enigmatic_title"] = update.Message.Text
					break
				}
			}

			if is_participant {
				db.DB_UPDATE_User(update_user_data)
				params.Text = "Я записала твой тайтл. Интересно, кому он выпадет?"

			} else {
				params.Text = "Ты не участвуешь в рулетке :("
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			params.Text = "Сбор тайтлов уже закончился"
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"
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
			params.Text = "Спасибо, я сохранила твой список."

		} else {
			params.Text = "Ты не участвуешь в рулетке :("
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		params.Text = "Сейчас аниме-рулетка не проводится"
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
	if !(update.Message != nil && update.Message.Chat.Type == models.ChatTypePrivate) {
		return
	}

	params := &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
	}

	params.Text = "Я не знаю такую команду." + "\n" +
		"Пожалуйста, используй команды из меню, я понимаю только их." + "\n" +
		"Для выхода в главное меню напиши /start"

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
			var formattedTime, formattedDate string
			is_participant := false

			for _, participant := range activity.Participants {
				if participant.UserTgID == update.CallbackQuery.From.ID {
					is_participant = true
					break
				}
			}

			// Определите желаемый формат дд.мм чч:мм

			// Используйте метод Format для форматирования времени
			formattedTime = activity.DateMeeting.Format("15:04")
			formattedDate = formatDate(activity.DateMeeting)

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

				params.Text = fmt.Sprintf("<b>%s</b>\n\n"+
					"%s\n\n"+
					"📅 <b>%s</b>\n"+
					"🕒 <b>%s</b>\n"+
					"📍 <b>%s</b>",
					activity.Title,
					activity.Description,
					formattedDate,
					formattedTime,
					activity.Location)

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

				params.Text = fmt.Sprintf("<b>%s</b>\n\n"+
					"%s\n\n"+
					"📅 <b>%s</b>\n"+
					"🕒 <b>%s</b>\n"+
					"📍 <b>%s</b>",
					activity.Title,
					activity.Description,
					formattedDate,
					formattedTime,
					activity.Location)

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
					db.DB_UPDATE_Activity_ADD_Participants(uint(activity_id), current_user.ID)
					params.Text = "Я записала тебя на мероприятие «" + activity.Title + "»"
					params.ReplyMarkup = keyboards.ListEvents

					db.DB_UPDATE_User(map[string]interface{}{
						"user_tg_id": current_user.UserTgID,
						"step":       config.STEP_DEFAULT,
					})
				}
			} else {
				params.Text = fmt.Sprintf("В прошлый раз ты указывал(а) номер %s.\n"+
					"В день мероприятия обязательно возьми телефон и паспорт с собой — с этого номера нужно позвонить на терминал для печати пропуска, а паспорт может попросить охрана.",
					current_user.PhoneNumber)
				params_load.ReplyMarkup = keyboards.CreateKeyboard_Cancel("back")
				params.ReplyMarkup = keyboards.CreateInlineKbd_RelevancePhoneNumber()

				fmt.Println(activity_id)

				db.DB_UPDATE_User(map[string]interface{}{
					"user_tg_id":       current_user.UserTgID,
					"step":             config.STEP_DEFAULT,
					"temp_activity_id": int(activity_id),
				})
			}
		} else {
			params.Text = "Кажется, мы с тобой ещё не знакомы. Подскажи, ты учишься/работаешь в ИТМО?"

			params_load.ReplyMarkup = keyboards.CreateKeyboard_Cancel("back")
			params.ReplyMarkup = keyboards.CreateInlineKbd_Appointment()

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
				params.Text = "Хорошо, я отменила твою запись на «" + activity.Title + "»"
				params.ReplyMarkup = keyboards.ListEvents

			case db.DB_ANSWER_OBJECT_NOT_FOUND:
				params.Text = "Такого мероприятия нет!"
				params.ReplyMarkup = keyboards.ListEvents

			case db.DB_ANSWER_OBJECT_EXISTS:
				params.Text = "Но ведь ты и так не записан(а) на это мероприятие..."
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
				db.DB_UPDATE_Activity_ADD_Participants(uint(activity.ID), current_user.ID)
				params.Text = "Я записала тебя на «" + activity.Title + "»"
				params.ReplyMarkup = keyboards.ListEvents
				update_user_data["step"] = config.STEP_DEFAULT

			} else {
				update_user_data["step"] = config.STEP_CHANGING_PHONE
				db.DB_UPDATE_User(update_user_data)

				params.Text = "Нажми «Отправить номер», чтобы поделиться со мной контактом"
				params.ReplyMarkup = keyboards.CreateKeyboard_RequestContact()
			}

			_, err_msg := b.SendMessage(ctx, params)
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "BotHandler_CallbackQuery_RELEVANC_PHONE", "b.SendMessage(ctx, params)", "Ошибка отправки сообщения", err_msg.Error())
			}

		}

	}
}
