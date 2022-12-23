package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/redstar01/sixtysec/internal/entity"
)

// quizRepo -.
type quizRepo struct {
	db *sqlx.DB
}

// New -.
func New(db *sqlx.DB) *quizRepo {
	return &quizRepo{db: db}
}

func (q quizRepo) GetRandomPuzzle(ctx context.Context) (*entity.Quiz, error) {
	var quiz entity.Quiz
	err := q.db.GetContext(ctx, &quiz, `SELECT id, created_at, updated_at, deleted_at, question  FROM quizzes ORDER BY RANDOM() LIMIT 1;`)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetRandomPuzzle: %w", err)
	}
	if err != nil {
		return nil, err
	}

	var answers []entity.Answer
	err = q.db.SelectContext(ctx, &answers, `SELECT id, created_at, updated_at, deleted_at, text, is_correct FROM answers WHERE quiz_id = ?`, quiz.ID)
	if err != nil {
		return nil, err
	}

	quiz.Answers = answers

	return &quiz, nil
}

func (q quizRepo) NewGame(ctx context.Context) (*entity.Game, error) {
	n := time.Now()
	game := entity.Game{
		TotalQuestions:       0,
		FailedQuestions:      0,
		NotAnsweredQuestions: 0,
		StartTime:            &n,
		Times:                entity.Times{CreatedAt: &n},
	}
	r, err := q.db.ExecContext(ctx,
		`INSERT INTO games
				(created_at, total_questions, failed_questions, not_answered_questions, start_time)
				VALUES(?, ?, ?, ?, ?)`, game.CreatedAt, game.TotalQuestions, game.FailedQuestions, game.NotAnsweredQuestions, game.StartTime)
	if err != nil {
		return nil, err
	}

	lastId, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	game.ID = lastId

	return &game, nil
}

func (q quizRepo) GetGameByID(ctx context.Context, gameID int) (*entity.Game, error) {
	var game entity.Game
	err := q.db.GetContext(ctx, &game, `SELECT * FROM games g WHERE g.id = ? AND g.deleted_on IS NULL`, gameID)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (q quizRepo) FinishGameByID(ctx context.Context, gameID int) (bool, error) {
	_, err := q.db.ExecContext(ctx, `UPDATE games g SET g.end_time = ? WHERE g.id = ?`, time.Now(), gameID)
	if err != nil {
		return false, err
	}

	return true, nil
}
