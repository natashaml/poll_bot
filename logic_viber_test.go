package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserFlow(t *testing.T) {
	s := newStorage()
	reply, err := generateReplyFor(s, newSubscribeCallback("123"))
	require.NoError(t, err)
	require.Equal(t, reply, "Добрый день, Vasya. Добро пожаловать")

	reply, err = generateReplyFor(s, newTextCallback("123", "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы гражданин республики Беларусь?")

	reply, err = generateReplyFor(s, newTextCallback("123", "Нет"))
	require.NoError(t, err)
	require.Equal(t, reply, "Надо ответить да!")

	reply, err = generateReplyFor(s, newTextCallback("123", "Да"))
	require.NoError(t, err)
	require.Equal(t, reply, "Ващ возраст?")

	reply, err = generateReplyFor(s, newTextCallback("123", "16"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вам должно быть 18 или больше")

	reply, err = generateReplyFor(s, newTextCallback("123", "39"))
	require.NoError(t, err)
	require.Equal(t, reply, "Укажите вас регион?")

	reply, err = generateReplyFor(s, newSeenCallback("123"))
	require.NoError(t, err)
	require.Equal(t, reply, "")

	reply, err = generateReplyFor(s, newTextCallback("123", "Берлин"))
	require.NoError(t, err)
	require.Equal(t, reply, "Какой ваш кандидат?")

	reply, err = generateReplyFor(s, newTextCallback("123", "Лукашенко"))
	require.NoError(t, err)
	require.Equal(t, reply, "Спасибо за голосование!")

	reply, err = generateReplyFor(s, newTextCallback("123", "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы уже проголосовали за Лукашенко")

	reply, err = generateReplyFor(s, newSubscribeCallback("123"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы уже проголосовали за Лукашенко")

	reply, err = generateReplyFor(s, newSeenCallback("123"))
	require.NoError(t, err)
	require.Equal(t, reply, "")
}

func newSubscribeCallback(id string) *ViberCallback {
	return &ViberCallback{
		Event: "subscribed",
		User: User{
			Id:   id,
			Name: "Vasya",
		},
	}
}

func newTextCallback(id string, text string) *ViberCallback {
	return &ViberCallback{
		Event: "message",
		User: User{
			Id:   id,
			Name: "Vasya",
		},
		Message: Message{
			Text: text,
		},
	}
}

func newSeenCallback(id string) *ViberCallback {
	return &ViberCallback{
		Event: "seen",
		User: User{
			Id: id,
		},
	}
}
