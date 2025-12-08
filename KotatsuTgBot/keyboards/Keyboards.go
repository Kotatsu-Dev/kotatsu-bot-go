// ------------------------------------
// RR IT 2024
//
// ------------------------------------
// Базовый движок для Котацу бота

package keyboards

import (
	//Внутренние пакеты проекта
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"time"

	//Сторонние библиотеки
	"github.com/go-telegram/bot/models"

	//Системные пакеты
	"fmt"
)

// ----------------------------------------------
//
//	Структуры
//
// ----------------------------------------------
// Клавиатура для главного меню
var Registration = &models.ReplyKeyboardMarkup{
	Keyboard: [][]models.KeyboardButton{
		{
			{Text: "🗃 Продолжить"},
		},
	},
	ResizeKeyboard:  true,  // Опционально: уменьшить клавиатуру до размера кнопок
	OneTimeKeyboard: false, // Опционально: скрыть клавиатуру после использования
}

var Keyboard_GenderSelect = &models.ReplyKeyboardMarkup{
	Keyboard: [][]models.KeyboardButton{
		{
			{Text: "Повелитель демонов"},
			{Text: "Девочка волшебница"},
		},
	},
	ResizeKeyboard:  true,  // Опционально: уменьшить клавиатуру до размера кнопок
	OneTimeKeyboard: false, // Опционально: скрыть клавиатуру после использования
}

var Keyboard_WasAtEvents = &models.ReplyKeyboardMarkup{
	Keyboard: [][]models.KeyboardButton{
		{
			{Text: "Да, я уже мандаринка"},
			{Text: "Ещё нет :("},
		},
	},
	ResizeKeyboard:  true, // Опционально: уменьшить клавиатуру до размера кнопок
	OneTimeKeyboard: true, // Опционально: скрыть клавиатуру после использования
}

var Keyboard_WasntAtEvents = &models.ReplyKeyboardMarkup{
	Keyboard: [][]models.KeyboardButton{
		{
			{Text: "Хорошо, заполню позже"},
			{Text: "Хочу продолжить"},
		},
	},
	ResizeKeyboard:  true, // Опционально: уменьшить клавиатуру до размера кнопок
	OneTimeKeyboard: true, // Опционально: скрыть клавиатуру после использования
}

// Клавиатура для незарегистрированных пользователей
func CreateKeyboard_MainMenuButtonsDefault(news_letter bool) *models.ReplyKeyboardMarkup {
	/*var news_letter_text string
	if news_letter {
		news_letter_text = "❌ Отписаться от рассылки"
	} else {
		news_letter_text = "📰 Подписаться на рассылку"
	}*/

	var keyboard = &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "⛩ Вступить в клуб"},
				{Text: "📝 Запись на мероприятия"},
			},
			/*/*{
				{Text: news_letter_text},
			},*/
		},
		ResizeKeyboard:  true,  // Опционально: уменьшить клавиатуру до размера кнопок
		OneTimeKeyboard: false, // Опционально: скрыть клавиатуру после использования
	}
	return keyboard
}

// Клавиатура для главного меню участника клуба
func CreateKeyboard_MainMenuButtonsClubMember(news_letter bool) *models.ReplyKeyboardMarkup {
	/*var news_letter_text string
	if news_letter {
		news_letter_text = "❌ Отписаться от рассылки"
	} else {
		news_letter_text = "📰 Подписаться на рассылку"
	}*/

	var keyboard = &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				//{Text: news_letter_text},
				{Text: "📝 Запись на мероприятия"},
				// {Text: "🤝 Акции и партнёры"},
			},
			{
				//{Text: news_letter_text},
				{Text: "🚪 Покинуть клуб"},
				{Text: "📂 Мои мероприятия"},
			},
		},
		ResizeKeyboard:  true,  // Опционально: уменьшить клавиатуру до размера кнопок
		OneTimeKeyboard: false, // Опционально: скрыть клавиатуру после использования
	}
	return keyboard
}

var CommunicationManager = &models.ReplyKeyboardMarkup{
	Keyboard: [][]models.KeyboardButton{
		{
			{Text: "⬅ Вернуться в главное меню"},
		},
	},
	ResizeKeyboard:  true,  // Опционально: уменьшить клавиатуру до размера кнопок
	OneTimeKeyboard: false, // Опционально: скрыть клавиатуру после использования
}

// Клавиатура для главного меню участника клуба
var ListEvents = &models.ReplyKeyboardMarkup{
	Keyboard: [][]models.KeyboardButton{
		{
			{Text: "🟡 Аниме-рулетка"},
		},
		{
			{Text: "⬅ Вернуться в главное меню"},
		},
	},
	ResizeKeyboard:  true,  // Опционально: уменьшить клавиатуру до размера кнопок
	OneTimeKeyboard: false, // Опционально: скрыть клавиатуру после использования
}

