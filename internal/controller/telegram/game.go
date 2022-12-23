package telegram

import (
	"context"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/redstar01/sixtysec/internal/usecase"
)

func (r *router) startGame(c telebot.Context) error {
	ctx, cancel := context.WithCancel(context.Background())

	r.gsLock.Lock()
	r.gameStoppers[c.Chat().ID] = cancel
	r.gsLock.Unlock()

	_ = c.Send("Начинаем игру, тебя ждут 10 вопросов. Чтобы остановить - /stopGame")
	time.Sleep(time.Second * 3)

	quizzes, err := r.ug.GameStart(ctx, usecase.Player(c.Chat().ID), 10)
	if err != nil {
		return err
	}

	for quiz := range quizzes {
		var options []telebot.PollOption

		var correctOption int
		for i, answer := range quiz.Answers {
			if answer.IsCorrect {
				correctOption = i
			}
			options = append(options, telebot.PollOption{Text: answer.Text})
		}

		poll := &telebot.Poll{
			Type:          telebot.PollQuiz,
			Question:      quiz.Question,
			Options:       options,
			CorrectOption: correctOption,
			Anonymous:     false,
			OpenPeriod:    r.cfg.GameSpeed,
		}

		_ = c.Send(poll)
	}

	_ = c.Send("Спасибо за игру, еще разок? /newGame")

	return nil
}

func (r *router) answerHandle(c telebot.Context) error {
	r.ug.Answer(usecase.Player(c.Sender().ID))
	return nil
}

func (r *router) stopGame(c telebot.Context) error {
	r.gsLock.RLock()
	if stopper, ok := r.gameStoppers[c.Sender().ID]; ok {
		stopper()
	}
	r.gsLock.RUnlock()

	return nil
}
