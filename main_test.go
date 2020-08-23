package main

import (
	"context"
	"entgo-example/ent"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func MockHandler() (Handler, *ent.Client) {
	client, err := ent.Open("sqlite3", "file:ent_test?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Failed opening connection to sqlite: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resource %v", err)
	}

	return NewHandler(client), client
}

func TestHandleUsersFetch(t *testing.T) {
	handler, client := MockHandler()
	defer client.Close()

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	handlerUsersFetch := http.HandlerFunc(handler.handleUsersFetch)
	handlerUsersFetch.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Users Fetch endpoint did not return %v, instead return %v", http.StatusOK, w.Code)
	}
}

func TestHandleUserCreate(t *testing.T) {
	handler, client := MockHandler()
	defer client.Close()

	req, err := http.NewRequest("POST", "/users/create", strings.NewReader(`{"name": "luis", "age": "19"}`))
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	handleUserCreate := http.HandlerFunc(handler.handleUserCreate)
	handleUserCreate.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("User Create endpoint did not return %v, instead return %v", http.StatusCreated, w.Code)
	}
}
