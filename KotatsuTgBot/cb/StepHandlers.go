package cb

import (
	"context"
	"regexp"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/keyboards"
	"rr/kotatsutgbot/rr_debug"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var fullNameRegexp = regexp.MustCompile(`^([А-Яа-яЁё]+)\s(([А-Яа-яЁё]+)\s)?([А-Яа-яЁё]+)$`)

func ITMO_EnterISU(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	var text string
	var step int
	if action == "join_club" {
		step = config.STEP_ITMO_ENTER_FULLNAME
	} else {
		step = config.STEP_APPOINTMENT_ITMO_ENTER_FULLNAME
	}
	if _, err := strconv.Atoi(update.Message.Text); err == nil {
		db.DB_UPDATE_User(map[string]any{
			"user_tg_id": current_user.UserTgID,
			"isu":        update.Message.Text,
			"step":       step,
		})

		text = config.T("request.enter_full_name")
	} else {
		text = config.T("request.not_isu_id")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.From.ID,
		ParseMode: models.ParseModeHTML,
		Text:      text,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterISU", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func ITMO_EnterFullName(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	var text string
	var keyboard models.ReplyMarkup
	matched := fullNameRegexp.MatchString(update.Message.Text)

	if !matched {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.From.ID,
			ParseMode: models.ParseModeHTML,
			Text:      config.T("request.incorrect_name_format"),
		})
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
		}
		return
	}

	if action == "join_club" {
		update_user_data := map[string]any{
			"user_tg_id":      update.Message.From.ID,
			"step":            config.STEP_DEFAULT,
			"is_sent_request": true,
			"is_filled_data":  true,
			"is_itmo":         true,
		}

		_, updated_user, _ := db.DB_UPDATE_User(update_user_data)

		db_answer_code := db.DB_CREATE_Request(current_user.ID)
		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			text = config.T("request.sent")

			_, err_msg := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
				ParseMode: models.ParseModeHTML,
				Text:      config.TT("request.notification", updated_user),
			})
			if err_msg != nil {
				rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err_msg.Error())
			}

		default:
			text = config.T("error.generic")
		}

		keyboard = keyboards.Keyboard_MainMenuButtonsDefault

	} else {
		update_user_data := map[string]any{
			"user_tg_id":     update.Message.From.ID,
			"full_name":      update.Message.Text,
			"step":           config.STEP_DEFAULT,
			"is_itmo":        true,
			"is_filled_data": true,
		}
		db.DB_UPDATE_User(update_user_data)

		db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))

		switch db_answer_code {
		case db.DB_ANSWER_SUCCESS:
			if activity.Status {
				// No ITMO check since we 100% from ITMO here
				db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

				text = config.TT("events.registered", activity)
				keyboard = keyboards.ListEvents
			} else {
				text = config.TT("events.non_existent", activity)
				keyboard = keyboards.ListEvents
			}
		}
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func NoITMO_EnterFullName(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	matched := fullNameRegexp.MatchString(update.Message.Text)

	if !matched {
		_, err_msg := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.From.ID,
			ParseMode: models.ParseModeHTML,
			Text:      config.T("request.incorrect_name_format"),
		})
		if err_msg != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_ITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err_msg.Error())
		}
		return
	}

	var step int
	if action == "join_club" {
		step = config.STEP_NOITMO_ENTER_PHONE
	} else {
		step = config.STEP_APPOINTMENT_NOITMO_ENTER_PHONE
	}

	db.DB_UPDATE_User(map[string]any{
		"user_tg_id": current_user.UserTgID,
		"full_name":  update.Message.Text,
		"step":       step,
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        config.T("request.enter_phone"),
		ReplyMarkup: keyboards.Keyboard_RequestContact,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_NoITMO_EnterFullName", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

// TODO
func NoITMO_EnterPhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, action string) {
	var text string
	var keyboard models.ReplyMarkup
	if update.Message.Contact != nil {
		keyboard = keyboards.Keyboard_ToMainMenu

		if action == "join_club" {
			_, updated_user, _ := db.DB_UPDATE_User(map[string]any{
				"user_tg_id":      update.Message.From.ID,
				"step":            config.STEP_DEFAULT,
				"is_sent_request": true,
				"is_filled_data":  true,
				"is_itmo":         false,
			})

			db_answer_code := db.DB_CREATE_Request(current_user.ID)
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:
				text = config.T("request.sent")

				_, err := b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
					ParseMode: models.ParseModeHTML,
					Text:      config.TT("request.notification", updated_user),
				})
				if err != nil {
					rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err.Error())
				}

			default:
				text = config.T("error.generic")
			}

			keyboard = keyboards.Keyboard_MainMenuButtonsDefault
		} else {
			db.DB_UPDATE_User(map[string]any{
				"user_tg_id":     update.Message.From.ID,
				"phone_number":   update.Message.Contact.PhoneNumber,
				"step":           config.STEP_DEFAULT,
				"is_itmo":        false,
				"is_filled_data": true,
			})

			db_answer_code, activity := db.DB_GET_Activity_BY_ID(uint(current_user.TempActivityID))
			switch db_answer_code {
			case db.DB_ANSWER_SUCCESS:
				if activity.Status {
					// No itmo check since we 100% not from ITMO here
					if activity.GuestRegistrationUntil != nil &&
						activity.GuestRegistrationUntil.Before(time.Now()) {
						text = config.T("events.registration_closed")
						keyboard = keyboards.ListEvents
					} else {
						db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)
						text = config.TT("events.registered", activity)
						keyboard = keyboards.ListEvents
					}
				} else {
					text = config.TT("events.non_existent", activity)
					keyboard = keyboards.ListEvents
				}
			}
		}

	} else {
		text = config.T("request.incorrect_phone_format")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_NoITMO_EnterPhoneNumber", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func EnterSecretCode(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON, status string) {
	var text string

	_, updated_user, _ := db.DB_UPDATE_User(map[string]any{
		"user_tg_id":      update.Message.From.ID,
		"secret_code":     update.Message.Text,
		"step":            config.STEP_DEFAULT,
		"is_itmo":         status == "itmo",
		"is_sent_request": true,
		"is_filled_data":  true,
	})

	db_answer_code := db.DB_CREATE_Request(current_user.ID)
	switch db_answer_code {
	case db.DB_ANSWER_SUCCESS:
		text = config.T("request.sent")

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
			ParseMode: models.ParseModeHTML,
			Text:      config.TT("request.notification", updated_user),
		})
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err.Error())
		}

	default:
		text = config.T("error.generic")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboards.Keyboard_MainMenuButtonsDefault,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_EnterSecretCode", "b.SendMessage(ctx, params_user)", "Ошибка отправки сообщения", err.Error())
	}
}

