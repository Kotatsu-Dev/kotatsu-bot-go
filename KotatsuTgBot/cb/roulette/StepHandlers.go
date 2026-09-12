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

func EnterEnigmaticTitleE(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStepE(user, config.STEP_DEFAULT),
		GetActiveRoulette().
			Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
				return OneOf(
					RouletteStateGuard(roulette, RouletteStateRegistration, NoTheme()),
					RouletteStateGuard(roulette, RouletteStateWishing, ITE(
						check_is_participant(user, roulette),
						Seq(
							WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
								return UpdateUserE(user, map[string]any{
									"enigmatic_title": update.Message.Text,
								})
							}),
							SendMessageME(
								"roulette.sent_title", nil,
							),
						),
						NotParticipant(),
					)),
					RouletteEnded(),
				)
			}).
			Otherwise(RouletteInactive()),
		MainE(user),
	)
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
