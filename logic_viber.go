package main

import (
	"errors"
	"fmt"
)

func generateReplyFor(s *Storage, c *ViberCallback) (string, error) {
	storageUser, err := s.Obtain(c.User.Id)
	if err != nil {
		return "", err
	}

	if c.Event == "subscribed" {
		message := getMessageForLevel(storageUser.Level, c)
		return message, nil
	}
	if c.Event == "message" {
		message := analyseAnswer(storageUser.Level, c)
		if message != "" {
			return message, nil
		}
		storageUser.Level++
		message = getMessageForLevel(storageUser.Level, c)
		return message, nil
	}
	return "", errors.New("Unknown message")
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

func analyseAnswer(level int, c *ViberCallback) string {
	if level == 0 {
		return ""
	}
	fmt.Println(c.Message.Text)
	if level == 1 && c.Message.Text != "Да" {
		return "Надо ответить да!"
	}
	if level == 2 && c.Message.Text != "Да" {
		return "Надо больше 18!"
	}
	return "Непоятно"
}
