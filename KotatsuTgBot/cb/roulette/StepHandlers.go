package roulette

import (
	"context"
	"regexp"
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var linkToListRegexp = regexp.MustCompile(`^((https://)?anilist\.co/user/[A-Za-z0-9]+(/)?|(https://)?myanimelist\.net/profile/[A-Za-z0-9]+|(https://)?shikimori.one/[^/]+)$`)

func EnterEnigmaticTitle(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_DEFAULT)

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
				SendMessageM(
					ctx, b, update,
					"roulette.sent_title", nil,
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

	Main(ctx, b, update, current_user)
}

func EnterEnigmaticTitleE(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStepE(user, config.STEP_DEFAULT),
		GetActiveRoulette().
			Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
				now := time.Now()
				if now.After(roulette.StartDate) && now.Before(roulette.AnnounceDate) {
					return NoTheme()
				} else if now.After(roulette.AnnounceDate) && now.Before(roulette.DistributionDate) {
					return ITE(
						check_is_participant(user, roulette),
						SendMessageME(
							"roulette.sent_title", nil,
						),
						NotParticipant(),
					)
				}
				return RouletteEnded()
			}).
			Otherwise(RouletteInactive()),
		MainE(user),
	)
}

func EnterLinkMyAnimeList(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_DEFAULT)

	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		if check_is_participant(current_user, current_anime_roulette) {
			link_to_list := update.Message.Text
			if linkToListRegexp.MatchString(link_to_list) {
				UpdateCurrentUser(current_user, map[string]any{
					"link_my_anime_list": link_to_list,
				})
				SendMessageM(
					ctx, b, update,
					"roulette.sent_list", nil,
				)
			} else {
				SendMessageM(
					ctx, b, update,
					"roulette.incorrect_list_format", nil,
				)
			}
		} else {
			SendMessageM(
				ctx, b, update,
				"roulette.not_participant", nil,
			)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}

	Main(ctx, b, update, current_user)
}

func EnterLinkMyAnimeListE(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStepE(user, config.STEP_DEFAULT),
		GetActiveRoulette().
			Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
				return ITE(
					check_is_participant(user, roulette),
					MatchText(linkToListRegexp).
						Then(func(text string) Executor {
							return Seq(
								UpdateUserE(user, map[string]any{
									"link_my_anime_list": text,
								}),
								SendMessageME(
									"roulette.sent_list", nil,
								),
							)
						}).
						Otherwise(SendMessageME(
							"roulette.incorrect_list_format", nil,
						)).(Executor),
					NotParticipant(),
				)
			}).
			Otherwise(RouletteInactive()),
		MainE(user),
	)
}
