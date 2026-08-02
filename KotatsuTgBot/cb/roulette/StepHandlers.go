package roulette

import (
	"context"
	"regexp"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var linkToListRegexp = regexp.MustCompile(`^((https://)?anilist\.co/user/[A-Za-z0-9]+(/)?|(https://)?myanimelist\.net/profile/[A-Za-z0-9]+|(https://)?shikimori.one/[^/]+)$`)

func EnterEnigmaticTitle(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			text = config.T("roulette.no_theme")
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
			if check_is_participant(current_user, current_anime_roulette) {
				text = config.T("roulette.sent_title")
			} else {
				text = config.T("roulette.not_participant")
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			text = config.T("roulette.ended")
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	db.DB_UPDATE_User(map[string]any{
		"user_tg_id": update.Message.From.ID,
		"step":       config.STEP_DEFAULT,
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
		Text:      text,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_AnimeRoulette_EnterEnigmaticTitle", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}

	Main(ctx, b, update, current_user)
}

func EnterLinkMyAnimeList(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string

	update_user_data := map[string]any{
		"user_tg_id": update.Message.From.ID,
		"step":       config.STEP_DEFAULT,
	}

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		if check_is_participant(current_user, current_anime_roulette) {
			link_to_list := update.Message.Text
			if linkToListRegexp.MatchString(link_to_list) {
				update_user_data["link_my_anime_list"] = link_to_list
				text = config.T("roulette.sent_list")
			} else {
				text = config.T("roulette.incorrect_list_format")
			}
		} else {
			text = config.T("roulette.not_participant")
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		text = config.T("roulette.inactive")
	}

	db.DB_UPDATE_User(update_user_data)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
		Text:      text,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_AnimeRoulette_EnterLinkMyAnimeList", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}

	Main(ctx, b, update, current_user)
}
