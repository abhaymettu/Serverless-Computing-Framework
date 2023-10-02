package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"github.com/gorilla/mux"
)

// Function represents a piece of code to be executed
type Function struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

// InMemoryStore is an in-memory storage for our functions
type InMemoryStore struct {
	mu       sync.Mutex
	functions map[string]Function
}

var store = &InMemoryStore{
	functions: make(map[string]Function),
}

func (s *InMemoryStore) Save(f Function) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.functions[f.ID] = f
}

func (s *InMemoryStore) Get(id string) (Function, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f, ok := s.functions[id]; ok {
		return f, nil
	}
	return Function{}, errors.New("function not found")
}

func createFunction(w http.ResponseWriter, r *http.Request) {
	var f Function
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.Save(f)
	w.WriteHeader(http.StatusCreated)
}

func invokeFunction(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Function ID is required", http.StatusBadRequest)
		return
	}

	f, err := store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Simulating function execution
	fmt.Fprintf(w, "Executing function with ID %s: %s", f.ID, f.Code)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/function", createFunction).Methods("POST")
	r.HandleFunc("/function/invoke", invokeFunction).Methods("POST")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
