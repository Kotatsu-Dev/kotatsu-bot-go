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

	if config.GetConfig().ROULETTES {
		go func() {
			for {
				check_roulette(b)
				<-minute_ticker.C
			}
		}()
	}

	go func() {
		for {
			check_step(b)
			<-halfhour_ticker.C
		}
	}()
}

func check_roulette(b *bot.Bot) {
	now := time.Now()
	a_hour_ago := now.Add(-1 * time.Minute)

	if db_answer_code, ended := db.DB_GET_AnimeRoulette_BY_Status(false); db_answer_code == db.DB_ANSWER_SUCCESS {
		if ended.EndDate.After(a_hour_ago) && ended.EndDate.Before(now) {
			rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Рулетка закончилась", "")
			for _, member := range ended.Participants {
				params := &bot.SendMessageParams{
					ChatID: member.UserTgID,
					Text:   config.T("roulette.messages.ended"),
				}

				b.SendMessage(context.TODO(), params)
			}
		}
	}

	db_answer_code, roulette := db.DB_GET_AnimeRoulette_BY_Status(true)
	if db_answer_code != db.DB_ANSWER_SUCCESS {
		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Нет рулетки", "")
		return
	}

	if roulette.AnnounceDate.After(a_hour_ago) && roulette.AnnounceDate.Before(now) {
		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Регистрация закончилась", "")
		for _, member := range roulette.Participants {
			params := &bot.SendMessageParams{
				ChatID: member.UserTgID,
				Text:   config.TT("roulette.messages.registration_ended", roulette),
			}

			b.SendMessage(context.TODO(), params)
		}
	} else if roulette.DistributionDate.After(a_hour_ago) && roulette.DistributionDate.Before(now) {
		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Сбор названий закончился", "")
		shuffled := make([]db.User_ReadJSON, len(roulette.Participants))
		copy(shuffled, roulette.Participants)
		rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		distr32 := make([]int32, len(shuffled))
		for i, member := range shuffled {
			distr32[i] = int32(member.ID)
		}

		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Перемудрили участников", "")
		res := db.DB_UPDATE_AnimeRoulette_SET_Distribution(roulette.ID, distr32)

		if res != db.DB_ANSWER_SUCCESS {
			rr_debug.PrintLOG("cron.go", "check_roulette", "ERROR", "Ошибка сохранения рулетки", fmt.Sprint(res))
			return
		}

		rr_debug.PrintLOG("cron.go", "check_roulette", "INFO", "Рассылаем приглашения", "")
		for i, member := range shuffled {
			next := shuffled[(i+1)%len(shuffled)]
			params := &bot.SendMessageParams{
				ChatID: member.UserTgID,
				Text:   config.TT("roulette.messages.selection_ended", next),
			}

			b.SendMessage(context.TODO(), params)
		}

	}
}

func check_step(b *bot.Bot) {
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
