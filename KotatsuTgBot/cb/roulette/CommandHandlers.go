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
	return SendMessageM(
		"roulette.inactive", nil,
	)
}

func StartMenu(is_participant bool) Executor {
	return SendMessageM(
		"roulette.menu", keyboards.CreateKeyboard_AnimeRouletteStart(is_participant),
	)
}

func OngoingMenu(is_participant bool) Executor {
	return SendMessageM(
		"roulette.menu", keyboards.CreateKeyboard_AnimeRouletteMenu(is_participant),
	)
}

func AlreadyParticipant() Executor {
	return SendMessageM(
		"roulette.already_participant", keyboards.CreateKeyboard_AnimeRouletteStart(true),
	)
}

func RegistrationEnd() Executor {
	return SendMessageM(
		"roulette.registration_end", nil,
	)
}

func NotParticipant() Executor {
	return SendMessageM(
		"roulette.not_participant", keyboards.CreateKeyboard_AnimeRouletteStart(false),
	)
}

func RouletteEnded() Executor {
	return SendMessageM(
		"roulette.ended", nil,
	)
}

func NoTheme() Executor {
	return SendMessageM(
		"roulette.no_theme", nil,
	)
}

func Main(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			is_participant := check_is_participant(user, roulette)
			return If(
				GetRouletteState(roulette) <= RouletteStateRegistration,
				StartMenu(is_participant),
				OngoingMenu(is_participant),
			)
		}).
		Otherwise(RouletteInactive())
}

func Participate(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			if check_is_participant(user, roulette) {
				return AlreadyParticipant()
			}
			return FirstMatch(
				RouletteStateGuard(roulette, RouletteStateRegistration, Do(
					AddRouletteParticipant(user),
					SendMessageM(
						"roulette.registered", keyboards.CreateKeyboard_AnimeRouletteStart(true),
					),
				)),
				RegistrationEnd(),
			)
		}).
		Otherwise(RouletteInactive())
}

func CancelParticipate(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			if check_is_participant(user, roulette) {
				return Do(
					RemoveRouletteParticipant(user),
					SendMessageM(
						"roulette.unregistered", keyboards.CreateKeyboard_AnimeRouletteStart(false),
					),
				)
			} else {
				return NotParticipant()
			}
		}).
		Otherwise(RouletteInactive())
}

func AnimeWish(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			return FirstMatch(
				RouletteStateGuard(roulette, RouletteStateRegistration, NoTheme()),
				RouletteStateGuard(roulette, RouletteStateWishing, If(
					check_is_participant(user, roulette),
					Do(
						UpdateStep(user, config.STEP_ANIME_RULETTE_ENTER_ENIGMATIC_TITLE),
						SendMessageM(
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

func LinkMyList(user *db.User_ReadJSON) Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			return If(
				check_is_participant(user, roulette),
				Do(
					UpdateStep(user, config.STEP_ANIME_RULETTE_ENTER_LINK_MY_ANIME_LIST),
					If(
						user.LinkMyAnimeList == "",
						SendMessageM(
							"roulette.send_list", nil,
						),
						SendMessageMT(
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

func Rules() Executor {
	return SendMessageM(
		"roulette.rules", nil,
	)
}

func MainTheme() Executor {
	return GetActiveRoulette().
		Then(func(roulette *db.AnimeRoulette_ReadJSON) Executor {
			return FirstMatch(
				RouletteStateGuard(roulette, RouletteStateRegistration, NoTheme()),
				RouletteStateGuard(roulette, RouletteStateWishing, If(
					roulette.Theme == "",
					SendMessageM(
						"roulette.almost_no_theme", nil,
					),
					SendMessageRawM(
						roulette.Theme, nil,
					),
				)),
				RouletteStateGuard(roulette, RouletteStateDistributed, RouletteEnded()),
				RegistrationEnd(),
			)
		}).
		Otherwise(RouletteInactive())
}
