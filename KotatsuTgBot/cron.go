package main

import (
	"context"
	"fmt"
	"math/rand"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"
	"time"

	"github.com/go-telegram/bot"
)

func StartCron(b *bot.Bot) {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			check_roulette(b)
			<-ticker.C
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
