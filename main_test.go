package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetTasks(t *testing.T) {
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

	body := []Task{}
	err := json.Unmarshal(rr.Body.Bytes(), &body)

	if err != nil {
		fmt.Println(err.Error())
	}

	expected := []Task{}

	if !reflect.DeepEqual(expected, body) {
		t.Errorf("Expected %v, got %v", expected, body)
	}
}
