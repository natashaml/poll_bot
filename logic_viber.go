package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func knownNotMessageEvent(c *ViberCallback) bool {
	return c.Event == "delivered" || c.Event == "seen" || c.Event == "subscribed" || c.Event == "conversation_started"
}

func generateReplyFor(s *Storage, c *ViberCallback) (string, error) {
	if c.Event == "delivered" || c.Event == "seen" {
		return "", nil
	}

	if c.Message.Text == "clear" {
		err := s.Clear(c.User.Id)
		return fmt.Sprintf("Your storage cleared with %v", err), nil
	}

	storageUser, err := s.Obtain(c.User.Id)
	if err != nil {
		return "", err
	}

	if c.Event == "message" {
		if storageUser.Candidate != "" {
			return "Вы уже проголосовали за " + strings.Title(storageUser.Candidate), nil
		}

		message := analyseAnswer(storageUser, c)
		if message != "" {
			return message, nil
		}
		if storageUser.Level == 4 {
			err = s.Persist(storageUser.Id)
			if err != nil {
				return "", err
			}

			return "Спасибо за голосование!", nil
		}

		storageUser.Level++
		err = s.Persist(storageUser.Id)
		if err != nil {
			return "", err
		}
		message = getMessageForLevel(storageUser.Level, c)
		return message, nil
	}

	if !storageUser.ConversationStarted {
		message := getMessageForLevel(storageUser.Level, c)
		storageUser.ConversationStarted = true
		err = s.Persist(storageUser.Id)
		if err != nil {
			return "", err
		}
		return message, nil
	}

	if knownNotMessageEvent(c) {
		return "", nil
	}

	return "", fmt.Errorf("Unknown message %v", c.Event)
}

func getMessageForLevel(level int, c *ViberCallback) string {
	if level == 0 {
		return "Добрый день, " + c.User.Name + ". Добро пожаловать"
	} else if level == 1 {
		return "Вы гражданин республики Беларусь?"
	} else if level == 2 {
		return "Ващ возраст?"
	} else if level == 3 {
		return "Укажите вас регион?"
	} else if level == 4 {
		return "Какой ваш кандидат?"
	} else {
		return "Непонятно"
	}
}

func analyseAnswer(u *StorageUser, c *ViberCallback) string {
	level := u.Level

	if level == 0 {
		return ""
	}
	if level == 1 {
		if c.Message.Text != "да" {
			return "Надо ответить да!"
		}
		return ""
	}
	if level == 2 {
		age, err := analyseAge(c.Message.Text)
		if err != nil {
			return "Вам должно быть 18 или больше"
		}
		u.Age = age
		return ""
	}
	if level == 4 {
		u.Candidate = c.Message.Text
	}
	return ""
}

func analyseAge(age string) (int, error) {
	i, err := strconv.Atoi(age)
	if err != nil {
		return 0, nil
	}
	if i < 18 {
		return 0, errors.New("Вам должно быть 18 или больше")
	}
	return i, nil
}
