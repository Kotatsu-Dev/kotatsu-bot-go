// ------------------------------------
// RR IT 2024
//
// ------------------------------------

//
// ----------------------------------------------------------------------------------
//
// 								Обработчики сообщений боту
//
// ----------------------------------------------------------------------------------
//

package main

import (
	"rr/kotatsutgbot/cb"

	//Сторонние библиотеки
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	//Системные пакеты
	"context"
)

//
// Главные процессы
//

func BotHandler_Default(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update == nil {
		return
	}

	if update.Message != nil {
		cb.MAIN.Execute(ctx, b, update)
	} else if update.CallbackQuery != nil {
		cb.CALLBACK_MAIN.Execute(ctx, b, update)
	}
}

//
//	Команды
//

// Главное меню
func BotHandler_Command_Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	cb.START.Execute(ctx, b, update)
}

func BotHandler_Command_Start_Link(ctx context.Context, b *bot.Bot, update *models.Update) {
	cb.START_LINK.Execute(ctx, b, update)
}
