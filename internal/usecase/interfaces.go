// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/redstar01/sixtysec/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Game - main game logic.
	Game interface {
		GameStart(ctx context.Context, p Player, questionCount int) (chan *entity.Quiz, error)
		AnswerCheck(quizID, answerID int64, p Player) error
		GetAndFlushProgress(p Player) GameState
	}

	// GameRepo - game repository.
	GameRepo interface {
		GetRandomPuzzle(context.Context) (*entity.Quiz, error)
		GetPuzzleByID(ctx context.Context, quizID int64) (*entity.Quiz, error)
	}
)
