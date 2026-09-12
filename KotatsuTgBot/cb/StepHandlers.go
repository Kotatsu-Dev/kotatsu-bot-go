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

func ITMO_EnterISU(user *db.User_ReadJSON, action string) Executor {
	return ParseTextInt().
		Then(func(i int) Executor {
			return Do(
				UpdateUser(user, map[string]any{
					"isu": strconv.Itoa(i),
					"step": ITE(
						action == "join_club",
						config.STEP_ITMO_ENTER_FULLNAME,
						config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME,
					),
				}),
				SendMessageM(
					"request.enter_full_name", nil,
				),
			)
		}).
		Otherwise(SendMessageM(
			"request.not_isu_id", nil,
		))
}

func ITMO_EnterFullName(user *db.User_ReadJSON, action string) Executor {
	return MatchText(fullNameRegexp).
		Then(func(text string) Executor {
			return Do(
				UpdateUser(user, map[string]any{
					"full_name":      text,
					"step":           config.STEP_DEFAULT,
					"is_itmo":        true,
					"is_filled_data": true,
				}),
				If(
					action == "join_club",
					CreateRequest(user).
						Then(func(user *db.User_ReadJSON) Executor {
							return Do(
								SendMessageT(
									config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
									"request.notification", user,
									nil,
								),
								SendMessageM(
									"request.sent", nil,
								),
								SendMainMenu(user),
							)
						}).
						Otherwise(Do(
							SendMessageM(
								"error.generic", nil,
							),
							SendMainMenu(user),
						)),
					GetActivityByID(uint(user.TempActivityID)).
						Then(func(activity *db.Activity_ReadJSON) Executor {
							if activity.Status {
								// No ITMO check since we 100% from ITMO here
								return Do(
									AddParticipant(activity, user),
									SendMessageMT(
										"events.registered", activity,
										keyboards.ListEvents,
									),
								)
							} else {
								return SendMessageM(
									"events.non_existent", keyboards.ListEvents,
								)
							}
						}).
						Otherwise(SendMessageM(
							"events.non_existent", keyboards.ListEvents,
						)),
				),
			)
		}).
		Otherwise(SendMessageM(
			"request.incorrect_name_format", nil,
		))
}

func NoITMO_EnterFullName(user *db.User_ReadJSON, action string) Executor {
	return MatchText(fullNameRegexp).
		Then(func(text string) Executor {
			return Do(
				UpdateUser(user, map[string]any{
					"full_name": text,
					"step": ITE(
						action == "join_club",
						config.STEP_NOITMO_ENTER_PHONE,
						config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE,
					),
				}),
				SendMessageM(
					"request.enter_phone", keyboards.Keyboard_RequestContact,
				),
			)
		}).
		Otherwise(SendMessageM(
			"request.incorrect_name_format", nil,
		))
}

func NoITMO_EnterPhoneNumber(user *db.User_ReadJSON, action string) Executor {
	return GetContact().
		Then(func(contact string) Executor {
			return Do(
				UpdateUser(user, map[string]any{
					"phone_number":   contact,
					"step":           config.STEP_DEFAULT,
					"is_itmo":        false,
					"is_filled_data": true,
				}),
				If(
					action == "join_club",
					CreateRequest(user).
						Then(func(user *db.User_ReadJSON) Executor {
							return Do(
								SendMessageT(
									config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
									"request.notification", user,
									nil,
								),
								SendMessageM(
									"request.sent", nil,
								),
								SendMainMenu(user),
							)
						}).
						Otherwise(Do(
							SendMessageM(
								"error.generic", nil,
							),
							SendMainMenu(user),
						)),
					GetActivityByID(uint(user.TempActivityID)).
						Then(func(activity *db.Activity_ReadJSON) Executor {
							if !activity.Status {
								return SendMessageM(
									"events.non_existent",
									keyboards.ListEvents,
								)
							}
							// No itmo check since we 100% not from ITMO here
							if activity.GuestRegistrationUntil != nil &&
								activity.GuestRegistrationUntil.Before(time.Now()) {
								return SendMessageM(
									"events.registration_closed", keyboards.ListEvents,
								)
							} else {
								return Do(
									AddParticipant(activity, user),
									SendMessageMT(
										"events.registered", activity,
										keyboards.ListEvents,
									),
								)
							}
						}).
						Otherwise(SendMessageM(
							"events.non_existent", keyboards.ListEvents,
						)),
				),
			)
		}).
		Otherwise(SendMessageM(
			"request.incorrect_phone_format", nil,
		))
}

func ChangePhoneNumber(user *db.User_ReadJSON) Executor {
	return GetContact().
		Then(func(contact string) Executor {
			return Do(
				UpdateUser(user, map[string]any{
					"phone_number": contact,
					"step":         config.STEP_DEFAULT,
				}),
				GetActivityByID(uint(user.TempActivityID)).
					Then(func(activity *db.Activity_ReadJSON) Executor {
						if activity.Status {
							if activity.GuestRegistrationUntil != nil &&
								!user.IsITMO &&
								activity.GuestRegistrationUntil.Before(time.Now()) {
								return SendMessageM(
									"events.registration_closed", keyboards.ListEvents,
								)
							} else {
								return Do(
									AddParticipant(activity, user),
									SendMessageMT(
										"events.saved_n_registered", activity,
										keyboards.ListEvents,
									),
								)
							}
						} else {
							return SendMessageM(
								"events.non_existent", keyboards.ListEvents,
							)
						}
					}).
					Otherwise(SendMessageM(
						"events.non_existent", keyboards.ListEvents,
					)),
			)
		}).
		Otherwise(SendMessageM(
			"request.incorrect_phone_format", nil,
		))

}

func LeavesClub(user *db.User_ReadJSON) Executor {
	return Do(
		UpdateUser(user, map[string]any{
			"is_club_member":  false,
			"is_sent_request": false,
		}),
		WithCtx(func(ctx context.Context, b *bot.Bot, update *models.Update) Executor {
			return SendMessageT(
				config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
				"leave_notification", &map[string]any{
					"user":   user,
					"reason": ITE(update.Message.Text == config.T("keyboard.skip"), "", update.Message.Text),
				}, nil,
			)
		}),
		SendMessageM(
			"leave_response", nil,
		),
		SendMainMenu(user),
	)
}
