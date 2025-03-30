package main

import (
	"app/be/internal/repositories"
	"app/be/internal/srv/user"
	"app/be/pkg/database"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	client, err := database.ConnectMongo("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())
	userCollection := database.GetCollection("bankingSystem", "users")
	userRepository := repositories.NewUserRepository(userCollection)
	authService := user.NewAuthService(userRepository)
	authHandler := user.NewAuthHandler(authService)

	err = userRepository.DeleteAllUsers(context.Background())
	if err != nil {
		log.Printf("Error clearing database: %v", err)
	}

	fs := http.FileServer(http.Dir("fe"))
	http.Handle("/fe/", http.StripPrefix("/fe/", fs))

	// API endpoints
	http.HandleFunc("/api/register", authHandler.RegisterHandler)
	http.HandleFunc("/api/login", authHandler.LoginHandler)
	http.HandleFunc("/api/logout", authHandler.Logout)
	http.HandleFunc("/api/debug/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepository.GetAllUsers(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	// Page handlers
	http.HandleFunc("/", serveTemplate("fe/html/home.html"))
	http.HandleFunc("/login", serveTemplate("fe/html/login.html"))
	http.HandleFunc("/create-account", serveTemplate("fe/html/createAcc.html"))

	log.Println("Server is running on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func serveTemplate(tmplPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		absPath, err := filepath.Abs(tmplPath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(absPath)
		if err != nil {
			http.Error(w, "Template not found", http.StatusNotFound)
			return
		}

		data := map[string]interface{}{
			"Title": "Banking System",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}
