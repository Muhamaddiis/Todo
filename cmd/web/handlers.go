package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type Todo struct {
	Title string `json:"title"`
	ID int	`json:"id"`
	Completed bool `json:"completed"`
}

var todos = []Todo{
    {ID: 1, Title: "Finish homework", Completed: false},
    {ID: 2, Title: "Go to the gym", Completed: true},
    {ID: 3, Title: "Buy groceries", Completed: false},
}

var mu sync.Mutex


func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/todos" {
		http.NotFound(w, r)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	newTodo.ID = len(todos) + 1
	todos = append(todos, newTodo)
}

func getById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	var foundTodo *Todo

	for _, todo := range todos {
		if todo.ID == id {
			foundTodo = &todo
			break
		}
	}

	if foundTodo == nil {
        http.NotFound(w, r)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundTodo)
}

func removeId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Metthod Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	var index = -1
	for i, todo := range todos {
		if todo.ID == id {
			index = i
			break
		}
	}
	todos = append(todos[:index], todos[index + 1:]...)
	w.Write([]byte("todo deleted successfully"))
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i, t := range todos {
		if t.ID == id {
			todos[i].Completed = true
			break
		}
	} 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}