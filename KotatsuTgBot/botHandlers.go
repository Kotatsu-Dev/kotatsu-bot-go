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
	"rr/kotatsutgbot/cb/helpers"
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

//
// Главные процессы
//

func BotHandler_Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil {
		return
	}

	if update.Message != nil {
		cb.MAIN.Execute(ctx, b, update)
	} else {
		cb.CALLBACK_MAIN.Execute(ctx, b, update)
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
		helpers.SendMessageMT(
			ctx, b, update,
			"welcome", update.Message.From,
			helpers.ITE(
				user.IsClubMember,
				keyboards.Keyboard_MainMenuButtonsClubMember,
				keyboards.Keyboard_MainMenuButtonsDefault,
			),
		)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		cb.ProccessRegistrationE().Execute(ctx, b, update)
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
		cb.ProccessRegistrationE().Execute(ctx, b, update)
	}
}

