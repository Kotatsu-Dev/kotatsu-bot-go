package main

import (
	"context"
	"math/rand"
	"rr/kotatsutgbot/db"
	"rr/kotatsutgbot/rr_debug"
	"time"

	"github.com/go-telegram/bot"
)

func StartCron(b *bot.Bot) {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for {
			check_roulette(b)
			<-ticker.C
		}
	}()
}

func check_roulette(b *bot.Bot) {
	db_answer_code, roulette := db.DB_GET_AnimeRoulette_BY_Status(true)

	if db_answer_code != db.DB_ANSWER_SUCCESS {
		return
	}

	now := time.Now()
	a_hour_ago := now.Add(-1 * time.Hour)

	if roulette.AnnounceDate.After(a_hour_ago) && roulette.AnnounceDate.Before(now) {
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
		par := roulette.Participants
		for i := range par {
			j := rand.Intn(i + 1)
			par[i], par[j] = par[j], par[i]
		}

		roulette.Participants = par
		res := db.DB_Database().Save(roulette)

		if res.Error != nil {
			rr_debug.PrintLOG("cron.go", "check_roulette", "ERROR", "Ошибка сохранения рулетки", res.Error.Error())
			return
		}

		for i, member := range roulette.Participants {
			next := roulette.Participants[(i+1)%len(roulette.Participants)]
			params := &bot.SendMessageParams{
				ChatID: member.UserTgID,
				Text: "[РУЛЕТКА]" + "\n" +
					"Завершился сбор названий. Теперь вам необходимо посмотреть тайтл:" + "\n" +
					next.EnigmaticTitle,
			}

			b.SendMessage(context.TODO(), params)
		}
	} else if roulette.EndDate.After(a_hour_ago) && roulette.EndDate.Before(now) {
		for _, member := range roulette.Participants {
			params := &bot.SendMessageParams{
				ChatID: member.UserTgID,
				Text: "[РУЛЕТКА]" + "\n" +
					"Рулетка завершилась!",
			}

			b.SendMessage(context.TODO(), params)
		}
	}
}
