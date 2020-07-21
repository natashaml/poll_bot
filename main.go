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

	s3Bucket := os.Getenv("VIBER_KEY")
	callback_URL := os.Getenv("CALLBACK_URL")
	v := viber.New(callback_URL, "Voting bot", "https://thumbs.dreamstime.com/z/human-hand-write-yes-vote-voting-paper-pen-flat-concept-illustration-man-s-red-pen-ballot-check-sign-88802664.jpg")
	go serve(v)
	v.SetWebhook("callback_URL, nil)
}
