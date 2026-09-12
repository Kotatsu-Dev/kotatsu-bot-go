package roulette

import (
	"context"
	"regexp"
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var linkToListRegexp = regexp.MustCompile(`^((https://)?anilist\.co/user/[A-Za-z0-9]+(/)?|(https://)?myanimelist\.net/profile/[A-Za-z0-9]+|(https://)?shikimori.one/[^/]+)$`)

func EnterEnigmaticTitle(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStep(user, config.STEP_DEFAULT),
		GetActiveRoulette().
			Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
				return OneOf(
					RouletteStateGuard(roulette, RouletteStateRegistration, NoTheme()),
					RouletteStateGuard(roulette, RouletteStateWishing, If(
						check_is_participant(user, roulette),
						Seq(
							WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
								return UpdateUser(user, map[string]any{
									"enigmatic_title": update.Message.Text,
								})
							}),
							SendMessageM(
								"roulette.sent_title", nil,
							),
						),
						NotParticipant(),
					)),
					RouletteEnded(),
				)
			}).
			Otherwise(RouletteInactive()),
		Main(user),
	)
}

func EnterLinkMyAnimeList(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStep(user, config.STEP_DEFAULT),
		GetActiveRoulette().
			Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
				return If(
					check_is_participant(user, roulette),
					MatchText(linkToListRegexp).
						Then(func(text string) Executor {
							return Seq(
								UpdateUser(user, map[string]any{
									"link_my_anime_list": text,
								}),
								SendMessageM(
									"roulette.sent_list", nil,
								),
							)
						}).
						Otherwise(SendMessageM(
							"roulette.incorrect_list_format", nil,
						)),
					NotParticipant(),
				)
			}).
			Otherwise(RouletteInactive()),
		Main(user),
	)
}
