package telegram

import "time"

type pollCacheItem struct {
	quizID             int64
	correctAnswerID    int64
	correctOptionIndex int
}

type PollCache interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string) (interface{}, bool)
}