func ChangePhoneNumber(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup
	if update.Message.Contact != nil {
		db.DB_UPDATE_User(map[string]any{
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
					text = config.T("events.registration_closed")
					keyboard = keyboards.ListEvents
				} else {
					db.DB_UPDATE_Activity_ADD_Participants(activity.ID, current_user.ID)

					text = config.TT("events.saved_n_registered", activity)
					keyboard = keyboards.ListEvents
				}
			} else {
				text = config.TT("events.non_existent", activity)
				keyboard = keyboards.ListEvents
			}
		}

	} else {
		text = config.T("request.incorrect_phone_format")
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_ChangePhoneNumber", "b.SendMessage", "Ошибка отправки сообщения", err.Error())
	}
}

func LeavesClub(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	var text string
	var keyboard models.ReplyMarkup

	switch update.Message.Text {
	case "Пропустить":
		update_user_data := map[string]any{
			"user_tg_id":      update.Message.From.ID,
			"is_club_member":  false,
			"is_sent_request": false,
		}

		text = config.T("leave_response")
		keyboard = keyboards.Keyboard_MainMenuButtonsDefault

		db.DB_UPDATE_User(update_user_data)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
			ParseMode: models.ParseModeHTML,
			LinkPreviewOptions: &models.LinkPreviewOptions{
				IsDisabled: bot.True(),
			},
			Text: config.TT("leave_notification", &map[string]any{
				"user":   current_user,
				"reason": "",
			}),
		})
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_LeavesClub", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err.Error())
		}

	default:
		update_user_data := map[string]any{
			"user_tg_id":      update.Message.From.ID,
			"is_club_member":  false,
			"is_sent_request": false,
		}

		text = config.T("leave_response")
		keyboard = keyboards.Keyboard_MainMenuButtonsDefault

		db.DB_UPDATE_User(update_user_data)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    config.GetConfig().CONFIG_ID_CHAT_SUPPORT,
			ParseMode: models.ParseModeHTML,
			LinkPreviewOptions: &models.LinkPreviewOptions{
				IsDisabled: bot.True(),
			},
			Text: config.TT("leave_notification", &map[string]any{
				"user":   current_user,
				"reason": update.Message.Text,
			}),
		})
		if err != nil {
			rr_debug.PrintLOG("botHandlers.go", "proccessStep_LeavesClub", "b.SendMessage(ctx, params_support)", "Ошибка отправки сообщения", err.Error())
		}
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.From.ID,
		ParseMode:   models.ParseModeHTML,
		Text:        text,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		rr_debug.PrintLOG("botHandlers.go", "proccessStep_LeavesClub", "b.SendMessage(ctx, params_user)", "Ошибка отправки сообщения", err.Error())
	}
}
