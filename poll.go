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
			if c.User.Name == "" {
				return "Добро пожаловать."
			}
			return "Добрый день, " + c.User.Name + ". Добро пожаловать"
		},
	}
	ret[item0.level] = &item0

	item1 := pollItem{
		level: 1,
		question: func(c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваше гражданство?"
		},
		possibleAnswers: []string{"Беларусь", "Россия", "Украина", "Казахстан", "Другая страна"},
		validateAnswer: func(answer string) error {
			if answer != "беларусь" {
				return errors.New("Только граждание Беларуси могут принимать участие!")
			}
			return nil
		},
	}
	ret[item1.level] = &item1

	item2 := pollItem{
		level: 2,
		question: func(c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш возраст"
		},
		possibleAnswers: []string{"18-24", "25-34", "35-44", "45-54", "55+"},
		validateAnswer: func(answer string) error {
			i, err := strconv.Atoi(answer)
			if err != nil {
				return nil
			}
			if i < 18 {
				return errors.New("Вам должно быть 18 или больше.")
			}
			return nil
		},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Age = answer
			return nil
		},
	}
	ret[item2.level] = &item2

	item3 := pollItem{
		level: 3,
		question: func(c *ViberCallback) string {
			return "Примете ли Вы участие в предстоящих выборах Президента?"
		},
		possibleAnswers: []string{"Да, приму обязательно", "Да, скорее приму", "Нет, скорее не приму", "Нет, не приму", "Затрудняюсь ответить"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.WillTakePart = answer
			return nil
		},
	}
	ret[item3.level] = &item3

	item4 := pollItem{
		level: 4,
		question: func(c *ViberCallback) string {
			return "Когда Вы планируете голосовать?"
		},
		possibleAnswers: []string{"Досрочно", "В основной день", "Затрудняюсь ответить"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.WhenVote = answer
			return nil
		},
	}
	ret[item4.level] = &item4

	item5 := pollItem{
		level: 5,
		question: func(c *ViberCallback) string {
			return "За кого Вы планируете проголосовать?"
		},
		possibleAnswers: []string{"Александр Лукашенко", "Сергей Черечень", "Анна Канопацкая", "Андрей Дмитриев", "Светлана Тихановская", "Против всех", "Затрудняюсь ответить"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Candidate = answer
			return nil
		},
	}
	ret[item5.level] = &item5

	item6 := pollItem{
		level: 6,
		question: func(c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш пол"
		},
		possibleAnswers: []string{"Мужской", "Женский"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Gender = answer
			return nil
		},
	}
	ret[item6.level] = &item6

	item7 := pollItem{
		level: 7,
		question: func(c *ViberCallback) string {
			return "Ваш уровень образования?"
		},
		possibleAnswers: []string{"Среднее общее (школа)", "Профессионально-техническое", "Среднее специальное (техникум, колледж)", "Неполное высшее", "Высшее, ученая степень"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Education = answer
			return nil
		},
	}
	ret[item7.level] = &item7

	return ret
}
