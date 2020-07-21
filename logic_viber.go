package main

import (
	"fmt"
	"strings"
)

func knownNotMessageEvent(c *ViberCallback) bool {
	return c.Event == "delivered" || c.Event == "seen" || c.Event == "subscribed" || c.Event == "conversation_started"
}

type viberReply struct {
	text    string
	options []string
}

func generateReplyFor(p poll, s *Storage, c *ViberCallback) (*viberReply, error) {
	if c.Event == "unsubscribed" {
		return nil, nil
	}

	if c.Event == "delivered" || c.Event == "seen" {
		return nil, nil
	}

	if c.Message.Text == "clear" {
		err := s.Clear(c.User.Id)
		return &viberReply{text: fmt.Sprintf("Your storage cleared with %v", err)}, nil
	}

	storageUser, err := s.Obtain(c.User.Id)
	if err != nil {
		return nil, err
	}

	if c.Event == "message" {
		if storageUser.Candidate != "" {
			return &viberReply{text: "Вы уже проголосовали за " + strings.Title(storageUser.Candidate)}, nil
		}

		err := analyseAnswer(p, storageUser, c)
		if err != nil {
			reply := getViberReplyForLevel(p, storageUser.Level, c)
			reply.text = err.Error() + " " + reply.text
			return reply, nil
		}
		if storageUser.Level == len(p)+1 {
			_ = s.Persist(storageUser.Id)
			if err != nil {
				return nil, err
			}

			return &viberReply{text: "Спасибо за голосование!"}, nil
		}

		storageUser.Level++
		_ = s.Persist(storageUser.Id)
		if err != nil {
			return nil, err
		}
		reply := getViberReplyForLevel(p, storageUser.Level, c)
		return reply, nil
	}

	if !storageUser.ConversationStarted {
		reply := getViberReplyForLevel(p, storageUser.Level, c)
		storageUser.ConversationStarted = true
		_ = s.Persist(storageUser.Id)
		if err != nil {
			return nil, err
		}
		return reply, nil
	}

	if knownNotMessageEvent(c) {
		return nil, nil
	}

	return nil, fmt.Errorf("Unknown message %v", c.Event)
}

func getViberReplyForLevel(p poll, level int, c *ViberCallback) *viberReply {
	item := p[level]
	reply := viberReply{text: "Непонятно"}
	if item != nil {
		reply.text = item.question(c)
		reply.options = item.possibleAnswers
	}
	return &reply
}

func analyseAnswer(p poll, u *StorageUser, c *ViberCallback) error {
	level := u.Level
	item := p[level]
	if item != nil && item.validateAnswer != nil {
		err := item.validateAnswer(c.Message.Text)
		if err != nil {
			return err
		}
	}
	if item != nil && item.persistAnswer != nil {
		err := item.persistAnswer(c.Message.Text, u)
		if err != nil {
			return err
		}
	}
	return nil
}
