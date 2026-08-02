package cb

import (
	"context"
	"os"
	"path/filepath"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SetGender(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, gender db.Gender) {
	db.DB_UPDATE_User(map[string]any{
		"user_tg_id": current_user.UserTgID,
		"gender":     gender,
	})

	var keyboard models.ReplyMarkup
	if current_user.IsClubMember {
		keyboard = keyboards.Keyboard_MainMenuButtonsClubMember
	} else {
		keyboard = keyboards.Keyboard_MainMenuButtonsDefault
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        config.T("main_menu"),
		ReplyMarkup: keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func WasAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, actually bool) {
	db.DB_UPDATE_User(map[string]any{
		"user_tg_id":        current_user.UserTgID,
		"is_visited_events": actually,
	})

	var text string
	var keyboard models.ReplyMarkup
	if actually {
		text = config.T("request.is_itmo")
		keyboard = keyboards.InlineKbd_JoinClub
	} else {
		text = config.T("request.not_enough_visits")
		keyboard = keyboards.Keyboard_WasntAtEvents
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func WasntAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, cont bool) {
	var text string
	var keyboard models.ReplyMarkup

	if cont {
		text = config.T("request.is_itmo")
		keyboard = keyboards.InlineKbd_JoinClub
	} else {
		text = config.T("main_menu")

		if current_user.IsClubMember {
			keyboard = keyboards.Keyboard_MainMenuButtonsClubMember
		} else {
			keyboard = keyboards.Keyboard_MainMenuButtonsDefault
		}
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func JoinClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup
	if current_user.IsSentRequest {
		text = config.T("request.in_progress")
		if current_user.IsClubMember {
			keyboard = keyboards.Keyboard_MainMenuButtonsClubMember
		} else {
			keyboard = keyboards.Keyboard_MainMenuButtonsDefault
		}
	} else if current_user.IsClubMember {
		text = config.T("request.already_accepted")
		keyboard = keyboards.Keyboard_MainMenuButtonsClubMember
	} else {
		text = config.T("request.rules")
		keyboard = keyboards.Keyboard_WasAtEvents
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func SigningUpForActivity(ctx context.Context, b *bot.Bot, update *models.Update) {
	var active_activities_list []db.Activity_ReadJSON

	activities_list := db.DB_GET_Activities()

	// TODO: Filter by query
	for _, activity := range activities_list {
		if activity.Status {
			active_activities_list = append(active_activities_list, activity)
		}
	}

	status, _ := db.DB_GET_AnimeRoulette_BY_Status(true)
	has_roulette := status == db.DB_ANSWER_SUCCESS

	var text string
	var keyboard models.ReplyMarkup
	if len(active_activities_list) > 0 || has_roulette {
		text = config.T("events.list")
		keyboard = keyboards.CreateInlineKbd_ActivitiesList(active_activities_list, update.Message.From.ID, has_roulette)
	} else {
		text = config.T("events.empty")
	}

	// TODO: Handle absence of calendar file
	directory := config.ByUI("./img/calendar_activities")
	files, err_dir := os.ReadDir(directory)
	if err_dir != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "os.ReadDir", "Ошибка поиска файла календаря", err_dir.Error())
	}

	fileInfo := files[0]
	filePath := filepath.Join(directory, fileInfo.Name())

	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()

		inputFile := &models.InputFileUpload{
			Filename: filepath.Base(filePath),
			Data:     file,
		}

		_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID:      update.Message.From.ID,
			ParseMode:   models.ParseModeHTML,
			Photo:       inputFile,
			Caption:     text,
			ReplyMarkup: keyboard,
		})
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "b.SendPhoto(ctx, params_photo)", "Ошибка отправки фото файла календаря", err.Error())
			return
		}
	} else {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_SigningUpForActivity", "os.Stat", "Ошибка проверки наличия изображения мероприятий", err.Error())
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			ParseMode:   models.ParseModeHTML,
			Text:        text,
			ReplyMarkup: keyboard,
		})
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessCommand_Unknown", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
		}
	}
}

func BackMainMenu(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db.DB_UPDATE_User(map[string]any{
		"user_tg_id": update.Message.From.ID,
		"step":       config.STEP_DEFAULT,
	})

	var keyboard models.ReplyMarkup
	if current_user.IsClubMember {
		keyboard = keyboards.Keyboard_MainMenuButtonsClubMember
	} else {
		keyboard = keyboards.Keyboard_MainMenuButtonsDefault
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        config.T("main_menu"),
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_BackMeinMenu", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func LeaveClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db.DB_UPDATE_User(map[string]any{
		"user_tg_id": update.Message.From.ID,
		"step":       config.STEP_USER_LEAVES_CLUB,
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        config.T("leave_reason"),
		ReplyMarkup: keyboards.Keyboard_Skip,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_LeaveClub", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func MyActivities(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup
	if len(current_user.MyActivities) == 0 {
		text = config.T("my_events.empty")
	} else {
		var active_activities_list []*db.Activity
		// TODO: Filter at db side
		for _, activity := range current_user.MyActivities {
			if activity.Status {
				active_activities_list = append(active_activities_list, activity)
			}
		}

		text = config.T("my_events.list")
		keyboard = keyboards.CreateInlineKbd_MyActivitiesList(active_activities_list)
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_MyActivities", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func NoPhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var keyboard models.ReplyMarkup
	if current_user.IsClubMember {
		keyboard = keyboards.Keyboard_MainMenuButtonsClubMember
	} else {
		keyboard = keyboards.Keyboard_MainMenuButtonsDefault
	}

	db.DB_UPDATE_User(map[string]any{
		"user_tg_id": update.Message.From.ID,
		"step":       config.STEP_DEFAULT,
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        config.T("request.no_phone_number"),
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_BackMeinMenu", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}