// Клавиатура для выбранного мероприятия
var SelectedEvent = &models.ReplyKeyboardMarkup{
	Keyboard: [][]models.KeyboardButton{
		{
			{Text: "❌ Отменить запись"},
		},
		{
			{Text: "🟡 Аниме рулетка"},
		},
		{
			{Text: "⬅ Вернуться к списку мероприятий"},
		},
	},
	ResizeKeyboard:  true,  // Опционально: уменьшить клавиатуру до размера кнопок
	OneTimeKeyboard: false, // Опционально: скрыть клавиатуру после использования
}

// Клавиатура для аниме рулетки перед её запуском
func CreateKeyboard_AnimeRouletteStart(is_member bool) *models.ReplyKeyboardMarkup {

	var keyboard = &models.ReplyKeyboardMarkup{}

	if is_member {
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "🚪 Покинуть рулетку"},
				},
				{
					{Text: "📋 Правила"},
					{Text: "📚 Мой список"},
				},
				{
					{Text: "⬅ Вернуться в главное меню"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	} else {
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "✅ Участвовать в рулетке"},
				},
				{
					{Text: "📋 Правила"},
					{Text: "📚 Мой список"},
				},
				{
					{Text: "⬅ Вернуться в главное меню"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	}

	return keyboard
}

// Клавиатура для аниме рулетки когда запущена
func CreateKeyboard_AnimeRouletteMenu(is_member bool) *models.ReplyKeyboardMarkup {

	var keyboard = &models.ReplyKeyboardMarkup{}

	if is_member {
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "❔ Загадать аниме"},
				},
				{
					{Text: "📋 Правила"},
					{Text: "📔 Тема"},
					{Text: "📚 Мой список"},
				},
				{
					{Text: "⬅ Вернуться в главное меню"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	} else {
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "✅ Участвовать в рулетке"},
				},
				{
					{Text: "📋 Правила"},
					{Text: "📔 Тема"},
				},
				{
					{Text: "⬅ Вернуться в главное меню"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	}

	return keyboard
}

// Клавиатура возврата назад
func CreateKeyboard_Cancel(cancel_type string) *models.ReplyKeyboardMarkup {

	var keyboard = &models.ReplyKeyboardMarkup{}

	switch cancel_type {
	case "skip":
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "Пропустить"},
				},
				{
					{Text: "⬅ Вернуться в главное меню"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}

	case "anime_roulette":
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "⬅ Вернуться в меню рулетки"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}

	default:
		keyboard = &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "⬅ Вернуться в главное меню"},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
		}
	}

	return keyboard
}

