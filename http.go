package main

import (
	"encoding/json"
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

func isJSON(s []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil

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
	if !isJSON(bytes) {
		http.Error(w, "Not json response", http.StatusBadRequest)
		return
	}

	log.Printf("Request body: %v", string(bytes))

	c, err := parseCallback(bytes)
	if err != nil {
		log.Printf("Error reading callback: %v", err)
		http.Error(w, "can't parse body", http.StatusBadRequest)
		return
	}

	if c.Event == "delivered" || c.Event == "seen" {

	}

	reply, err := generateReplyFor(s, c)
	if err != nil {
		log.Printf("Error generating reply: %v", err)
		http.Error(w, "can't reply", http.StatusBadRequest)
		return
	}
	if reply != "" {
		v.SendTextMessage(c.User.Id, reply)
	}
}
