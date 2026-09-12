package cb

import (
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/cb/roulette"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
)

var START = Start()

var START_LINK = StartLink()

var MAIN = PrivateMessagesGuard(GetCurrentUser().
	Then(func(user *db.User_ReadJSON) Executor {
		return FirstMatch(
			TextGuard("keyboard.gender_male", Lazy2(SetGender, user, "male")),
			TextGuard("keyboard.gender_female", Lazy2(SetGender, user, "female")),
			TextGuard("keyboard.visited_enough", Lazy2(WasAtEvents, user, true)),
			TextGuard("keyboard.not_visited_enough", Lazy2(WasAtEvents, user, false)),
			TextGuard("keyboard.fill_back_later", Lazy2(WasntAtEvents, user, false)),
			TextGuard("keyboard.fill_now", Lazy2(WasntAtEvents, user, true)),
			TextGuard("keyboard.join_club", Lazy(JoinClub, user)),
			TextGuard("keyboard.event_registration", Lazy(SigningUpForActivity, user)),
			TextGuard("keyboard.to_main_menu", Lazy(BackMainMenu, user)),
			TextGuard("keyboard.leave_club", Lazy(LeaveClub, user)),
			TextGuard("keyboard.to_roulette_menu", Lazy(roulette.Main, user)),
			TextGuard("keyboard.participate_roulette", Lazy(roulette.Participate, user)),
			TextGuard("keyboard.leave_roulette", Lazy(roulette.CancelParticipate, user)),
			TextGuard("keyboard.send_title", Lazy(roulette.AnimeWish, user)),
			TextGuard("keyboard.roulette_rules", Lazy0(roulette.Rules)),
			TextGuard("keyboard.roulette_theme", Lazy0(roulette.MainTheme)),
			TextGuard("keyboard.roulette_list", Lazy(roulette.LinkMyList, user)),
			TextGuard("keyboard.my_events", Lazy(MyActivities, user)),
			TextGuard("keyboard.not_my_number", Lazy(NoPhoneNumber, user)),

			StepGuard(user, config.STEP_ITMO_ENTER_ISU, Lazy2(ITMO_EnterISU, user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_ITMO_ENTER_ISU, Lazy2(ITMO_EnterISU, user, "activity")),
			StepGuard(user, config.STEP_ITMO_ENTER_FULLNAME, Lazy2(ITMO_EnterFullName, user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME, Lazy2(ITMO_EnterFullName, user, "activity")),
			StepGuard(user, config.STEP_NOITMO_ENTER_FULLNAME, Lazy2(NoITMO_EnterFullName, user, "join_club")),
			StepGuard(user, config.STEP_APPOINTMENT_NOITMO_ENTER_FULLNAME, Lazy2(NoITMO_EnterFullName, user, "activity")),
			StepGuard(user, config.STEP_NOITMO_ENTER_PHONE, Lazy2(NoITMO_EnterPhoneNumber, user, "join_club")),
			StepGuard(user, config.STEP_CHANGING_PHONE, Lazy(ChangePhoneNumber, user)),
			StepGuard(user, config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE, Lazy2(NoITMO_EnterPhoneNumber, user, "activity")),
			StepGuard(user, config.STEP_USER_LEAVES_CLUB, Lazy(LeavesClub, user)),
			StepGuard(user, config.STEP_ANIME_RULETTE_ENTER_ENIGMATIC_TITLE, Lazy(roulette.EnterEnigmaticTitle, user)),
			StepGuard(user, config.STEP_ANIME_RULETTE_ENTER_LINK_MY_ANIME_LIST, Lazy(roulette.EnterLinkMyAnimeList, user)),

			SendMessageM("unknown_command", nil),
		)
	}).
	Otherwise(ProccessRegistration()))

// Ignore unknown commands - it's either unregistered, either unauthorized
var CALLBACK_MAIN = GetCurrentUser().Then(dispatchQuery)

func dispatchQuery(user *db.User_ReadJSON) Executor {
	return FirstMatch(
		QueryGuard("JOIN_CLUB", Lazy(JoinClubQuery, user)),
		QueryGuard("APPOINTMENT", Lazy(AppointQuery, user)),
		QueryGuard("ACTIVITIES", Lazy(ActivitiesQuery, user)),
		QueryGuard("MY_ACTIVITIES", Lazy(ActivitiesQuery, user)),
		QueryGuard("ACTIVITY_SUBSCRIBE", Lazy(SubscribeQuery, user)),
		QueryGuard("ACTIVITY_UNSUBSCRIBE", Lazy(UnsubscribeQuery, user)),
		QueryGuard("RELEVANC_PHONE", Lazy(RelevancePhoneQuery, user)),
		QueryGuard("ROULETTES", Lazy(roulette.MainQuery, user)),
	)
}
