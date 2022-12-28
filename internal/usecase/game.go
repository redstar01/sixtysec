package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/redstar01/sixtysec/config"
	"github.com/redstar01/sixtysec/internal/entity"
)

type Player int64

// gameUseCase -.
type gameUseCase struct {
	cfg *config.Config

	repo GameRepo

	pLock sync.RWMutex
	pr    ProgressRepo

	qLock       sync.RWMutex
	cancelState map[Player]chan struct{}
}

// New -.
func New(cfg *config.Config, q GameRepo, pr ProgressRepo) Game {
	return &gameUseCase{
		cfg:         cfg,
		repo:        q,
		cancelState: map[Player]chan struct{}{},
		pr:          pr,
	}
}

func (q *gameUseCase) GameStart(ctx context.Context, p Player, questionCount int) (chan *entity.Quiz, error) {
	quizzes := make(chan *entity.Quiz, questionCount)

	go func() {
		for i := 0; i < questionCount; i++ {
			quiz, err := q.repo.GetRandomPuzzle(context.Background())
			if err != nil {
				const timeToSleep = 10

				time.Sleep(time.Second * timeToSleep)

				continue
			}

			select {
			case <-ctx.Done():
				close(quizzes)

				return
			default:
			}

			select {
			case <-ctx.Done():
				close(quizzes)

				return
			case quizzes <- quiz:
				q.qLock.Lock()
				if _, ok := q.cancelState[p]; !ok {
					q.cancelState[p] = make(chan struct{})
				}
				q.qLock.Unlock()

				select {
				case <-ctx.Done():
					close(quizzes)

					return
				case <-time.After(time.Second * time.Duration(q.cfg.GameSpeed)):
					continue
				case <-q.cancelState[p]:
					continue
				}
			}
		}
		close(quizzes)

		return
	}()

	return quizzes, nil
}

func (q *gameUseCase) AnswerCheck(quizID, answerID int64, p Player) error {
	quiz, err := q.repo.GetPuzzleByID(context.Background(), quizID)
	if err != nil {
		return err
	}

	isCorrect := false

	for _, answer := range quiz.Answers {
		if answer.IsCorrect && answer.ID == answerID {
			isCorrect = true
		}
	}

	if err := q.progressUpdate(p, isCorrect); err != nil {
		return err
	}

	q.pushNextQuestion(p)

	return nil
}

func (q *gameUseCase) GetAndFlushProgress(p Player) entity.GameProgress {
	q.pLock.Lock()
	defer q.pLock.Unlock()

	gp, err := q.pr.Get(p)
	if err != nil {
		return entity.GameProgress{}
	}

	if err = q.pr.Delete(p); err != nil {
		return gp
	}

	return gp
}

func (q *gameUseCase) progressUpdate(p Player, isCorrect bool) error {
	q.pLock.Lock()
	defer q.pLock.Unlock()

	gp, err := q.pr.Get(p)
	if err != nil {
	}

	if isCorrect {
		gp.SuccessQuestions++
	} else {
		gp.FailedQuestions++
	}

	if err := q.pr.Set(p, gp); err != nil {
		return err
	}

	return nil
}

func (q *gameUseCase) pushNextQuestion(p Player) {
	q.qLock.RLock()
	if _, ok := q.cancelState[p]; ok {
		q.cancelState[p] <- struct{}{}
	}
	q.qLock.RUnlock()
}
