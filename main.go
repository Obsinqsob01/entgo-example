package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"entgo-example/ent"

	_ "github.com/mattn/go-sqlite3"
)

const (
	PORT = 8000
)

// userBody to handle post requests, it is better than encode all body into User struct
type userBody struct {
	Name string `json:"name"`
	Age  int    `json:"age,string"`
}

type Handler struct {
	client *ent.Client
}

func NewHandler(client *ent.Client) Handler {
	return Handler{client}
}

func (h *Handler) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	userBody := &userBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.client.User.Create().SetName(userBody.Name).SetAge(userBody.Age).Save(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJSON)
}

func (h *Handler) handleUsersFetch(w http.ResponseWriter, r *http.Request) {
	users, err := h.client.User.Query().All(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(usersJSON)
}

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resource %v", err)
	}

	handler := NewHandler(client)

	http.HandleFunc("/users", handler.handleUsersFetch)
	http.HandleFunc("/users/create", handler.handleUserCreate)

	log.Printf("Server running at %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
