package ui

import (
	"log"
	"net/http"
)

func (ui *UI) Users(w http.ResponseWriter, _ *http.Request) {
	users, err := ui.src.Query()

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		if _, err := w.Write([]byte(user.Acct + "\n")); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
