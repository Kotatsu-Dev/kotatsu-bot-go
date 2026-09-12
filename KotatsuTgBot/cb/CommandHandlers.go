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

func SendMainMenuE(user *db.User_ReadJSON) Executor {
	return SendMessageME("main_menu", MainMenuKeyboard(user))
}

func StartE() Executor {
	return GetCurrentUserE().
		Then(func(user *db.User_ReadJSON) Executor {
			return WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
				return SendMessageMTE("welcome", update.Message.From, MainMenuKeyboard(user))
			})
		}).
		Otherwise(ProccessRegistrationE())
}

func StartLinkE() Executor {
	return GetCurrentUserE().
		Then(func(user *db.User_ReadJSON) Executor {
			return CommandArgUintE("/start").
				Then(func(activity_id uint64) Executor {
					return ShowActivityE(user, uint(activity_id))
				}).
				Otherwise(Func(func() (bool, error) {
					rr_debug.PrintLOG("CommandHandlers.go", "StartLinkE", "CommandArgUintE", "Ошибка конвертации строки в uint", "")
					return true, nil
				}))
		}).
		Otherwise(ProccessRegistrationE())
}

func SetGenderE(user *db.User_ReadJSON, gender db.Gender) Executor {
	return Seq(
		UpdateGenderE(user, gender),
		SendMainMenuE(user),
	)
}

// ---

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

func WasntAtEventsE(user *db.User_ReadJSON, cont bool) Executor {
	return ITE(
		cont,
		SendMessageME("request.is_itmo", keyboards.InlineKbd_JoinClub),
		SendMainMenuE(user),
	)
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

func BackMainMenuE(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateStepE(user, config.STEP_DEFAULT),
		SendMainMenuE(user),
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
			return Seq(
				SendDocumentME("CAACAgIAAx0CbgUG4QACCWpostfAVRPNDHNAWu8vcIbjv0nuagACrXQAAl8iQUmAFQIjshq4bTYE"),
				SendMessageRawE(
					update.Message.From.ID,
					config.T("hello")+"\n"+config.T("personal_data"),
					keyboards.Registration,
				),
			)
		}),
	)
}
