package cb

import (
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/cb/roulette"
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
			TextGuard("keyboard.to_main_menu", BackMainMenuE(user)),
			TextGuard("keyboard.not_my_number", NoPhoneNumberE(user)),
		)
	})
