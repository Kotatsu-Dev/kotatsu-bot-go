package cb

import (
	"context"
	"regexp"
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var fullNameRegexp = regexp.MustCompile(`^([А-Яа-яЁё]+)\s(([А-Яа-яЁё]+)\s)?([А-Яа-яЁё]+)$`)

func ITMO_EnterISUE(user *db.User_ReadJSON, action string) Executor {
	return ParseTextInt().
		Then(func(i int) Executor {
			return Seq(
				UpdateUserE(user, map[string]any{
					"isu": strconv.Itoa(i),
					"step": ITE(
						action == "join_club",
						config.STEP_ITMO_ENTER_FULLNAME,
						config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME,
					),
				}),
				SendMessageME(
					"request.enter_full_name", nil,
				),
			)
		}).
		Otherwise(SendMessageME(
			"request.not_isu_id", nil,
		))
}

func ITMO_EnterFullNameE(user *db.User_ReadJSON, action string) Executor {
	return MatchText(fullNameRegexp).
		Then(func(text string) Executor {
			return Seq(
				UpdateUserE(user, map[string]any{
					"full_name":      text,
					"step":           config.STEP_DEFAULT,
					"is_itmo":        true,
					"is_filled_data": true,
				}),
				ITE(
					action == "join_club",
					CreateRequest(user).
						Then(func(user *db.User_ReadJSON) Executor {
							return Seq(
								SendMessageTE(
									config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
									"request.notification", user,
									nil,
								),
								SendMessageME(
									"request.sent", nil,
								),
								SendMainMenuE(user),
							)
						}).
						Otherwise(SendMessageME(
							"error.generic", nil,
						)).(Executor),
					GetActivityByID(uint(user.TempActivityID)).
						Then(func(activity *db.Activity_ReadJSON) Executor {
							if activity.Status {
								// No ITMO check since we 100% from ITMO here
								return Seq(
									AddParticipant(activity, user),
									SendMessageMTE(
										"events.registered", activity,
										keyboards.ListEvents,
									),
								)
							} else {
								return SendMessageME(
									"events.non_existent", keyboards.ListEvents,
								)
							}
						}).
						Otherwise(SendMessageME(
							"events.non_existent", keyboards.ListEvents,
						)).(Executor),
				),
			)
		}).
		Otherwise(SendMessageME(
			"request.incorrect_name_format", nil,
		))
}

func NoITMO_EnterFullNameE(user *db.User_ReadJSON, action string) Executor {
	return MatchText(fullNameRegexp).
		Then(func(text string) Executor {
			return Seq(
				UpdateUserE(user, map[string]any{
					"full_name": text,
					"step": ITE(
						action == "join_club",
						config.STEP_NOITMO_ENTER_PHONE,
						config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE,
					),
				}),
				SendMessageME(
					"request.enter_phone", keyboards.Keyboard_RequestContact,
				),
			)
		}).
		Otherwise(SendMessageME(
			"request.incorrect_name_format", nil,
		))
}

func NoITMO_EnterPhoneNumberE(user *db.User_ReadJSON, action string) Executor {
	return GetContact().
		Then(func(contact string) Executor {
			return Seq(
				UpdateUserE(user, map[string]any{
					"phone_number":   contact,
					"step":           config.STEP_DEFAULT,
					"is_itmo":        false,
					"is_filled_data": true,
				}),
				ITE(
					action == "join_club",
					Seq(
						UpdateUserE(user, map[string]any{
							"is_sent_request": true,
						}),
						CreateRequest(user).
							Then(func(user *db.User_ReadJSON) Executor {
								return Seq(
									SendMessageTE(
										config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
										"request.notification", user,
										nil,
									),
									SendMessageME(
										"request.sent", nil,
									),
									SendMainMenuE(user),
								)
							}).
							Otherwise(SendMessageME(
								"error.generic", nil,
							)),
					),
					GetActivityByID(uint(user.TempActivityID)).
						Then(func(activity *db.Activity_ReadJSON) Executor {
							if !activity.Status {
								return SendMessageME(
									"events.non_existent",
									keyboards.ListEvents,
								)
							}
							// No itmo check since we 100% not from ITMO here
							if activity.GuestRegistrationUntil != nil &&
								activity.GuestRegistrationUntil.Before(time.Now()) {
								return SendMessageME(
									"events.registration_closed", keyboards.ListEvents,
								)
							} else {
								return Seq(
									AddParticipant(activity, user),
									SendMessageMTE(
										"events.registered", activity,
										keyboards.ListEvents,
									),
								)
							}
						}).
						Otherwise(SendMessageME(
							"events.non_existent", keyboards.ListEvents,
						)).(Executor),
				),
			)
		}).
		Otherwise(SendMessageME(
			"request.incorrect_phone_format", nil,
		))
}

func ChangePhoneNumberE(user *db.User_ReadJSON) Executor {
	return GetContact().
		Then(func(contact string) Executor {
			return Seq(
				UpdateUserE(user, map[string]any{
					"phone_number": contact,
					"step":         config.STEP_DEFAULT,
				}),
				GetActivityByID(uint(user.TempActivityID)).
					Then(func(activity *db.Activity_ReadJSON) Executor {
						if activity.Status {
							if activity.GuestRegistrationUntil != nil &&
								!user.IsITMO &&
								activity.GuestRegistrationUntil.Before(time.Now()) {
								return SendMessageME(
									"events.registration_closed", keyboards.ListEvents,
								)
							} else {
								return Seq(
									AddParticipant(activity, user),
									SendMessageME(
										"events.saved_n_registered", keyboards.ListEvents,
									),
								)
							}
						} else {
							return SendMessageME(
								"events.non_existent", keyboards.ListEvents,
							)
						}
					}).
					Otherwise(SendMessageME(
						"events.non_existent", keyboards.ListEvents,
					)),
			)
		}).
		Otherwise(SendMessageME(
			"request.incorrect_phone_format", nil,
		))

}

func LeavesClubE(user *db.User_ReadJSON) Executor {
	return Seq(
		UpdateUserE(user, map[string]any{
			"is_club_member":  false,
			"is_sent_request": false,
		}),
		WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
			return SendMessageTE(
				config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
				"leave_notification", &map[string]any{
					"user":   user,
					"reason": ITE(update.Message.Text == config.T("keyboard.skip"), "", update.Message.Text),
				}, nil,
			)
		}),
		SendMessageME(
			"leave_response", nil,
		),
		SendMainMenuE(user),
	)
}
