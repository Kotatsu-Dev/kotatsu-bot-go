package roulette

import (
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
)

func check_is_participant(user *db.User_ReadJSON, roulette *db.AnimeRoulette_ReadJSON) bool {
	for _, participant := range roulette.Participants {
		if user.UserTgID == participant.UserTgID {
			return true
		}
	}
	return false
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

func RulesE() Executor {
	return SendMessageME(
		"roulette.rules", nil,
	)
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
