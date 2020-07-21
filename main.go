package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mileusna/viber"
)

func main() {

	<-make(chan int)
}

func execute() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	viberKey := os.Getenv("VIBER_KEY")
	callbackURL := os.Getenv("CALLBACK_URL")
	v := viber.New(viberKey, "Voting bot", "https://thumbs.dreamstime.com/z/human-hand-write-yes-vote-voting-paper-pen-flat-concept-illustration-man-s-red-pen-ballot-check-sign-88802664.jpg")
	go serve(v)
	_, err = v.SetWebhook(callbackURL, nil)
	return err
}
