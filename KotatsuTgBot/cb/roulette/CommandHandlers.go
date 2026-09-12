package roulette

import (
	"context"
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"

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

		if GetRouletteState(current_anime_roulette) == RouletteStateRegistration {
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

func RouletteInactive() Executor {
	return SendMessageME(
		"roulette.inactive", nil,
	)
}

func StartMenu(is_participant bool) Executor {
	return SendMessageME(
		"roulette.menu", keyboards.CreateKeyboard_AnimeRouletteStart(is_participant),
	)
}

func OngoingMenu(is_participant bool) Executor {
	return SendMessageME(
		"roulette.menu", keyboards.CreateKeyboard_AnimeRouletteMenu(is_participant),
	)
}

func AlreadyParticipant() Executor {
	return SendMessageME(
		"roulette.already_participant", keyboards.CreateKeyboard_AnimeRouletteStart(true),
	)
}

func RegistrationEnd() Executor {
	return SendMessageME(
		"roulette.registration_end", nil,
	)
}

func NotParticipant() Executor {
	return SendMessageME(
		"roulette.not_participant", keyboards.CreateKeyboard_AnimeRouletteStart(false),
	)
}

func RouletteEnded() Executor {
	return SendMessageME(
		"roulette.ended", nil,
	)
}

func NoTheme() Executor {
	return SendMessageME(
		"roulette.no_theme", nil,
	)
}

func MainE(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			is_participant := check_is_participant(user, roulette)
			return ITE(
				GetRouletteState(roulette) == RouletteStateRegistration,
				StartMenu(is_participant),
				OngoingMenu(is_participant),
			)
		}).
		Otherwise(RouletteInactive())
}

func Participate(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		if check_is_participant(current_user, current_anime_roulette) {
			SendMessageM(
				ctx, b, update,
				"roulette.already_participant", keyboards.CreateKeyboard_AnimeRouletteStart(true),
			)
			return
		}

		if GetRouletteState(current_anime_roulette) == RouletteStateRegistration {
			db.DB_UPDATE_AnimeRoulette_ADD_Participants(current_user.ID)
			SendMessageM(
				ctx, b, update,
				"roulette.registered", keyboards.CreateKeyboard_AnimeRouletteStart(true),
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

func ParticipateE(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			if check_is_participant(user, roulette) {
				return AlreadyParticipant()
			}
			return OneOf(
				RouletteStateGuard(roulette, RouletteStateRegistration, Seq(
					AddRouletteParticipant(user),
					SendMessageME(
						"roulette.registered", keyboards.CreateKeyboard_AnimeRouletteStart(true),
					),
				)),
				RegistrationEnd(),
			)
		}).
		Otherwise(RouletteInactive())
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
				"roulette.not_participant", keyboards.CreateKeyboard_AnimeRouletteStart(false),
			)
		}

	case db.DB_ANSWER_OBJECT_NOT_FOUND:
		SendMessageM(
			ctx, b, update,
			"roulette.inactive", nil,
		)
	}
}

func CancelParticipateE(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			if check_is_participant(user, roulette) {
				return Seq(
					RemoveRouletteParticipant(user),
					SendMessageME(
						"roulette.unregistered", keyboards.CreateKeyboard_AnimeRouletteStart(false),
					),
				)
			} else {
				return NotParticipant()
			}
		}).
		Otherwise(RouletteInactive())
}

func AnimeWish(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		switch GetRouletteState(current_anime_roulette) {
		case RouletteStateRegistration:
			SendMessageM(
				ctx, b, update,
				"roulette.no_theme", nil,
			)
		case RouletteStateWishing:
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
		default:
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

func AnimeWishE(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			return OneOf(
				RouletteStateGuard(roulette, RouletteStateRegistration, NoTheme()),
				RouletteStateGuard(roulette, RouletteStateWishing, ITE(
					check_is_participant(user, roulette),
					Seq(
						UpdateStepE(user, config.STEP_ANIME_RULETTE_ENTER_ENIGMATIC_TITLE),
						SendMessageME(
							"roulette.send_title", keyboards.Keyboard_CancelAnimeRoulette,
						),
					),
					NotParticipant(),
				)),
				RouletteEnded(),
			)
		}).
		Otherwise(RouletteInactive())
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

func LinkMyListE(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			return ITE(
				check_is_participant(user, roulette),
				Seq(
					UpdateStepE(user, config.STEP_ANIME_RULETTE_ENTER_LINK_MY_ANIME_LIST),
					ITE(
						user.LinkMyAnimeList == "",
						SendMessageME(
							"roulette.send_list", nil,
						),
						SendMessageMTE(
							"roulette.your_list", user,
							keyboards.Keyboard_CancelAnimeRoulette,
						),
					),
				),
				NotParticipant(),
			)
		}).
		Otherwise(RouletteInactive())
}

func Rules(ctx context.Context, b *bot.Bot, update *models.Update) {
	SendMessageM(
		ctx, b, update,
		"roulette.rules", nil,
	)
}

func RulesE() Executor {
	return SendMessageME(
		"roulette.rules", nil,
	)
}

func MainTheme(ctx context.Context, b *bot.Bot, update *models.Update) {
	db_answer_code, current_anime_roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		switch GetRouletteState(current_anime_roulette) {
		case RouletteStateRegistration:
			SendMessageM(
				ctx, b, update,
				"roulette.no_theme", nil,
			)
		case RouletteStateWishing:
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
		case RouletteStateDistributed:
			SendMessageM(
				ctx, b, update,
				"roulette.ended", nil,
			)
		default:
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

func MainThemeE() Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			return OneOf(
				RouletteStateGuard(roulette, RouletteStateRegistration, NoTheme()),
				RouletteStateGuard(roulette, RouletteStateWishing, ITE(
					roulette.Theme == "",
					SendMessageME(
						"roulette.almost_no_theme", nil,
					),
					SendMessageRawME(
						roulette.Theme, nil,
					),
				)),
				RouletteStateGuard(roulette, RouletteStateDistributed, RouletteEnded()),
				RegistrationEnd(),
			)
		}).
		Otherwise(RouletteInactive())
}
