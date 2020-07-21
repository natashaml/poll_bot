package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserFlowCaseSensitive(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	reply, err := generateReplyFor(s, newSubscribeCallback(t, "123"))
	require.NoError(t, err)
	require.Equal(t, reply, "Добрый день, Vasya. Добро пожаловать")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы гражданин республики Беларусь?")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "да"))
	require.NoError(t, err)
	require.Equal(t, reply, "Ващ возраст?")
}

func TestUserFlow(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	reply, err := generateReplyFor(s, newSubscribeCallback(t, "123"))
	require.NoError(t, err)
	require.Equal(t, reply, "Добрый день, Vasya. Добро пожаловать")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы гражданин республики Беларусь?")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Нет"))
	require.NoError(t, err)
	require.Equal(t, reply, "Надо ответить да!")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Да"))
	require.NoError(t, err)
	require.Equal(t, reply, "Ващ возраст?")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "16"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вам должно быть 18 или больше")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "39"))
	require.NoError(t, err)
	require.Equal(t, reply, "Укажите вас регион?")

	reply, err = generateReplyFor(s, newSeenCallback(t, "123"))
	require.NoError(t, err)
	require.Equal(t, reply, "")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Берлин"))
	require.NoError(t, err)
	require.Equal(t, reply, "Какой ваш кандидат?")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Лукашенко"))
	require.NoError(t, err)
	require.Equal(t, reply, "Спасибо за голосование!")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы уже проголосовали за Лукашенко")

	reply, err = generateReplyFor(s, newTextCallback(t, "123", "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы уже проголосовали за Лукашенко")

	reply, err = generateReplyFor(s, newSubscribeCallback(t, "123"))
	require.NoError(t, err)
	require.Equal(t, reply, "")

	reply, err = generateReplyFor(s, newSeenCallback(t, "123"))
	require.NoError(t, err)
	require.Equal(t, reply, "")
}

func newSubscribeCallback(t *testing.T, id string) *ViberCallback {
	c := &ViberCallback{
		Event: "subscribed",
		User: User{
			Id:   id,
			Name: "Vasya",
		},
	}

	b, err := json.Marshal(c)
	require.NoError(t, err)

	ret, err := parseCallback(b)
	require.NoError(t, err)

	return ret
}

func newTextCallback(t *testing.T, id string, text string) *ViberCallback {
	c := &ViberCallback{
		Event: "message",
		User: User{
			Id:   id,
			Name: "Vasya",
		},
		Message: Message{
			Text: text,
		},
	}
	b, err := json.Marshal(c)
	require.NoError(t, err)

	ret, err := parseCallback(b)
	require.NoError(t, err)

	return ret
}

func newSeenCallback(t *testing.T, id string) *ViberCallback {
	c := &ViberCallback{
		Event: "seen",
		User: User{
			Id: id,
		},
	}
	b, err := json.Marshal(c)
	require.NoError(t, err)

	ret, err := parseCallback(b)
	require.NoError(t, err)

	return ret
}
