package roulette

import (
	"context"
	. "rr/kotatsutgbot/cb"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
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
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		is_participant := check_is_participant(current_user, current_anime_roulette)

		if current_anime_roulette.AnnounceDate.After(time.Now()) {
			SendMessageM(
				ctx, b, update,
				"roulette.menu", keyboards.CreateKeyboard_AnimeRouletteStart(is_participant),
			)
		} else {
			SendMessageM(
				ctx, b, update,
				"roulette.menu", keyboards.CreateKeyboard_AnimeRouletteMenu(is_participant),
			)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}
}

func Participate(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		is_participant := check_is_participant(current_user, current_anime_roulette)
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			if is_participant {
				SendMessageM(
					ctx, b, update,
					"roulette.already_participant", keyboards.CreateKeyboard_AnimeRouletteStart(true),
				)
			} else {
				db.DB_UPDATE_AnimeRoulette_ADD_Participants(current_user.ID)
				SendMessageM(
					ctx, b, update,
					"roulette.registered", keyboards.CreateKeyboard_AnimeRouletteStart(true),
				)
			}
		} else {
			if is_participant {
				SendMessageM(
					ctx, b, update,
					"roulette.already_participant", keyboards.CreateKeyboard_AnimeRouletteStart(true),
				)
			} else {
				SendMessageM(
					ctx, b, update,
					"roulette.registration_end", nil,
				)
			}
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}
}

func CancelParticipate(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		is_participant := check_is_participant(current_user, current_anime_roulette)

		if is_participant {
			db.DB_UPDATE_AnimeRoulette_REMOVE_Participants(current_user.ID)
			SendMessageM(
				ctx, b, update,
				"roulette.unregistered", keyboards.CreateKeyboard_AnimeRouletteStart(false),
			)
		} else {
			SendMessageM(
				ctx, b, update,
				"roulette.not_participant", keyboards.CreateKeyboard_AnimeRouletteStart(true),
			)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}
}

func AnimeWish(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			SendMessageM(
				ctx, b, update,
				"roulette.no_theme", nil,
			)
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
			if check_is_participant(current_user, current_anime_roulette) {
				UpdateStep(current_user, config.STEP_ANIME_RULETTE_ENTER_ENIGMATIC_TITLE)
				SendMessageM(
					ctx, b, update,
					"roulette.send_title", keyboards.Keyboard_CancelAnimeRoulette,
				)
			} else {
				SendMessageM(
					ctx, b, update,
					"roulette.not_participant", nil,
				)
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			SendMessageM(
				ctx, b, update,
				"roulette.ended", nil,
			)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}
}

func LinkMyList(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		if check_is_participant(current_user, current_anime_roulette) {
			UpdateStep(current_user, config.STEP_ANIME_RULETTE_ENTER_LINK_MY_ANIME_LIST)

			if current_user.LinkMyAnimeList == "" {
				SendMessageM(
					ctx, b, update,
					"roulette.send_list", nil,
				)
			} else {
				SendMessageMT(
					ctx, b, update,
					"roulette.your_list", current_user,
					keyboards.Keyboard_CancelAnimeRoulette,
				)
			}

		} else {
			SendMessageM(
				ctx, b, update,
				"roulette.not_participant", keyboards.Keyboard_CancelAnimeRoulette,
			)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}
}

func Rules(ctx context.Context, b *bot.Bot, update *models.Update) {
	SendMessageM(
		ctx, b, update,
		"roulette.rules", nil,
	)
}

func MainTheme(ctx context.Context, b *bot.Bot, update *models.Update) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		now := time.Now()
		if now.After(current_anime_roulette.StartDate) && now.Before(current_anime_roulette.AnnounceDate) {
			SendMessageM(
				ctx, b, update,
				"roulette.no_theme", nil,
			)
		} else if now.After(current_anime_roulette.AnnounceDate) && now.Before(current_anime_roulette.DistributionDate) {
			if current_anime_roulette.Theme == "" {
				SendMessageM(
					ctx, b, update,
					"roulette.almost_no_theme", nil,
				)
			} else {
				SendMessageRaw(
					ctx, b, update.Message.From.ID,
					current_anime_roulette.Theme, nil,
				)
			}
		} else if now.After(current_anime_roulette.DistributionDate) && now.Before(current_anime_roulette.EndDate) {
			SendMessageM(
				ctx, b, update,
				"roulette.ended", nil,
			)
		} else {
			SendMessageM(
				ctx, b, update,
				"roulette.registration_end", nil,
			)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}
}
