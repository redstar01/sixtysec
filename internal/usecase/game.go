package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/redstar01/sixtysec/config"
	"github.com/redstar01/sixtysec/internal/entity"
)

type GameState struct {
	SuccessQuestions int
	FailedQuestions  int
}

type Player int64

// gameUseCase -.
type gameUseCase struct {
	cfg  *config.Config
	repo GameRepo

	qLock       sync.RWMutex
	cancelState map[Player]chan struct{}

	gpLock       sync.RWMutex
	gameProgress map[Player]GameState
}

// New -.
func New(cfg *config.Config, q GameRepo) Game {
	return &gameUseCase{
		cfg:          cfg,
		repo:         q,
		cancelState:  map[Player]chan struct{}{},
		gameProgress: make(map[Player]GameState),
	}
}

func (q *gameUseCase) GameStart(ctx context.Context, p Player, questionCount int) (chan *entity.Quiz, error) {
	quizzes := make(chan *entity.Quiz, questionCount)

	go func() {
		for i := 0; i < questionCount; i++ {
			quiz, err := q.repo.GetRandomPuzzle(context.Background())
			if err != nil {
				time.Sleep(time.Second * 10)

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

	correctAnswered := false

	for _, answer := range quiz.Answers {
		if answer.IsCorrect && answer.ID == answerID {
			correctAnswered = true
		}
	}

	q.gpLock.Lock()
	if gp, ok := q.gameProgress[p]; ok {
		if correctAnswered {
			gp.SuccessQuestions++
			q.gameProgress[p] = gp
		} else {
			gp.FailedQuestions++
			q.gameProgress[p] = gp
		}
	} else {
		if correctAnswered {
			q.gameProgress[p] = GameState{SuccessQuestions: 1}
		} else {
			q.gameProgress[p] = GameState{FailedQuestions: 1}
		}
	}
	q.gpLock.Unlock()

	q.qLock.RLock()
	if _, ok := q.cancelState[p]; ok {
		q.cancelState[p] <- struct{}{}
	}
	q.qLock.RUnlock()

	return nil
}

func (q *gameUseCase) GetAndFlushProgress(p Player) GameState {
	q.gpLock.Lock()
	defer q.gpLock.Unlock()

	if gs, ok := q.gameProgress[p]; ok {
		delete(q.gameProgress, p)

		return gs
	}

	return GameState{}
}
