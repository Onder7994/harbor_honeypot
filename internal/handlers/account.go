package handlers

import (
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
)

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, err := r.Cookie("fingerprint")
		if err != nil {
			fingerprint := rand.Intn(900_000_000) + 100_000_000
			http.SetCookie(w, &http.Cookie{
				Name:     "fingerprint",
				Value:    strconv.Itoa(fingerprint),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   3600 * 24 * 365, // one year
			})
		}

		tmpl, err := template.ParseFiles("static/index.html")

		if err != nil {
			http.Error(w, "HTML template error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
