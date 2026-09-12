package cb

import (
	"context"
	"io"
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func MainMenuKeyboard(user *db.User_ReadJSON) models.ReplyMarkup {
	return ITE(
		user.IsClubMember,
		keyboards.Keyboard_MainMenuButtonsClubMember,
		keyboards.Keyboard_MainMenuButtonsDefault,
	)
}

func SendMainMenu(user *db.User_ReadJSON) Executor {
	return SendMessageM("main_menu", MainMenuKeyboard(user))
}

func Start() Executor {
	return GetCurrentUser().
		Then(func(user *db.User_ReadJSON) Executor {
			return WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
				return SendMessageMT("welcome", update.Message.From, MainMenuKeyboard(user))
			})
		}).
		Otherwise(ProccessRegistration())
}

func StartLink() Executor {
	return GetCurrentUser().
		Then(func(user *db.User_ReadJSON) Executor {
			return CommandArgUint("/start").
				Then(func(activity_id uint64) Executor {
					return ShowActivity(user, uint(activity_id))
				}).
				Otherwise(Func(func() (bool, error) {
					rr_debug.PrintLOG("CommandHandlers.go", "StartLink", "CommandArgUint", "Ошибка конвертации строки в uint", "")
					return true, nil
				}))
		}).
		Otherwise(ProccessRegistration())
}

func SetGender(user *db.User_ReadJSON, gender db.Gender) Executor {
	return Do(
		UpdateGender(user, gender),
		SendMainMenu(user),
	)
}

// ---

func WasAtEvents(current_user *db.User_ReadJSON, actually bool) Executor {
	return Do(
		UpdateVisited(current_user, actually),
		If(
			actually,
			SendMessageM(
				"request.is_itmo", keyboards.InlineKbd_JoinClub,
			),
			SendMessageM(
				"request.not_enough_visits", keyboards.Keyboard_WasntAtEvents,
			),
		),
	)
}

func WasntAtEvents(user *db.User_ReadJSON, cont bool) Executor {
	return If(
		cont,
		SendMessageM("request.is_itmo", keyboards.InlineKbd_JoinClub),
		SendMainMenu(user),
	)
}

func JoinClub(user *db.User_ReadJSON) Executor {
	if user.IsSentRequest {
		return Do(
			SendMessageM(
				"request.in_progress", nil,
			),
			SendMainMenu(user),
		)
	}
	if user.IsClubMember {
		return Do(
			SendMessageM(
				"request.already_accepted", nil,
			),
			SendMainMenu(user),
		)
	}
	return SendMessageM(
		"request.rules", keyboards.Keyboard_WasAtEvents,
	)
}

func SigningUpForActivity(user *db.User_ReadJSON) Executor {
	return GetActiveActivities().
		Then(func(activities []db.Activity_ReadJSON) Executor {
			return HasActiveRoulette().
				Then(func(has_roulette bool) Executor {
					text, keyboard := "events.empty", models.ReplyMarkup(nil)
					if len(activities) > 0 || has_roulette {
						text = "events.list"
						keyboard = keyboards.CreateInlineKbd_ActivitiesList(activities, user.UserTgID, has_roulette)
					}

					return OpenCalendar().
						Then(func(file io.Reader) Executor {
							return SendPhotoM(file, text, keyboard)
						}).
						Otherwise(SendMessageM(text, keyboard))
				})
		})
}

func BackMainMenu(user *db.User_ReadJSON) Executor {
	return Do(
		UpdateStep(user, config.STEP_DEFAULT),
		SendMainMenu(user),
	)
}

func LeaveClub(current_user *db.User_ReadJSON) Executor {
	return Do(
		UpdateStep(current_user, config.STEP_USER_LEAVES_CLUB),
		SendMessageM(
			"leave_reason", keyboards.Keyboard_Skip,
		),
	)
}

func MyActivities(user *db.User_ReadJSON) Executor {
	return GetUserActiveActivities(user).
		Then(func(activities []db.Activity_ReadJSON) Executor {
			return If(
				len(activities) == 0,
				SendMessageM(
					"my_events.empty", nil,
				),
				SendMessageM(
					"my_events.list", keyboards.CreateInlineKbd_MyActivitiesList(activities),
				),
			)
		})
}

func NoPhoneNumber(user *db.User_ReadJSON) Executor {
	return Do(
		UpdateStep(user, config.STEP_DEFAULT),
		SendMessageM(
			"request.no_phone_number", nil,
		),
		SendMainMenu(user),
	)
}

func ProccessRegistration() Executor {
	return FirstMatch(
		TextGuard("keyboard.continue", CreateOrGetUser().
			Then(func(user *db.User_ReadJSON) Executor {
				return SendMessageM("gender_select", keyboards.Keyboard_GenderSelect)
			}).
			Otherwise(SendMessageM("error.database", nil))),
		WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
			// TODO: Fix after merge with dev, avoiding conflict from editing locale
			return Do(
				SendDocumentM("CAACAgIAAx0CbgUG4QACCWpostfAVRPNDHNAWu8vcIbjv0nuagACrXQAAl8iQUmAFQIjshq4bTYE"),
				SendMessageRaw(
					update.Message.From.ID,
					config.T("hello")+"\n"+config.T("personal_data"),
					keyboards.Registration,
				),
			)
		}),
	)
}
