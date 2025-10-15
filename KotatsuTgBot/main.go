package main

import (
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/routes"
	"rr/kotatsutgbot/rr_debug"

	//Сторонние библиотеки
	"github.com/go-telegram/bot"

	//Системные пакеты
	"context"
	"os"
	"os/signal"
)

func main() {

	db.DB_Init()

	// Запускаем роутер
	go routes.RunServer()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(BotHandler_Default),
	}

	b, err := bot.New(config.GetConfig().CONFIG_BOT_TOKEN, opts...)
	if err != nil {
		rr_debug.PrintLOG("main.go", "main()", "bot.New", "Ошибка инициализации бота", err.Error())
		return
	}

	StartCron(b)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, BotHandler_Command_Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, BotHandler_Command_Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/login", bot.MatchTypeExact, BotHandler_Command_Login)
	b.Start(ctx)
}
