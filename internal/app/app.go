// Package app configures and runs application.
package app

import (
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/telebot.v3"

	"github.com/redstar01/sixtysec/config"
	"github.com/redstar01/sixtysec/internal/controller/telegram"
	"github.com/redstar01/sixtysec/internal/usecase"
	"github.com/redstar01/sixtysec/internal/usecase/repo"
	"github.com/redstar01/sixtysec/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.LogLevel)

	db, err := sqlx.Connect("sqlite3", "_data/main.db")
	if err != nil {
		l.Fatal(err)
	}

	defer func() { _ = db.Close() }()

	pref := telebot.Settings{
		Token:  cfg.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		l.Fatal(err)
	}

	quizRepo := repo.New(db)

	ucq := usecase.New(cfg, quizRepo)

	telegram.NewRouter(cfg, b, ucq)

	b.Start()
}
