package cb

import (
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/cb/roulette"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
)

// TODO: otherwise
var MAIN = GetCurrentUserE().
	Then(func(user *db.User_ReadJSON) Executor {
		return OneOf(
			TextGuard("keyboard.gender_male", SetGenderE(user, "male")),
			TextGuard("keyboard.gender_female", SetGenderE(user, "female")),
			TextGuard("keyboard.visited_enough", WasAtEventsE(user, true)),
			TextGuard("keyboard.not_visited_enough", WasAtEventsE(user, false)),
			TextGuard("keyboard.fill_back_later", WasntAtEventsE(user, false)),
			TextGuard("keyboard.fill_now", WasntAtEventsE(user, true)),
			TextGuard("keyboard.join_club", JoinClubE(user)),
			TextGuard("keyboard.event_registration", SigningUpForActivityE(user)),
			TextGuard("keyboard.to_main_menu", BackMainMenuE(user)),
			TextGuard("keyboard.leave_club", LeaveClubE(user)),
			TextGuard("keyboard.to_roulette_menu", roulette.MainE(user)),
			TextGuard("keyboard.participate_roulette", roulette.ParticipateE(user)),
			TextGuard("keyboard.leave_roulette", roulette.CancelParticipateE(user)),
			TextGuard("keyboard.send_title", roulette.AnimeWishE(user)),
			TextGuard("keyboard.roulette_rules", roulette.RulesE()),
			TextGuard("keyboard.roulette_theme", roulette.MainThemeE()),
			TextGuard("keyboard.roulette_list", roulette.LinkMyListE(user)),
			TextGuard("keyboard.my_events", MyActivitiesE(user)),
			TextGuard("keyboard.not_my_number", NoPhoneNumberE(user)),

			StepGuard(user, config.STEP_ITMO_ENTER_ISU, ITMO_EnterISUE(user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_ITMO_ENTER_ISU, ITMO_EnterISUE(user, "activity")),
			StepGuard(user, config.STEP_ITMO_ENTER_FULLNAME, ITMO_EnterFullNameE(user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME, ITMO_EnterFullNameE(user, "activity")),
			StepGuard(user, config.STEP_NOITMO_ENTER_FULLNAME, NoITMO_EnterFullNameE(user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME, NoITMO_EnterFullNameE(user, "activity")),
			StepGuard(user, config.STEP_NOITMO_ENTER_PHONE, NoITMO_EnterPhoneNumberE(user, "join_club")),
			StepGuard(user, config.STEP_CHANGING_PHONE, ChangePhoneNumberE(user)),
			StepGuard(user, config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE, NoITMO_EnterPhoneNumberE(user, "activity")),
			StepGuard(user, config.STEP_USER_LEAVES_CLUB, LeavesClubE(user)),
			StepGuard(user, config.STEP_ANIME_RULETTE_ENTER_ENIGMATIC_TITLE, roulette.EnterEnigmaticTitleE(user)),
			StepGuard(user, config.STEP_ANIME_RULETTE_ENTER_LINK_MY_ANIME_LIST, roulette.EnterLinkMyAnimeListE(user)),
		)
	})

	// TODO: Otherwise
var CALLBACK_MAIN = GetCurrentUserE().
	Then(func(user *db.User_ReadJSON) Executor {
		return OneOf(
			QueryGuard("JOIN_CLUB", JoinClubQueryE(user)),
			QueryGuard("APPOINTMENT", AppointQueryE(user)),
			QueryGuard("ACTIVITIES", ActivitiesQueryE(user)),
			QueryGuard("MY_ACTIVITIES", ActivitiesQueryE(user)),
			QueryGuard("ACTIVITY_SUBSCRIBE", SubscribeQueryE(user)),
			QueryGuard("ACTIVITY_UNSUBSCRIBE", UnsubscribeQueryE(user)),
			QueryGuard("RELEVANC_PHONE", RelevancePhoneQueryE(user)),
			QueryGuard("ROULETTES", roulette.MainQueryE(user)),
		)
	})
