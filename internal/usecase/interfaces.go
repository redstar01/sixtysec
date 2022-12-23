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
		Answer(p Player)
	}

	// GameRepo - game repository.
	GameRepo interface {
		GetRandomPuzzle(context.Context) (*entity.Quiz, error)
		NewGame(context.Context) (*entity.Game, error)
		GetGameByID(ctx context.Context, gameID int) (*entity.Game, error)
		FinishGameByID(ctx context.Context, gameID int) (bool, error)
	}
)