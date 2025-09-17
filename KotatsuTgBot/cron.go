package main

import (
	"context"
	"fmt"
	"math/rand"
	"rr/kotatsutgbot/config"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"
	"time"

	"github.com/go-telegram/bot"
)

func StartCron(b *bot.Bot) {
	minute_ticker := time.NewTicker(1 * time.Minute)
	halfhour_ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for {
			check_roulette(b)
			<-minute_ticker.C
		}
	}()

	go func() {
		for {
			check_step(b)
			<-halfhour_ticker.C
		}
	}()
}

func check_roulette(b *bot.Bot) {
	rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Начинаем проверку рулеток", "")

	db_answer_code, roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	now := time.Now()
	a_hour_ago := now.Add(-1 * time.Minute)

	if db_answer_code != db.DB_ANSWER_SUCCESS {
		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Нет рулетки", "")
		db_answer_code, roulette = db.DB_GET_AnimeRoulette_BY_Status(false)
		if db_answer_code != db.DB_ANSWER_SUCCESS {
			rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Точно нет", "")
			return
		}
		if roulette.EndDate.After(a_hour_ago) && roulette.EndDate.Before(now) {
			rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Рулетка закончилась", "")
			for _, member := range roulette.Participants {
				params := &bot.SendMessageParams{
					ChatID: member.UserTgID,
					Text: "[РУЛЕТКА]" + "\n" +
						"Рулетка завершилась!",
				}

				b.SendMessage(context.TODO(), params)
			}
		}
		return
	}

	if roulette.AnnounceDate.After(a_hour_ago) && roulette.AnnounceDate.Before(now) {
		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Регистрация закончилась", "")
		for _, member := range roulette.Participants {
			params := &bot.SendMessageParams{
				ChatID: member.UserTgID,
				Text: "[РУЛЕТКА]" + "\n" +
					"Завершилась регистрация на рулетку. Теперь необходимо загадать тайтл на тему:" + "\n" +
					roulette.Theme,
			}

			b.SendMessage(context.TODO(), params)
		}
	} else if roulette.DistributionDate.After(a_hour_ago) && roulette.DistributionDate.Before(now) {
		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Сбор названий закончился", "")
		distr := rand.Perm(len(roulette.Participants))
		distr32 := make([]int32, len(distr))
		for i := range distr {
			distr32[i] = int32(distr[i])
		}

		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Перемудрили участников", "")
		res := db.DB_UPDATE_AnimeRoulette_SET_Distribution(distr32)

		if res != db.DB_ANSWER_SUCCESS {
			rr_debug.PrintLOG("cron.go", "check_roulette", "ERROR", "Ошибка сохранения рулетки", fmt.Sprint(res))
			return
		}

		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Рассылаем приглашения", "")
		for _, j := range distr {
			member := roulette.Participants[j]
			next := roulette.Participants[distr[(j+1)%len(distr)]]
			params := &bot.SendMessageParams{
				ChatID: member.UserTgID,
				Text: "[РУЛЕТКА]" + "\n" +
					"Завершился сбор названий. Теперь вам необходимо посмотреть тайтл:" + "\n" +
					next.EnigmaticTitle,
			}

			b.SendMessage(context.TODO(), params)
		}

	}
}

func check_step(b *bot.Bot) {
	rr_debug.PrintLOG("cron.go", "check_step", "INFO", "Начинаем проверку шагов", "")

	users_outdated := db.DB_GET_Users_BY_Step(config.STEP_ACTIVITY_OUTDATED)

	for _, user := range users_outdated {
		text := "Ты прочитал(а) описание мероприятия, но не нажал(а) «Запиши меня».\n" +
			"Если хочешь записаться на мероприятие, выбери его ещё раз и не забудь нажать кнопку под описанием."

		if !user.IsITMO {
			text += " Без записи я не смогу попросить для тебя пропуск в Университет."
		}

		params := &bot.SendMessageParams{
			ChatID: user.UserTgID,
			Text:   text,
		}

		b.SendMessage(context.TODO(), params)
		db.DB_UPDATE_User(map[string]interface{}{
			"user_tg_id": user.UserTgID,
			"step":       config.STEP_DEFAULT,
		})
	}

	users_active := db.DB_GET_Users_BY_Step(config.STEP_ACTIVITY)
	for _, user := range users_active {
		db.DB_UPDATE_User(map[string]interface{}{
			"user_tg_id": user.UserTgID,
			"step":       config.STEP_ACTIVITY_OUTDATED,
		})
	}
}
