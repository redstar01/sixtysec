package telegram

import (
	"context"
	"sync"

	"gopkg.in/telebot.v3"

	"github.com/redstar01/sixtysec/config"
	"github.com/redstar01/sixtysec/internal/usecase"
)

type (
	router struct {
		cfg *config.Config
		b   *telebot.Bot
		ug  usecase.Game

		gsLock       sync.RWMutex
		gameStoppers map[int64]context.CancelFunc

		c PollCache
	}
)

// NewRouter - creates telegram command to handler router.
func NewRouter(cfg *config.Config, b *telebot.Bot, ug usecase.Game, c PollCache) {
	r := &router{
		cfg:          cfg,
		b:            b,
		ug:           ug,
		gameStoppers: make(map[int64]context.CancelFunc),
		c:            c,
	}
	r.b.Handle("/help", r.help)
	r.b.Handle("/start", r.help)
	r.b.Handle("/newgame", r.startGame)
	r.b.Handle("/stopgame", r.stopGame)

	r.b.Handle(telebot.OnPollAnswer, r.answerHandle)
}
