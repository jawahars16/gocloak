package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jawahars16/gocloak/config"
	"github.com/jawahars16/gocloak/infra/db"
	"github.com/jawahars16/gocloak/infra/token"
	"github.com/jawahars16/gocloak/internal/auth"
	"github.com/jawahars16/gocloak/internal/user"
)

func main() {
	cfg := config.Load()
	db, err := db.New(cfg.DB)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to DB: %v", err))
	}

	crypto := auth.NewCrypto()
	tokenGenerator := token.NewGenerator()

	userManager := user.NewManager(&db, &crypto)
	authManager := auth.NewManager(&db, &crypto, &tokenGenerator, cfg.Auth)

	userHandler := user.NewHandler(&userManager)
	authHandler := auth.NewHandler(&authManager)

	r := mux.NewRouter()
	r.HandleFunc("/user", userHandler.AddUser).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	fmt.Println("Starting server on port 3000")
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
