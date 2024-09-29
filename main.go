package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go-sqlc-sqlite-crud/db"
)

var store *db.Queries

// Create users table if not exists
func createTable(conn *sql.DB) {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );
    `
	_, err := conn.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Table 'users' exists or was created successfully.")
}

func main() {
	conn, err := sql.Open("sqlite3", "mydb.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create users table if it doesn't exist
	createTable(conn)

	store = db.New(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", usersHandler)
	mux.HandleFunc("/users/", userHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server running on port 8080")
	log.Fatal(server.ListenAndServe())
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch r.Method {
	case http.MethodGet:
		readUsers(w, ctx)
	case http.MethodPost:
		createUser(w, r, ctx)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := strconv.ParseInt(r.URL.Path[len("/users/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		readUserByID(w, ctx, id)
	case http.MethodPut:
		updateUser(w, r, ctx, id)
	case http.MethodDelete:
		deleteUser(w, ctx, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func readUsers(w http.ResponseWriter, ctx context.Context) {
	users, err := store.GetAllUsers(ctx)
	if err != nil {
		http.Error(w, "Could not get users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func readUserByID(w http.ResponseWriter, ctx context.Context, id int64) {
	user, err := store.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Could not get user", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request, ctx context.Context, id int64) {
	var params db.UpdateUserByIDParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	params.ID = id

	err := store.UpdateUserByID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Could not update user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, ctx context.Context, id int64) {
	err := store.DeleteUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Could not delete user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createUser(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var params db.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	if params.Name == "" || params.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	err := store.CreateUser(ctx, params)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, fmt.Sprintf("Could not create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
