package main

import (
	"net/http"
	"log"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/todos", home)
	mux.HandleFunc("/todo/create", createTodo)
	mux.HandleFunc("/todo", getById)
	mux.HandleFunc("/todo/update/", updateTodo)
	mux.HandleFunc("/todo/remove", removeId)

	log.Print("----- server starting on port :4000 -----")

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	} 

}