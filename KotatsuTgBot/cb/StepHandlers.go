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

func ITMO_EnterISU(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	if _, err := strconv.Atoi(update.Message.Text); err == nil {
		UpdateCurrentUser(current_user, map[string]any{
			"user_tg_id": current_user.UserTgID,
			"isu":        update.Message.Text,
			"step": ITE(
				action == "join_club",
				config.STEP_ITMO_ENTER_FULLNAME,
				config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME,
			),
		})

		SendMessageM(
			ctx, b, update,
			"request.enter_full_name", nil,
		)
	} else {
		SendMessageM(
			ctx, b, update,
			"request.not_isu_id", nil,
		)
	}
}

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

func ITMO_EnterFullName(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	matched := fullNameRegexp.MatchString(update.Message.Text)

	if !matched {
		SendMessageM(
			ctx, b, update,
			"request.incorrect_name_format", nil,
		)
		return
	}

	UpdateCurrentUser(current_user, map[string]any{
		"full_name":      update.Message.Text,
		"step":           config.STEP_DEFAULT,
		"is_itmo":        true,
		"is_filled_data": true,
	})
	if action == "join_club" {
		_, updated_user, _ := UpdateCurrentUser(current_user, map[string]any{
			"is_sent_request": true,
		})

		db_answer_code := db.DB_CREATE_Request(current_user.ID)
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			SendMessageT(
				ctx, b,
				config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
				"request.notification", updated_user,
				nil,
			)
			SendMessageM(
				ctx, b, update,
				"request.sent", nil,
			)
			SendMainMenu(ctx, current_user, b, update)

		default:
			SendMessageM(
				ctx, b, update,
				"error.generic", nil,
			)
			return
		}

	} else {
		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))

		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			if activity.Status {
				// No ITMO check since we 100% from ITMO here
				db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

				SendMessageMT(
					ctx, b, update,
					"events.registered", activity,
					keyboards.ListEvents,
				)
			} else {
				SendMessageM(
					ctx, b, update,
					"events.non_existent", keyboards.ListEvents,
				)
			}
		default:
			SendMessageM(
				ctx, b, update,
				"events.non_existent", keyboards.ListEvents,
			)
		}
	}
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
									"request.notification", user, // TODO: Updated user
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
							return Seq(
								AddParticipant(activity, user),
								SendMessageMTE(
									"events.registered", activity,
									keyboards.ListEvents,
								),
							)
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

func NoITMO_EnterFullName(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	matched := fullNameRegexp.MatchString(update.Message.Text)

	if !matched {
		SendMessageM(
			ctx, b, update,
			"request.incorrect_name_format", nil,
		)
		return
	}

	UpdateCurrentUser(current_user, map[string]any{
		"full_name": update.Message.Text,
		"step": ITE(
			action == "join_club",
			config.STEP_NOITMO_ENTER_PHONE,
			config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE,
		),
	})

	SendMessageM(
		ctx, b, update,
		"request.enter_phone", keyboards.Keyboard_RequestContact,
	)
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

func NoITMO_EnterPhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	if update.Message.Contact != nil {
		UpdateCurrentUser(current_user, map[string]any{
			"phone_number":   update.Message.Contact.PhoneNumber,
			"step":           config.STEP_DEFAULT,
			"is_itmo":        false,
			"is_filled_data": true,
		})

		if action == "join_club" {
			_, updated_user, _ := UpdateCurrentUser(current_user, map[string]any{
				"is_sent_request": true,
			})

			db_answer_code := db.DB_CREATE_Request(current_user.ID)
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:
				SendMessageT(
					ctx, b,
					config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
					"request.notification", updated_user,
					nil,
				)
				SendMessageM(
					ctx, b, update,
					"request.sent", nil,
				)
				SendMainMenu(ctx, current_user, b, update)
			default:
				SendMessageM(
					ctx, b, update,
					"error.generic", nil,
				)
				return
			}
		} else {
			db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:
				if activity.Status {
					// No itmo check since we 100% not from ITMO here
					if activity.GuestRegistrationUntil != nil &&
						activity.GuestRegistrationUntil.Before(time.Now()) {
						SendMessageM(
							ctx, b, update,
							"events.registration_closed", keyboards.ListEvents,
						)
					} else {
						db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)
						SendMessageMT(
							ctx, b, update,
							"events.registered", activity,
							keyboards.ListEvents,
						)
					}
				} else {
					SendMessageMT(
						ctx, b, update,
						"events.non_existent", activity,
						keyboards.ListEvents,
					)
				}
			}
		}

	} else {
		SendMessageM(
			ctx, b, update,
			"request.incorrect_phone_format", nil,
		)
	}
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
										"request.notification", user, // TODO: Updated user
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
							// No itmo check since we 100% not from ITMO here
							if activity.GuestRegistrationUntil != nil &&
								activity.GuestRegistrationUntil.Before(time.Now()) {
								return SendMessageME(
									"events.registration_closed", keyboards.ListEvents,
								)
							} else {
								db.DB_UPDATE_Activity_ADD_Participants(activity.ID, user.ID)
								return SendMessageMTE(
									"events.registered", activity,
									keyboards.ListEvents,
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

func ChangePhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	if update.Message.Contact != nil {
		UpdateCurrentUser(current_user, map[string]any{
			"user_tg_id":   update.Message.From.ID,
			"phone_number": update.Message.Contact.PhoneNumber,
			"step":         config.STEP_DEFAULT,
		})

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			if activity.Status {
				if activity.GuestRegistrationUntil != nil &&
					!current_user.IsITMO &&
					activity.GuestRegistrationUntil.Before(time.Now()) {
					SendMessageM(
						ctx, b, update,
						"events.registration_closed", keyboards.ListEvents,
					)
				} else {
					db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

					SendMessageM(
						ctx, b, update,
						"events.saved_n_registered", keyboards.ListEvents,
					)
				}
			} else {
				SendMessageMT(
					ctx, b, update,
					"events.non_existent", activity,
					keyboards.ListEvents,
				)
			}
		}

	} else {
		SendMessageM(
			ctx, b, update,
			"request.incorrect_phone_format", nil,
		)
	}
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
								db.DB_UPDATE_Activity_ADD_Participants(activity.ID, user.ID)

								return SendMessageME(
									"events.saved_n_registered", keyboards.ListEvents,
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

func LeavesClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	UpdateCurrentUser(current_user, map[string]any{
		"is_club_member":  false,
		"is_sent_request": false,
	})

	SendMessageT(
		ctx, b, config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
		"leave_notification", &map[string]any{
			"user":   current_user,
			"reason": ITE(update.Message.Text == config.T("keyboard.skip"), "", update.Message.Text),
		}, nil,
	)
	SendMessageM(
		ctx, b, update,
		"leave_response", nil,
	)
	SendMainMenu(ctx, current_user, b, update)
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
