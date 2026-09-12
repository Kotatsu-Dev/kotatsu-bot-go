package cb

import (
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/cb/roulette"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
)

var START = Start()

var START_LINK = StartLink()

// TODO: otherwise
var MAIN = PrivateMessagesGuard(GetCurrentUser().
	Then(func(user *db.User_ReadJSON) Executor {
		return OneOf(
			TextGuard("keyboard.gender_male", SetGender(user, "male")),
			TextGuard("keyboard.gender_female", SetGender(user, "female")),
			TextGuard("keyboard.visited_enough", WasAtEvents(user, true)),
			TextGuard("keyboard.not_visited_enough", WasAtEvents(user, false)),
			TextGuard("keyboard.fill_back_later", WasntAtEvents(user, false)),
			TextGuard("keyboard.fill_now", WasntAtEvents(user, true)),
			TextGuard("keyboard.join_club", JoinClub(user)),
			TextGuard("keyboard.event_registration", SigningUpForActivity(user)),
			TextGuard("keyboard.to_main_menu", BackMainMenu(user)),
			TextGuard("keyboard.leave_club", LeaveClub(user)),
			TextGuard("keyboard.to_roulette_menu", roulette.Main(user)),
			TextGuard("keyboard.participate_roulette", roulette.Participate(user)),
			TextGuard("keyboard.leave_roulette", roulette.CancelParticipate(user)),
			TextGuard("keyboard.send_title", roulette.AnimeWish(user)),
			TextGuard("keyboard.roulette_rules", roulette.Rules()),
			TextGuard("keyboard.roulette_theme", roulette.MainTheme()),
			TextGuard("keyboard.roulette_list", roulette.LinkMyList(user)),
			TextGuard("keyboard.my_events", MyActivities(user)),
			TextGuard("keyboard.not_my_number", NoPhoneNumber(user)),

			StepGuard(user, config.STEP_ITMO_ENTER_ISU, ITMO_EnterISU(user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_ITMO_ENTER_ISU, ITMO_EnterISU(user, "activity")),
			StepGuard(user, config.STEP_ITMO_ENTER_FULLNAME, ITMO_EnterFullName(user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME, ITMO_EnterFullName(user, "activity")),
			StepGuard(user, config.STEP_NOITMO_ENTER_FULLNAME, NoITMO_EnterFullName(user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME, NoITMO_EnterFullName(user, "activity")),
			StepGuard(user, config.STEP_NOITMO_ENTER_PHONE, NoITMO_EnterPhoneNumber(user, "join_club")),
			StepGuard(user, config.STEP_CHANGING_PHONE, ChangePhoneNumber(user)),
			StepGuard(user, config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE, NoITMO_EnterPhoneNumber(user, "activity")),
			StepGuard(user, config.STEP_USER_LEAVES_CLUB, LeavesClub(user)),
			StepGuard(user, config.STEP_ANIME_RULETTE_ENTER_ENIGMATIC_TITLE, roulette.EnterEnigmaticTitle(user)),
			StepGuard(user, config.STEP_ANIME_RULETTE_ENTER_LINK_MY_ANIME_LIST, roulette.EnterLinkMyAnimeList(user)),

			SendMessageM("unknown_command", nil),
		)
	}).Otherwise(ProccessRegistration()))

// Ignore unknown commands - it's either unregistered, either unauthorized
var CALLBACK_MAIN = GetCurrentUser().
	Then(func(user *db.User_ReadJSON) Executor {
		return OneOf(
			QueryGuard("JOIN_CLUB", JoinClubQuery(user)),
			QueryGuard("APPOINTMENT", AppointQuery(user)),
			QueryGuard("ACTIVITIES", ActivitiesQuery(user)),
			QueryGuard("MY_ACTIVITIES", ActivitiesQuery(user)),
			QueryGuard("ACTIVITY_SUBSCRIBE", SubscribeQuery(user)),
			QueryGuard("ACTIVITY_UNSUBSCRIBE", UnsubscribeQuery(user)),
			QueryGuard("RELEVANC_PHONE", RelevancePhoneQuery(user)),
			QueryGuard("ROULETTES", roulette.MainQuery(user)),
		)
	})
