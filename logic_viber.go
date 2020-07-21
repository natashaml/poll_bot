package main

import (
	"errors"
	"strconv"
)

func generateReplyFor(s *Storage, c *ViberCallback) (string, error) {
	if c.Event == "delivered" || c.Event == "seen" {
		return "", nil
	}

	storageUser, err := s.Obtain(c.User.Id)
	if err != nil {
		return "", err
	}

	if storageUser.Candidate != "" {
		return "Вы уже проголосовали за " + storageUser.Candidate, nil
	}

	if c.Event == "subscribed" {
		message := getMessageForLevel(storageUser.Level, c)
		return message, nil
	}
	if c.Event == "message" {
		message := analyseAnswer(storageUser, c)
		if message != "" {
			return message, nil
		}
		if storageUser.Level == 4 {
			return "Спасибо за голосование!", nil
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

func analyseAnswer(u *StorageUser, c *ViberCallback) string {
	level := u.Level

	if level == 0 {
		return ""
	}
	if level == 1 {
		if c.Message.Text != "Да" {
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
