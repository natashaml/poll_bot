package main

import (
	"errors"
	"strconv"
)

type pollItem struct {
	level           int
	question        func(c *ViberCallback) string
	possibleAnswers []string
	validateAnswer  func(string) error
	persistAnswer   func(string, *StorageUser) error
}

type poll map[int]*pollItem

func generateOurPoll() poll {
	ret := poll{}
	item0 := pollItem{
		level: 0,
		question: func(c *ViberCallback) string {
			return "Добрый день, " + c.User.Name + ". Добро пожаловать"
		},
	}
	ret[item0.level] = &item0

	item1 := pollItem{
		level: 1,
		question: func(c *ViberCallback) string {
			return "Вы гражданин республики Беларусь?"
		},
		possibleAnswers: []string{"Да", "Нет"},
		validateAnswer: func(answer string) error {
			if answer != "да" {
				return errors.New("Надо ответить да!")
			}
			return nil
		},
	}
	ret[item1.level] = &item1

	item2 := pollItem{
		level: 2,
		question: func(c *ViberCallback) string {
			return "Ваш возраст?"
		},
		possibleAnswers: []string{"18", "19", "20"},
		validateAnswer: func(answer string) error {
			i, err := strconv.Atoi(answer)
			if err != nil {
				return err
			}
			if i < 18 {
				return errors.New("Вам должно быть 18 или больше")
			}
			return nil
		},
		persistAnswer: func(answer string, u *StorageUser) error {
			i, err := strconv.Atoi(answer)
			if err != nil {
				return err
			}
			u.Age = i
			return nil
		},
	}
	ret[item2.level] = &item2

	item3 := pollItem{
		level: 3,
		question: func(c *ViberCallback) string {
			return "Укажите вас регион?"
		},
	}
	ret[item3.level] = &item3

	item4 := pollItem{
		level: 4,
		question: func(c *ViberCallback) string {
			return "Какой ваш кандидат?"
		},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Candidate = answer
			return nil
		},
	}
	ret[item4.level] = &item4

	return ret
}
