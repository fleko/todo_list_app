package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fleko/todo_list_app/models"
	"github.com/fleko/todo_list_app/responses"
	"github.com/gorilla/mux"
)

// CreateTask parses request, validates data and saves the new task
func (a *App) CreateTask(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Task successfully created"}

	task := &models.Task{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	task.Completed = false

	task.Prepare() // strip away any white spaces

	if err = task.Validate(); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	taskCreated, err := task.Save(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["task"] = taskCreated
	responses.JSON(w, http.StatusCreated, resp)
	return
}

func (a *App) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := models.GetTasks(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, tasks)
	return
}

func (a *App) ToggleTaskStatus(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Task status updated successfully"}

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	task, err := models.GetTaskById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	task.Completed = !task.Completed

	// FIXME: toggle task status only
	_, err = task.UpdateTask(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}

func (a *App) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{"status": "success", "message": "Task deleted successfully"}

	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	_, err := models.GetTaskById(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = models.DeleteTask(id, a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, resp)
	return
}
