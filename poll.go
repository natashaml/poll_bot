package main

import (
	"errors"
)

type pollItem struct {
	level           int
	question        func(user *StorageUser, c *ViberCallback) string
	possibleAnswers []string
	validateAnswer  func(string) error
	persistAnswer   func(string, *StorageUser) error
}

type poll struct {
	items map[int]*pollItem
	size  int
}

func (p *poll) add(item *pollItem) {
	item.level = p.size
	p.items[item.level] = item
	p.size++
}

func (p *poll) getLevel(level int) *pollItem {
	return p.items[level]
}

func generateOurPoll() poll {
	ret := poll{
		items: map[int]*pollItem{},
	}

	ret.add(&pollItem{
		question: func(user *StorageUser, c *ViberCallback) string {
			var welcome string
			if user.Properties["ConversationStarted"] != "true" {
				if c.User.Name == "" {
					welcome = "Добро пожаловать. "
				} else {
					welcome = "Добрый день, " + c.User.Name + ". Добро пожаловать. "
				}
			}

			return welcome + "Укажите, пожалуйста, Ваше гражданство?"
		},
		possibleAnswers: []string{"Беларусь", "Россия", "Украина", "Казахстан", "Другая страна"},
		validateAnswer: func(answer string) error {
			if answer != "Беларусь" {
				return errors.New("Только граждание Беларуси могут принимать участие!")
			}
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *StorageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш возраст"
		},
		possibleAnswers: []string{"меньше 18", "18-24", "25-34", "35-44", "45-54", "55+"},
		validateAnswer: func(answer string) error {
			if answer == "меньше 18" {
				return errors.New("Вам должно быть 18 или больше.")
			}
			return nil
		},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Properties["age"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *StorageUser, c *ViberCallback) string {
			return "Примете ли Вы участие в предстоящих выборах Президента?"
		},
		possibleAnswers: []string{"Да, приму обязательно", "Да, скорее приму", "Нет, скорее не приму", "Нет, не приму", "Затрудняюсь ответить"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Properties["will_take_part"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *StorageUser, c *ViberCallback) string {
			return "Когда Вы планируете голосовать?"
		},
		possibleAnswers: []string{"Досрочно", "В основной день", "Затрудняюсь ответить"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Properties["when_vote"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *StorageUser, c *ViberCallback) string {
			return "За кого Вы планируете проголосовать?"
		},
		possibleAnswers: []string{"Александр Лукашенко", "Сергей Черечень", "Анна Канопацкая", "Андрей Дмитриев", "Светлана Тихановская", "Против всех", "Затрудняюсь ответить"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Candidate = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *StorageUser, c *ViberCallback) string {
			return "Укажите, пожалуйста, Ваш пол"
		},
		possibleAnswers: []string{"Мужской", "Женский"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Properties["gender"] = answer
			u.isChanged = true
			return nil
		},
	})

	ret.add(&pollItem{
		question: func(user *StorageUser, c *ViberCallback) string {
			return "Ваш уровень образования?"
		},
		possibleAnswers: []string{"Среднее общее (школа)", "Профессионально-техническое", "Среднее специальное (техникум, колледж)", "Неполное высшее", "Высшее, ученая степень"},
		persistAnswer: func(answer string, u *StorageUser) error {
			u.Properties["education"] = answer
			u.isChanged = true
			return nil
		},
	})

	return ret
}
