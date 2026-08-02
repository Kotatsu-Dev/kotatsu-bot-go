package roulette

import (
	"context"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func check_is_participant(user *db.User_ReadJSON, roulette *db.AnimeRoulette_ReadJSON) bool {
	for _, participant := range roulette.Participants {
		if user.UserTgID == participant.UserTgID {
			return true
		}
	}
	return false
}

func Main(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var chat_id int64
	if update.Message != nil {
		chat_id = update.Message.From.ID
	} else {
		chat_id = update.CallbackQuery.From.ID
	}

	var text string
	var keyboard models.ReplyMarkup

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		text = config.T("roulette.menu")

		is_participant := check_is_participant(current_user, current_anime_roulette)

		if current_anime_roulette.AnnounceDate.After(time.Now()) {
			keyboard = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)
		} else {
			keyboard = keyboards.CreateKeyboard_AnimeRouletteMenu(is_participant)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chat_id,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func Participate(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup
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
				text = config.T("roulette.already_participant")
			} else {
				db.DB_UPDATE_AnimeRoulette_ADD_Participants(current_user.ID)
				text = config.T("roulette.registered")
				is_participant = true
			}

			keyboard = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)
		} else {
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

			if is_participant {
				text = config.T("roulette.already_participant")
				keyboard = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)
			} else {
				text = config.T("roulette.registration_end")
			}
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette_Participate", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func CancelParticipate(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup
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

		if !is_participant {
			text = config.T("roulette.not_participant")
		} else {
			db.DB_UPDATE_AnimeRoulette_REMOVE_Participants(current_user.ID)
			text = config.T("roulette.unregistered")
			is_participant = false
		}

		keyboard = keyboards.CreateKeyboard_AnimeRouletteStart(is_participant)

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	_, err_msg := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette_CancelParticipate", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

func AnimeWish(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup
	is_participant := false

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			text = config.T("roulette.no_theme")
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
			for _, participant := range current_anime_roulette.Participants {
				if current_user.UserTgID == participant.UserTgID {
					is_participant = true
					break
				}
			}

			if is_participant {
				db.DB_UPDATE_User(map[string]any{
					"user_tg_id": update.Message.From.ID,
					"step":       config.STEP_ANIME_RUOLETTE_ENTER_ENIGMATIC_TITLE,
				})

				text = config.T("roulette.send_title")
				keyboard = keyboards.Keyboard_CancelAnimeRoulette
			} else {
				text = config.T("roulette.not_participant")
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			text = config.T("roulette.ended")
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	_, err_msg := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err_msg != nil {
		rr_debug.PrintLOG("botHandlers.go", "processText_AnimeRoulette_AnimeWish", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
	}
}

func LinkMyList(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup
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
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id": update.Message.From.ID,
				"step":       config.STEP_ANIME_RUOLETTE_ENTER_LINK_MY_ANIME_LIST,
			})

			if current_user.LinkMyAnimeList == "" {
				text = config.T("roulette.send_list")
			} else {
				text = config.TT("roulette.your_list", current_user)
				keyboard = keyboards.Keyboard_CancelAnimeRoulette
			}

		} else {
			text = config.T("roulette.not_participant")
			keyboard = keyboards.Keyboard_CancelAnimeRoulette
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_AnimeRoulette_LinkMyList", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func Rules(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: bot.True(),
		},
		Text: config.T("roulette.rules"),
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_AnimeRoulette_Rules", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func MainTheme(ctx context.Context, b *bot.Bot, update *models.Update) {
	var text string
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			text = config.T("roulette.no_theme")
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
			if current_anime_roulette.Theme == "" {
				text = config.T("roulette.almost_no_theme")
			} else {
				text = current_anime_roulette.Theme
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			text = config.T("roulette.ended")
		} else {
			text = config.T("roulette.registration_end")
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
		Text:      text,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessText_AnimeRoulette_MainTheme", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}
