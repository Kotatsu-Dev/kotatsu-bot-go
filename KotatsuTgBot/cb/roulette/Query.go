package roulette

import (
	"context"
	. "rr/kotatsutgbot/cb"
	"rr/kotatsutgbot/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func MainQuery(ctx context.Context, b *bot.Bot, update *models.Update, current_user *db.User_ReadJSON) {
	AnswerQuery(ctx, b, update)
	Main(ctx, b, update, current_user)
}
