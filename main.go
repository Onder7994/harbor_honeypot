package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Onder7994/harbor_honeypot/internal/handlers"
	"github.com/Onder7994/harbor_honeypot/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	middleware.InitLogger()
	defer middleware.CloseLogger()

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		fmt.Println("APP_PORT variable not set, using default port 8080")
		appPort = "8080"
	}
	addr := fmt.Sprintf(":%s", appPort)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/account/sign-in", http.StatusFound)
			return
		}
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	})
	mux.HandleFunc("/account/sign-in", handlers.AccountHandler)
	mux.HandleFunc("/c/login", handlers.LoginPostHandler)

	loggedMux := middleware.LoggingMiddleware(mux)
	fmt.Printf("Honeypot server running at %s port\n", appPort)
	log.Fatal(http.ListenAndServe(addr, loggedMux))
}
