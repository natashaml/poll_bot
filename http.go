package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/mileusna/viber"
)

func serve(v *viber.Viber) error {
	s := newStorage()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleMain(v, s, w, r)
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

func handleMain(v *viber.Viber, s *Storage, w http.ResponseWriter, r *http.Request) {
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
	log.Printf("Request body: %v", string(bytes))

	c, err := parseCallback(bytes)
	if err != nil {
		log.Printf("Error reading callback: %v", err)
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	reply, err := generateReplyFor(s, c)
	if err != nil {
		log.Printf("Error generating reply: %v", err)
		http.Error(w, "can't reply", http.StatusBadRequest)
		return
	}
	v.SendTextMessage(c.User.Id, reply)
}
