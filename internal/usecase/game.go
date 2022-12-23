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
	cfg  *config.Config
	repo GameRepo

	qLock       sync.RWMutex
	cancelState map[Player]chan struct{}
}

// New -.
func New(cfg *config.Config, q GameRepo) Game {
	return &gameUseCase{
		cfg:         cfg,
		repo:        q,
		cancelState: map[Player]chan struct{}{},
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
	}()

	return quizzes, nil
}

func (q *gameUseCase) Answer(p Player) {
	q.qLock.RLock()
	if _, ok := q.cancelState[p]; ok {
		q.cancelState[p] <- struct{}{}
	}
	q.qLock.RUnlock()
}
