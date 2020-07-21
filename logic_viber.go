package main

func generateReplyFor(c *ViberCallback) error {
	if c.Event == "subscribed" {
		message := "Добрый день,  " + c.User.Name
		return message, nil
	}
	if c.Event == "message" {
		message := c.User.Name + ", Вы гражданин республики беларусь?"
		return message, nil
	}
	return "Я что то вас не понял"
}
