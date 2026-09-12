package roulette

import (
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/db"
)

func MainQueryE(user *db.User_ReadJSON) Executor {
	return Seq(
		AnswerQueryE(),
		MainE(user),
	)
}
