package telegram

import (
	"gopkg.in/telebot.v3"
)

func (r *router) help(c telebot.Context) error {
	return c.Send(`Пccc, в 60 секунд сыграть хочешь?! 😏 Кликай /newGame и погнали🔥`)
}
