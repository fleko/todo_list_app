package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
)

// getTasks is the handler to get all tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("[JSON Encoding Error] %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getTask is the handler to get a specific task by ID
func getTask(w http.ResponseWriter, r *http.Request) {
	// mux.Vars gets a map of path variables by name. here "id" matches the {id} path
	id, ok := mux.Vars(r)["id"]

	if !ok {
		http.Error(w, "id missing in URL path", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			if err := json.NewEncoder(w).Encode(task); err != nil {
				log.Printf("[JSON Encoding Error] %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "no such task", http.StatusNotFound)
}

// createTask is the handler to create a task with a given name (name uniqueness not enforced)
func createTask(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "name missing in form params", http.StatusBadRequest)
		return
	}

	if u, err := uuid.NewV4(); err != nil {
		log.Printf("[UUID generation Error] %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		task := Task{
			u.String(),
			name,
			false,
		}

		tasks = append(tasks, task)

		if err := json.NewEncoder(w).Encode(task); err != nil {
			log.Printf("[JSON Encoding Error] %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// toggleTaskStatus is the handler to toggle the completion status of a task with a specific ID
func toggleTaskStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "ID missing in URL path", http.StatusBadRequest)
		return
	}

	for index, task := range tasks {
		if task.ID == id {

			tasks[index] = Task{
				task.ID,
				task.Name,
				!task.Completed,
			}

			if err := json.NewEncoder(w).Encode(tasks[index]); err != nil {
				log.Printf("[JSON Encoding Error] %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
	}

	http.Error(w, "no such task", http.StatusNotFound)
}

// deleteTask is the handler to delete a specific task by ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "ID missing in URL path", http.StatusBadRequest)
		return
	}

	if index := pos(id); index == -1 {
		http.Error(w, "no such task", http.StatusNotFound)
	} else {
		tasks = remove(index)
		fmt.Fprintf(w, "Success")
	}
}

// removes element from tasks
func remove(index int) []Task {
	return append(tasks[:index], tasks[index+1:]...)
}

// returns positions of task with given ID in tasks, -1 if not found
func pos(id string) int {
	for p, task := range tasks {
		if id == task.ID {
			return p
		}
	}
	return -1
}
