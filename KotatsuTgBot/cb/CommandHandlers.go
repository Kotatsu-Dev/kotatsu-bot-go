package cb

import (
	"context"
	"io"
	"maps"
	"os"
	"path/filepath"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendMessageRaw(ctx context.Context, b *bot.Bot, chat_id int64, text string, keyboard models.ReplyMarkup) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chat_id,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendMessageRaw", "bot.SendMessage", "Ошибка отправки сообщения", err.Error())
	}

	return err
}

func SendMessage(ctx context.Context, b *bot.Bot, chat_id int64, text string, keyboard models.ReplyMarkup) error {
	return SendMessageRaw(ctx, b, chat_id, config.T(text), keyboard)
}

func SendMessageT(ctx context.Context, b *bot.Bot, chat_id int64, text string, data any, keyboard models.ReplyMarkup) error {
	return SendMessageRaw(ctx, b, chat_id, config.TT(text, data), keyboard)
}

func SendMessageM(ctx context.Context, b *bot.Bot, update *models.Update, text string, keyboard models.ReplyMarkup) error {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	return SendMessage(ctx, b, chat_id, text, keyboard)
}

func SendMessageMT(ctx context.Context, b *bot.Bot, update *models.Update, text string, data any, keyboard models.ReplyMarkup) error {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	return SendMessageT(ctx, b, chat_id, text, data, keyboard)
}

func SendPhotoRaw(ctx context.Context, b *bot.Bot, chat_id int64, filename string, file io.Reader, text string, keyboard models.ReplyMarkup) error {
	_, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:    chat_id,
		ParseMode: models.ParseModeHTML,
		Photo: &models.InputFileUpload{
			Filename: filename,
			Data:     file,
		},
		Caption:     text,
		ReplyMarkup: keyboard,
	})

	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "SendPhotoRaw", "bot.SendPhoto", "Ошибка отправки фотографии", err.Error())
	}

	return err
}

func SendPhoto(ctx context.Context, b *bot.Bot, chat_id int64, filename string, file io.Reader, text string, keyboard models.ReplyMarkup) error {
	return SendPhotoRaw(ctx, b, chat_id, filename, file, config.T(text), keyboard)
}

func SendPhotoM(ctx context.Context, b *bot.Bot, update *models.Update, filename string, file io.Reader, text string, keyboard models.ReplyMarkup) error {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}
	return SendPhoto(ctx, b, chat_id, filename, file, text, keyboard)
}

func UpdateCurrentUser(user *db.User_ReadJSON, update map[string]any) (int, *db.User, bool) {
	update_user := maps.Clone(update)
	update_user["user_tg_id"] = user.UserTgID
	return db.DB_UPDATE_User(update_user)
}

func UpdateGender(user *db.User_ReadJSON, gender db.Gender) (int, *db.User, bool) {
	return UpdateCurrentUser(user, map[string]any{
		"gender": gender,
	})
}

func UpdateVisited(user *db.User_ReadJSON, is_visited_events bool) (int, *db.User, bool) {
	return UpdateCurrentUser(user, map[string]any{
		"is_visited_events": is_visited_events,
	})
}

func UpdateStep(user *db.User_ReadJSON, step int) (int, *db.User, bool) {
	return UpdateCurrentUser(user, map[string]any{
		"step": step,
	})
}

func ITE[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func SendMainMenu(ctx context.Context, current_user *db.User_ReadJSON, b *bot.Bot, update *models.Update) {
	if current_user.IsClubMember {
		SendMessageM(
			ctx, b, update,
			"main_menu", keyboards.Keyboard_MainMenuButtonsClubMember,
		)
	} else {
		SendMessageM(
			ctx, b, update,
			"main_menu", keyboards.Keyboard_MainMenuButtonsDefault,
		)
	}
}

func SetGender(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, gender db.Gender) {
	UpdateGender(current_user, gender)
	SendMainMenu(ctx, current_user, b, update)
}

func OpenCalendar() (*os.File, string) {
	directory := config.ByUI("./img/calendar_activities")
	files, err := os.ReadDir(directory)
	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "OpenCalendar", "os.ReadDir", "Ошибка поиска файла календаря", err.Error())
	}

	if len(files) <= 0 {
		return nil, ""
	}

	fileInfo := files[0]
	filePath := filepath.Join(directory, fileInfo.Name())

	file, err := os.Open(filePath)
	if err != nil {
		rr_debug.PrintLOG("CommandHandlers.go", "OpenCalendar", "os.Open", "Ошибка чтения файла календаря", err.Error())
		return nil, ""
	}

	return file, fileInfo.Name()
}

// ---

func WasAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, actually bool) {
	UpdateVisited(current_user, actually)

	if actually {
		SendMessageM(
			ctx, b, update,
			"request.is_itmo", keyboards.InlineKbd_JoinClub,
		)
	} else {
		SendMessageM(
			ctx, b, update,
			"request.not_enough_visits", keyboards.Keyboard_WasntAtEvents,
		)
	}

}

func WasntAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, cont bool) {
	if cont {
		SendMessageM(
			ctx, b, update,
			"request.is_itmo", keyboards.InlineKbd_JoinClub,
		)
	} else {
		SendMainMenu(ctx, current_user, b, update)
	}
}

func JoinClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	if current_user.IsSentRequest {
		SendMessageM(
			ctx, b, update,
			"request.in_progress", nil,
		)
		SendMainMenu(ctx, current_user, b, update)
	} else if current_user.IsClubMember {
		SendMessageM(
			ctx, b, update,
			"request.already_accepted", nil,
		)
		SendMainMenu(ctx, current_user, b, update)
	} else {
		SendMessageM(
			ctx, b, update,
			"request.rules", keyboards.Keyboard_WasAtEvents,
		)
	}
}

func SigningUpForActivity(ctx context.Context, b *bot.Bot, update *models.Update) {
	activities_list := db.DB_GET_Active_Activities()

	status, _ := db.DB_GET_AnimeRoulette_BY_Status(true)
	has_roulette := status == db.DB_ANSWER_SUCCESS

	var text string
	var keyboard models.ReplyMarkup
	if len(activities_list) > 0 || has_roulette {
		text = "events.list"
		keyboard = keyboards.CreateInlineKbd_ActivitiesList(activities_list, update.Message.From.ID, has_roulette)
	} else {
		text = "events.empty"
	}

	file, name := OpenCalendar()
	if file != nil {
		SendPhotoM(
			ctx, b, update,
			name, file,
			text, keyboard,
		)
		file.Close()
	} else {
		SendMessageM(
			ctx, b, update,
			text, keyboard,
		)
	}
}

func BackMainMenu(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_DEFAULT)
	SendMainMenu(ctx, current_user, b, update)
}

func LeaveClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_USER_LEAVES_CLUB)
	SendMessageM(
		ctx, b, update,
		"leave_reason", keyboards.Keyboard_Skip,
	)
}

func MyActivities(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	activities := db.DB_GET_User_Active_Activities(current_user.ID)
	if len(activities) == 0 {
		SendMessageM(
			ctx, b, update,
			"my_events.empty", nil,
		)
	} else {
		SendMessageM(
			ctx, b, update,
			"my_events.list", keyboards.CreateInlineKbd_MyActivitiesList(activities),
		)
	}
}

func NoPhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_DEFAULT)
	SendMessageM(
		ctx, b, update,
		"request.no_phone_number", nil,
	)
	SendMainMenu(ctx, current_user, b, update)
}
