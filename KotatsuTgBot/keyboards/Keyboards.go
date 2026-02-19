// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

package keyboards

import (
	"rr/kotatsutgbot/db"
	"time"

	"github.com/go-telegram/bot/models"

	"fmt"
)

var Registration = Default().
	TextT("keyboard.continue").
	Build()

var Keyboard_GenderSelect = Default().
	TextT("keyboard.gender_male").
	TextT("keyboard.gender_female").
	Build()

var Keyboard_WasAtEvents = Default().
	TextT("keyboard.visited_enough").
	TextT("keyboard.not_visited_enough").
	Build()

var Keyboard_WasntAtEvents = Default().
	TextT("keyboard.fill_back_later").
	TextT("keyboard.fill_now").
	Build()

var Keyboard_MainMenuButtonsDefault = Default().
	TextT("keyboard.join_club").
	TextT("keyboard.event_registration").
	Build()

var Keyboard_MainMenuButtonsClubMember = Default().
	TextT("keyboard.event_registration").
	Row().
	TextT("keyboard.leave_club").
	TextT("keyboard.my_events").
	Build()

var CommunicationManager = Default().
	TextT("keyboard.to_main_menu").
	Build()

var ListEvents = Default().
	TextT("keyboard.to_main_menu").
	Build()

// Probably unused
var SelectedEvent = Default().
	TextT("keyboard.cancel_registration").
	TextT("keyboard.to_events_list")

// Клавиатура для аниме рулетки перед её запуском
func CreateKeyboard_AnimeRouletteStart(is_member bool) *models.ReplyKeyboardMarkup {
	return Default().
		TextTIf(is_member, "keyboard.leave_roulette", "keyboard.participate_roulette").
		Row().
		TextT("keyboard.roulette_rules").
		TextTC("keyboard.roulette_list", is_member).
		Row().
		TextT("keyboard.to_main_menu").
		Build()
}

// Клавиатура для аниме рулетки когда запущена
func CreateKeyboard_AnimeRouletteMenu(is_member bool) *models.ReplyKeyboardMarkup {
	return Default().
		TextTIf(is_member, "keyboard.send_title", "keyboard.participate_roulette").
		Row().
		TextT("keyboard.roulette_rules").
		TextT("keyboard.roulette_theme").
		TextTC("keyboard.roulette_list", is_member).
		Row().
		TextT("keyboard.to_main_menu").
		Build()
}

var Keyboard_Skip = Default().
	TextT("keyboard.skip").
	Row().
	TextT("keyboard.to_main_menu").
	Build()

// Probably unused
var Keyboard_CancelAnimeRoulette = Default().
	TextT("keyboard.to_roulette_menu").
	Build()

var Keyboard_ToMainMenu = Default().
	TextT("keyboard.to_main_menu").
	Build()

// For date formatting
var date_format = "02.01 15:04"
var loc, _ = time.LoadLocation("Europe/Moscow")

func CreateInlineKbd_ActivitiesList(activities []db.Activity_ReadJSON, user_tg_id int64, has_roulette bool) *models.InlineKeyboardMarkup {
	k := DefaultInline()

	for _, activity := range activities {

		is_participant := false

		for _, participant := range activity.Participants {
			if participant.UserTgID == user_tg_id {
				is_participant = true
				break
			}
		}

		title_format := "[%s] %s"
		if is_participant {
			title_format = "✅ [%s] %s"
		}

		k.
			Row().
			Data(
				fmt.Sprintf(
					title_format,
					activity.DateMeeting.In(loc).Format(date_format),
					activity.Title,
				),
				fmt.Sprintf("ACTIVITIES::%d", activity.ID),
			)
	}

	if has_roulette {
		k.
			Row().
			DataT("roulette.title", "ROULETTES")
	}

	return k.Build()
}

var InlineKbd_PartnersList = DefaultInline().Row().Build()

func CreateInlineKbd_MyActivitiesList(my_activities []*db.Activity) *models.InlineKeyboardMarkup {
	k := DefaultInline()

	for _, activity := range my_activities {
		k.
			Row().
			Data(
				fmt.Sprintf(
					"[%s] %s",
					activity.DateMeeting.In(loc).Format(date_format),
					activity.Title,
				),
				fmt.Sprintf("MY_ACTIVITIES::%d", activity.ID),
			)
	}

	return k.Build()
}

func CreateInlineKbd_SubscribeActivity(activity_id int) *models.InlineKeyboardMarkup {
	return DefaultInline().
		DataT(
			"keyboard.register_event",
			fmt.Sprintf("ACTIVITY_SUBSCRIBE::%d", activity_id)).
		Build()
}

func CreateInlineKbd_UnsubscribeActivity(activity_id int) *models.InlineKeyboardMarkup {
	return DefaultInline().
		DataT(
			"keyboard.cancel_registration",
			fmt.Sprintf("ACTIVITY_UNSUBSCRIBE::%d", activity_id)).
		Build()
}

var InlineKbd_JoinClub = DefaultInline().
	DataT("keyboard.is_student", "JOIN_CLUB::from_ITMO_student").
	Row().
	DataT("keyboard.is_graduate", "JOIN_CLUB::from_ITMO_graduate").
	Row().
	DataT("keyboard.is_employee", "JOIN_CLUB::from_ITMO_employee").
	Row().
	DataT("keyboard.is_student_employee", "JOIN_CLUB::from_ITMO_student_employee").
	Row().
	DataT("keyboard.is_graduate_employee", "JOIN_CLUB::from_ITMO_graduate_employee").
	Row().
	DataT("keyboard.is_not_itmo", "JOIN_CLUB::not_from_ITMO").
	Build()

var InlineKbd_RelevancePhoneNumber = DefaultInline().
	DataT("keyboard.actual_number", "RELEVANC_PHONE::yes").
	Row().
	DataT("keyboard.not_actual_number", "RELEVANC_PHONE::not").
	Build()

var InlineKbd_Appointment = DefaultInline().
	DataT("keyboard.is_student", "APPOINTMENT::from_ITMO").
	Row().
	DataT("keyboard.is_not_itmo", "APPOINTMENT::not_from_ITMO").
	Build()

var Keyboard_RequestContact = Default().
	RequestContactT("keyboard.send_my_number").
	TextT("keyboard.not_my_number").
	Row().
	TextT("keyboard.to_main_menu").
	Build()
