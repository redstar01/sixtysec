package entity

import "time"

type Game struct {
	ID                   int64
	Player               int64
	TotalQuestions       int        `db:"total_questions"`
	FailedQuestions      int        `db:"failed_questions"`
	NotAnsweredQuestions int        `db:"not_answered_questions"`
	StartTime            *time.Time `db:"start_time"`
	EndTime              *time.Time `db:"end_time"`
	Times
}

type Quiz struct {
	ID       int64
	Question string
	Answers  []Answer
	Times
}

type Answer struct {
	ID        int64
	Text      string
	IsCorrect bool `db:"is_correct"`
	Times
}

type Times struct {
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type GameProgress struct {
	SuccessQuestions int
	FailedQuestions  int
}
