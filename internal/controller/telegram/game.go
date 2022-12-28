package telegram

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/redstar01/sixtysec/internal/usecase"
)

func (r *router) startGame(c telebot.Context) error {
	r.gsLock.Lock()
	if _, ok := r.gameStoppers[c.Chat().ID]; ok {
		_ = c.Send("Игра уже начаталась, остановите ее /stopgame если хочешь новую")
		r.gsLock.Unlock()

		return nil
	}
	r.gsLock.Unlock()

	ctx, cancel := context.WithCancel(context.Background())

	r.gsLock.Lock()
	r.gameStoppers[c.Chat().ID] = cancel
	r.gsLock.Unlock()

	_ = c.Send("Начинаем игру, тебя ждут 10 вопросов. Чтобы остановить - /stopgame")

	time.Sleep(time.Second * 1)

	quizzes, err := r.ug.GameStart(ctx, usecase.Player(c.Chat().ID), r.cfg.GameQuestionsCount)
	if err != nil {
		return err
	}

	for quiz := range quizzes {
		var options []telebot.PollOption

		var (
			correctOptionID    int64
			correctOptionIndex int
		)

		for i, answer := range quiz.Answers {
			if answer.IsCorrect {
				correctOptionIndex = i
				correctOptionID = answer.ID
			}

			options = append(options, telebot.PollOption{Text: answer.Text})
		}

		poll := &telebot.Poll{
			Type:          telebot.PollQuiz,
			Question:      quiz.Question,
			Options:       options,
			CorrectOption: correctOptionIndex,
			Anonymous:     false,
			OpenPeriod:    r.cfg.GameSpeed,
		}

		msg, err := r.b.Send(c.Sender(), poll)
		if err != nil {
			continue
		}

		pci := pollCacheItem{
			quizID:             quiz.ID,
			correctAnswerID:    correctOptionID,
			correctOptionIndex: correctOptionIndex,
		}

		r.c.Set(msg.Poll.ID, pci, time.Duration(r.cfg.CacheDefaultExpiration)*time.Minute)
	}

	gameProgress := r.ug.GetAndFlushProgress(usecase.Player(c.Chat().ID))
	_ = c.Send(fmt.Sprintf("Успешных ответов %d, неправильных ответов %d, еще разок? /newgame", gameProgress.SuccessQuestions, gameProgress.FailedQuestions))

	return nil
}

var (
	errCacheNotFound       = errors.New("answerHandle - poll is not found in the cache")
	errCacheTypeUnexpected = errors.New("answerHandle - poll has not appropriate type")
)

func (r *router) answerHandle(c telebot.Context) error {
	cache, ok := r.c.Get(c.PollAnswer().PollID)
	if !ok {
		return errCacheNotFound
	}

	pollCache, ok := cache.(pollCacheItem)
	if !ok {
		return errCacheTypeUnexpected
	}

	var correctAnswerID int64 = -1

	for _, chosenOption := range c.PollAnswer().Options {
		if pollCache.correctOptionIndex == chosenOption {
			correctAnswerID = pollCache.correctAnswerID

			break
		}
	}

	err := r.ug.AnswerCheck(pollCache.quizID, correctAnswerID, usecase.Player(c.Sender().ID))
	if err != nil {
		return err
	}

	return nil
}

func (r *router) stopGame(c telebot.Context) error {
	r.gsLock.Lock()
	if stopper, ok := r.gameStoppers[c.Sender().ID]; ok {
		stopper()
		delete(r.gameStoppers, c.Sender().ID)
	} else {
		_ = c.Send(fmt.Sprintf("Текущая игра уже была закончена, поиграть жми /newgame"))
	}
	r.gsLock.Unlock()

	return nil
}
