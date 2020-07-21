package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mileusna/viber"
)

func serve(v *viber.Viber) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleMain(v, w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}
	return nil
}

func handleMain(v *viber.Viber, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	c, err := parseCallback(bytes)
	if err != nil {
		log.Printf("Error reading callback: %v", err)
		http.NotFound(w, r)
		return
	}

	if c.Event == "subscribed" {
		message := "Добрый день,  " + c.User.Name
		v.SendTextMessage(c.User.Id, message)
	} else if c.Event == "message" {
		message := c.User.Name + ", Вы гражданин республики беларусь?"
		v.SendTextMessage(c.User.Id, message)
	}
}
