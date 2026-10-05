package roulette

import (
	"context"
	"time"

	. "rr/kotatsutgbot/cb/helpers"
	"rr/kotatsutgbot/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type RouletteStateT int

const (
	RouletteStateBeforeStart  RouletteStateT = iota // not open yet
	RouletteStateRegistration                       // joining is open
	RouletteStateWishing                            // participants submit their enigmatic title
	RouletteStateDistributed                        // titles handed out, watch period
	RouletteStateEnded                              // finished
)

func GetRouletteState(roulette *db.AnimeRoulette_ReadJSON) RouletteStateT {
	now := time.Now()
	switch {
	case now.Before(roulette.StartDate):
		return RouletteStateBeforeStart
	case now.Before(roulette.AnnounceDate):
		return RouletteStateRegistration
	case now.Before(roulette.DistributionDate):
		return RouletteStateWishing
	case now.Before(roulette.EndDate):
		return RouletteStateDistributed
	default:
		return RouletteStateEnded
	}
}

func RouletteIsActive(state RouletteStateT) bool {
	switch state {
	case RouletteStateRegistration, RouletteStateWishing, RouletteStateDistributed:
		return true
	default:
		return false
	}
}

func RouletteStateGuard(roulette *db.AnimeRoulette_ReadJSON, state RouletteStateT, executor Executor) *GuardS {
	return Guard(func(ctx context.Context, b *bot.Bot, u *models.Update) bool {
		return GetRouletteState(roulette) == state
	})(executor)
}

func RouletteActiveGuard(roulette *db.AnimeRoulette_ReadJSON, executor Executor) *GuardS {
	return Guard(func(ctx context.Context, b *bot.Bot, u *models.Update) bool {
		return RouletteIsActive(GetRouletteState(roulette))
	})(executor)
}

func RouletteInactiveGuard(roulette *db.AnimeRoulette_ReadJSON, executor Executor) *GuardS {
	return Guard(func(ctx context.Context, b *bot.Bot, u *models.Update) bool {
		return !RouletteIsActive(GetRouletteState(roulette))
	})(executor)
}
