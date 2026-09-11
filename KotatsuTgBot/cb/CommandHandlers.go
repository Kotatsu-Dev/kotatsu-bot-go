package cb

import (
	"context"
	"io"
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendMainMenu(ctx context.Context, current_user *db.User_ReadJSON, b *bot.Bot, update *models.Update) {
	if current_user.IsClubMember {
		SendMessageM(
			ctx, b, update,
			"main_menu", keyboards.Keyboard_MainMenuButtonsClubMember,
		)
	} else {
		SendMessageM(
			ctx, b, update,
			"main_menu", keyboards.Keyboard_MainMenuButtonsDefault,
		)
	}
}

func SendMainMenuE(user *db.User_ReadJSON) Executor {
	return SendMessageME(
		"main_menu",
		ITE(
			user.IsClubMember,
			keyboards.Keyboard_MainMenuButtonsClubMember,
			keyboards.Keyboard_MainMenuButtonsDefault,
		),
	)
}

func SetGender(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, gender db.Gender) {
	UpdateGender(current_user, gender)
	SendMainMenu(ctx, current_user, b, update)
}

func SetGenderE(user *db.User_ReadJSON, gender db.Gender) Executor {
	return Seq(
		UpdateGenderE(user, gender),
		SendMainMenuE(user),
	)
}

// ---

func WasAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, actually bool) {
	UpdateVisited(current_user, actually)

	if actually {
		SendMessageM(
			ctx, b, update,
			"request.is_itmo", keyboards.InlineKbd_JoinClub,
		)
	} else {
		SendMessageM(
			ctx, b, update,
			"request.not_enough_visits", keyboards.Keyboard_WasntAtEvents,
		)
	}

}

func WasAtEventsE(current_user *db.User_ReadJSON, actually bool) Executor {
	return Seq(
		UpdateVisitedE(current_user, actually),
		ITE(
			actually,
			SendMessageME(
				"request.is_itmo", keyboards.InlineKbd_JoinClub,
			),
			SendMessageME(
				"request.not_enough_visits", keyboards.Keyboard_WasntAtEvents,
			),
		),
	)
}

func WasntAtEvents(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, cont bool) {
	if cont {
		SendMessageM(
			ctx, b, update,
			"request.is_itmo", keyboards.InlineKbd_JoinClub,
		)
	} else {
		SendMainMenu(ctx, current_user, b, update)
	}
}

func WasntAtEventsE(user *db.User_ReadJSON, cont bool) Executor {
	return ITE(
		cont,
		SendMessageME("request.is_itmo", keyboards.InlineKbd_JoinClub),
		SendMainMenuE(user),
	)
}

func JoinClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	if current_user.IsSentRequest {
		SendMessageM(
			ctx, b, update,
			"request.in_progress", nil,
		)
		SendMainMenu(ctx, current_user, b, update)
	} else if current_user.IsClubMember {
		SendMessageM(
			ctx, b, update,
			"request.already_accepted", nil,
		)
		SendMainMenu(ctx, current_user, b, update)
	} else {
		SendMessageM(
			ctx, b, update,
			"request.rules", keyboards.Keyboard_WasAtEvents,
		)
	}
}

func JoinClubE(user *db.User_ReadJSON) Executor {
	if user.IsSentRequest {
		return Seq(
			SendMessageME(
				"request.in_progress", nil,
			),
			SendMainMenuE(user),
		)
	}
	if user.IsClubMember {
		return Seq(
			SendMessageME(
				"request.already_accepted", nil,
			),
			SendMainMenuE(user),
		)
	}
	return SendMessageME(
		"request.rules", keyboards.Keyboard_WasAtEvents,
	)
}

func SigningUpForActivity(ctx context.Context, b *bot.Bot, update *models.Update) {
	activities_list := db.DB_GET_Active_Activities()

	status, _ := db.DB_GET_AnimeRoulette_BY_Status(true)
	has_roulette := status == db.DB_ANSWER_SUCCESS

	var text string
	var keyboard models.ReplyMarkup
	if len(activities_list) > 0 || has_roulette {
		text = "events.list"
		keyboard = keyboards.CreateInlineKbd_ActivitiesList(activities_list, update.Message.From.ID, has_roulette)
	} else {
		text = "events.empty"
	}

	file, name := OpenCalendar()
	if file != nil {
		SendPhotoM(
			ctx, b, update,
			name, file,
			text, keyboard,
		)
		file.Close()
	} else {
		SendMessageM(
			ctx, b, update,
			text, keyboard,
		)
	}
}

func SigningUpForActivityE(user *db.User_ReadJSON) Executor {
	return GetActiveActivities().
		Then(func(activities []db.Activity_ReadJSON) Executor {
			return HasActiveRoulette().
				Then(func(has_roulette bool) Executor {
					text, keyboard := "events.empty", models.ReplyMarkup(nil)
					if len(activities) > 0 || has_roulette {
						text = "events.list"
						keyboard = keyboards.CreateInlineKbd_ActivitiesList(activities, user.UserTgID, has_roulette)
					}

					return OpenCalendarE().
						Then(func(file io.Reader) Executor {
							return SendPhotoME(file, text, keyboard)
						}).
						Otherwise(SendMessageME(text, keyboard))
				})
		})
}

func BackMainMenu(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_DEFAULT)
	SendMainMenu(ctx, current_user, b, update)
}

func BackMainMenuE(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStepE(user, config.STEP_DEFAULT),
		SendMainMenuE(user),
	)
}

func LeaveClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_USER_LEAVES_CLUB)
	SendMessageM(
		ctx, b, update,
		"leave_reason", keyboards.Keyboard_Skip,
	)
}

func LeaveClubE(current_user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStepE(current_user, config.STEP_USER_LEAVES_CLUB),
		SendMessageME(
			"leave_reason", keyboards.Keyboard_Skip,
		),
	)
}

func MyActivities(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	activities := db.DB_GET_User_Active_Activities(current_user.ID)
	if len(activities) == 0 {
		SendMessageM(
			ctx, b, update,
			"my_events.empty", nil,
		)
	} else {
		SendMessageM(
			ctx, b, update,
			"my_events.list", keyboards.CreateInlineKbd_MyActivitiesList(activities),
		)
	}
}

func MyActivitiesE(user *db.User_ReadJSON) Executor {
	return GetUserActiveActivities(user).
		Then(func(activities []db.Activity_ReadJSON) Executor {
			return ITE(
				len(activities) == 0,
				SendMessageME(
					"my_events.empty", nil,
				),
				SendMessageME(
					"my_events.list", keyboards.CreateInlineKbd_MyActivitiesList(activities),
				),
			)
		})
}

func NoPhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateStep(current_user, config.STEP_DEFAULT)
	SendMessageM(
		ctx, b, update,
		"request.no_phone_number", nil,
	)
	SendMainMenu(ctx, current_user, b, update)
}

func NoPhoneNumberE(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStepE(user, config.STEP_DEFAULT),
		SendMessageME(
			"request.no_phone_number", nil,
		),
		SendMainMenuE(user),
	)
}

func ProccessRegistrationE() Executor {
	return OneOf(
		TextGuard("keyboard.continue", CreateOrGetUserE().
			Then(func(user *db.User_ReadJSON) Executor {
				return SendMessageME("gender_select", keyboards.Keyboard_GenderSelect)
			}).
			Otherwise(SendMessageME("error.database", nil))),
		WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
			// TODO: Fix after merge with dev, avoiding conflict from editing locale
			return SendMessageRawE(
				update.Message.From.ID,
				config.T("hello")+"\n"+config.T("personal_data"),
				keyboards.Registration,
			)
		}),
	)
}
