package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
)

func TestGetEmptyTasks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tasks", nil)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getTasks)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "[]"
	jsonassert.New(t).Assertf(rr.Body.String(), expected)
}

func TestGetTasks(t *testing.T) {

	task := Task{
		"uuid",
		"name",
		false,
	}

	tasks = append(tasks, task)

	req, _ := http.NewRequest("GET", "/tasks", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getTasks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":"uuid","name":"name","completed":false}]`

	jsonassert.New(t).Assertf(rr.Body.String(), expected)
}

func TestGetTaskMatch(t *testing.T) {

	task := Task{
		"uuid",
		"name",
		false,
	}

	tasks = append(tasks, task)

	req, _ := http.NewRequest("GET", "/tasks/uuid", nil)

	// fake gorilla/mux vars
	vars := map[string]string{
		"id": "uuid",
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"uuid","name":"name","completed":false}`

	jsonassert.New(t).Assertf(rr.Body.String(), expected)
}

func TestGetTaskNoMatch(t *testing.T) {

	task := Task{
		"uuid",
		"name",
		false,
	}

	tasks = append(tasks, task)

	req, _ := http.NewRequest("GET", "/tasks/wrong", nil)

	// fake gorilla/mux vars
	vars := map[string]string{
		"id": "wrong",
	}

	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestCreateTask(t *testing.T) {
	reader := strings.NewReader("name=taskName")

	req, _ := http.NewRequest("POST", "/tasks", reader)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":"<<PRESENCE>>","name":"taskName","completed":false}`

	jsonassert.New(t).Assertf(rr.Body.String(), expected)
}

func TestCreateTaskMissingName(t *testing.T) {
	req, _ := http.NewRequest("POST", "/tasks", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createTask)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := []byte("name missing in form params\n")

	assert.Equal(t, rr.Body.Bytes(), expected)
}
