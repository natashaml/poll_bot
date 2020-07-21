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

func generateReplyFor(p poll, s *Storage, c *ViberCallback) (string, error) {
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

		message := analyseAnswer(p, storageUser, c)
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
		message = getMessageForLevel(p, storageUser.Level, c)
		return message, nil
	}

	if !storageUser.ConversationStarted {
		message := getMessageForLevel(p, storageUser.Level, c)
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

func getMessageForLevel(p poll, level int, c *ViberCallback) string {
	item := p[level]
	if item != nil {
		return item.question(c)
	}
	return "Непонятно"
}

func analyseAnswer(p poll, u *StorageUser, c *ViberCallback) string {
	level := u.Level
	item := p[level]
	if item != nil && item.validateAnswer != nil {
		err := item.validateAnswer(c.Message.Text)
		if err != nil {
			return err.Error()
		}
	}
	if item != nil && item.persistAnswer != nil {
		err := item.persistAnswer(c.Message.Text, u)
		if err != nil {
			return err.Error()
		}
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
