package roulette

import (
	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/db"
)

func MainQuery(user *db.User_ReadJSON) Executor {
	return Do(
		AnswerQuery(),
		Main(user),
	)
}
