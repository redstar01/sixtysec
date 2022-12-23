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
	}
)

// NewRouter - creates telegram command to handler router
func NewRouter(cfg *config.Config, b *telebot.Bot, ug usecase.Game) {
	r := &router{
		cfg:          cfg,
		b:            b,
		ug:           ug,
		gameStoppers: make(map[int64]context.CancelFunc),
	}
	r.b.Handle("/help", r.help)
	r.b.Handle("/newGame", r.startGame)
	r.b.Handle("/stopGame", r.stopGame)

	r.b.Handle(telebot.OnPollAnswer, r.answerHandle)
}
