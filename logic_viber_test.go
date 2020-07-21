package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserFlowCaseSensitive(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userId := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, reply, "Добрый день, Vasya. Добро пожаловать")

	text := newTextCallback(t, userId, "Привет")
	require.Equal(t, text.User.Id, userId)
	reply, err = generateReplyFor(p, s, text)
	require.NoError(t, err)
	require.Equal(t, reply, "Вы гражданин республики Беларусь?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "да"))
	require.NoError(t, err)
	require.Equal(t, reply, "Ваш возраст?")
}

func TestUserFlow(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userId := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, reply, "Добрый день, Vasya. Добро пожаловать")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы гражданин республики Беларусь?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Нет"))
	require.NoError(t, err)
	require.Equal(t, reply, "Надо ответить да!")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Да"))
	require.NoError(t, err)
	require.Equal(t, reply, "Ваш возраст?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "16"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вам должно быть 18 или больше")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "39"))
	require.NoError(t, err)
	require.Equal(t, reply, "Укажите вас регион?")

	user, err := s.fromPersisted(userId)
	require.NoError(t, err)

	require.Equal(t, user.Id, userId)
	require.Equal(t, user.Age, 39)
	require.Equal(t, user.Level, 3)

	reply, err = generateReplyFor(p, s, newSeenCallback(t, userId))
	require.NoError(t, err)
	require.Equal(t, reply, "")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Берлин"))
	require.NoError(t, err)
	require.Equal(t, reply, "Какой ваш кандидат?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Лукашенко"))
	require.NoError(t, err)
	require.Equal(t, reply, "Спасибо за голосование!")

	user, err = s.fromPersisted(userId)
	require.NoError(t, err)

	require.Equal(t, user.Id, userId)
	require.Equal(t, user.Age, 39)
	require.Equal(t, user.Level, 4)
	require.Equal(t, user.Candidate, "лукашенко")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы уже проголосовали за Лукашенко")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userId, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply, "Вы уже проголосовали за Лукашенко")

	subscribe := newSubscribeCallback(t, userId)
	user, err = s.Obtain(userId)
	require.NoError(t, err)
	reply, err = generateReplyFor(p, s, subscribe)
	require.NoError(t, err)
	require.Equal(t, reply, "")

	reply, err = generateReplyFor(p, s, newSeenCallback(t, userId))
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
	json := `{"event":"message","sender":{"id":"%s","Name":"Vasya"},"message":{"type":"text","text":"%s"}}`

	validJson := fmt.Sprintf(json, id, text)

	ret, err := parseCallback([]byte(validJson))
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
