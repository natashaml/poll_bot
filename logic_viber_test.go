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
	require.Equal(t, reply, "Надо ответить да!")
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