// Inline-клавиатура - Список мероприятий
func CreateInlineKbd_ActivitiesList(activities []db.Activity_ReadJSON, user_tg_id int64, has_roulette bool) *models.InlineKeyboardMarkup {
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	var title string
	var formattedTime string

	// Определите желаемый формат дд.мм чч:мм
	format := "02.01, 15:04"

	for _, activity := range activities {

		is_participant := false

		for _, participant := range activity.Participants {
			if participant.UserTgID == user_tg_id {
				is_participant = true
				break
			}
		}

		loc, _ := time.LoadLocation("Europe/Moscow")
		formattedTime = activity.DateMeeting.In(loc).Format(format)
		if is_participant {
			title = "✅ [" + formattedTime + "] " + activity.Title
		} else {
			title = "[" + formattedTime + "] " + activity.Title
		}

		row := []models.InlineKeyboardButton{
			{
				Text:         title,
				CallbackData: fmt.Sprintf("ACTIVITIES::%d", activity.ID),
			},
		}
		inlineKeyboard = append(inlineKeyboard, row)
	}

	if has_roulette {
		inlineKeyboard = append(inlineKeyboard, []models.InlineKeyboardButton{
			{
				Text:         config.T("roulette.title"),
				CallbackData: "ROULETTES",
			},
		})
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

// Inline-клавиатура - Список акций и партнёров
func CreateInlineKbd_PartnersList() *models.InlineKeyboardMarkup {
	//Создаем экземпляр структуры
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	row_1 := []models.InlineKeyboardButton{
		{
			Text:         "☕️ Кафе «Тайяки»",
			CallbackData: fmt.Sprintf("PARTNERS::%s", "cafeTaiyaki"),
		},
	}

	row_2 := []models.InlineKeyboardButton{
		{
			Text:         "🌟 Фестиваль GemFest [11.11]",
			CallbackData: fmt.Sprintf("PARTNERS::%s", "gemfest"),
		},
	}

	row_back := []models.InlineKeyboardButton{
		{
			Text:         "◀️ Назад",
			CallbackData: fmt.Sprintf("PARTNERS::%s", "back"),
		},
	}

	inlineKeyboard = append(inlineKeyboard, row_1)
	inlineKeyboard = append(inlineKeyboard, row_2)
	inlineKeyboard = append(inlineKeyboard, row_back)

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

// Inline-клавиатура - Список моих мероприятий
func CreateInlineKbd_MyActivitiesList(my_activities []*db.Activity) *models.InlineKeyboardMarkup {
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	var title string
	var formattedTime string

	// Определите желаемый формат дд.мм чч:мм
	format := "02.01 15:04"
	loc, _ := time.LoadLocation("Europe/Moscow")

	for _, activity := range my_activities {

		formattedTime = activity.DateMeeting.In(loc).Format(format)
		title = "[" + formattedTime + "] " + activity.Title

		row := []models.InlineKeyboardButton{
			{
				Text:         title,
				CallbackData: fmt.Sprintf("MY_ACTIVITIES::%d", activity.ID),
			},
		}
		inlineKeyboard = append(inlineKeyboard, row)
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

// Inline-клавиатура - Подписаться на мероприятие
func CreateInlineKbd_SubscribeActivity(activity_id int) *models.InlineKeyboardMarkup {
	//Создаем экземпляр структуры
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	row := []models.InlineKeyboardButton{
		{
			Text:         "✅ Запиши меня",
			CallbackData: fmt.Sprintf("ACTIVITY_SUBSCRIBE::%d", activity_id),
		},
	}

	inlineKeyboard = append(inlineKeyboard, row)

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

// Inline-клавиатура - Отписаться от мероприятия
func CreateInlineKbd_UnsubscribeActivity(activity_id int) *models.InlineKeyboardMarkup {
	//Создаем экземпляр структуры
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	row := []models.InlineKeyboardButton{
		{
			Text:         "❌ Отмени запись",
			CallbackData: fmt.Sprintf("ACTIVITY_UNSUBSCRIBE::%d", activity_id),
		},
	}

	inlineKeyboard = append(inlineKeyboard, row)

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

// Inline-клавиатура - Вступить в клуб
func CreateInlineKbd_JoinClub() *models.InlineKeyboardMarkup {
	//Создаем экземпляр структуры
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	row_1 := []models.InlineKeyboardButton{
		{
			Text:         "Я студент/сотрудник/выпускник ИТМО",
			CallbackData: fmt.Sprintf("JOIN_CLUB::%s", "from_ITMO"),
		},
	}

	row_2 := []models.InlineKeyboardButton{
		{
			Text:         "Я не из ИТМО",
			CallbackData: fmt.Sprintf("JOIN_CLUB::%s", "not_from_ITMO"),
		},
	}

	inlineKeyboard = append(inlineKeyboard, row_1)
	inlineKeyboard = append(inlineKeyboard, row_2)

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

// Inline-клавиатура - Проверка актуальности номера телефона пользователя
func CreateInlineKbd_RelevancePhoneNumber() *models.InlineKeyboardMarkup {
	//Создаем экземпляр структуры
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	row_1 := []models.InlineKeyboardButton{
		{
			Text:         "Номер актуальный, паспорт возьму",
			CallbackData: fmt.Sprintf("RELEVANC_PHONE::%s", "yes"),
		},
	}

	row_2 := []models.InlineKeyboardButton{
		{
			Text:         "У меня поменялся номер телефона",
			CallbackData: fmt.Sprintf("RELEVANC_PHONE::%s", "no"),
		},
	}

	inlineKeyboard = append(inlineKeyboard, row_1)
	inlineKeyboard = append(inlineKeyboard, row_2)

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

// Inline-клавиатура - Запись на мероприятие (для не участников клуба)
func CreateInlineKbd_Appointment() *models.InlineKeyboardMarkup {
	//Создаем экземпляр структуры
	inlineKeyboard := [][]models.InlineKeyboardButton{}

	row_1 := []models.InlineKeyboardButton{
		{
			Text:         "Я студент/сотрудник/выпускник ИТМО",
			CallbackData: fmt.Sprintf("APPOINTMENT::%s", "from_ITMO"),
		},
	}

	row_2 := []models.InlineKeyboardButton{
		{
			Text:         "Я не из ИТМО",
			CallbackData: fmt.Sprintf("APPOINTMENT::%s", "not_from_ITMO"),
		},
	}

	inlineKeyboard = append(inlineKeyboard, row_1)
	inlineKeyboard = append(inlineKeyboard, row_2)

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}
}

func CreateKeyboard_RequestContact() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{
					Text:           "Отправить свой номер",
					RequestContact: true,
				},
				{Text: "Я не пользуюсь номером, к которому привязан Telegram"},
			}, {
				{Text: "⬅ Вернуться в главное меню"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}
