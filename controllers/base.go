package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fleko/todo_list_app/middlewares"
	"github.com/fleko/todo_list_app/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize connect to the database and wire up routes
func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error
	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	a.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Printf("\n Cannot connect to database %s", DbName)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the database %s", DbName)
	}

	a.DB.Debug().AutoMigrate(&models.Task{}, &models.User{}) //database migration

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middlewares.SetContentTypeMiddleware) // setting content-type to json

	// configure the router to always run this handler when it couldn't match a request to any other handler
	a.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("%s not found\n", r.URL)))
	})

	a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")

	a.Router.HandleFunc("/api/tasks", a.GetTasks).Methods("GET")

	a.Router.HandleFunc("/api/tasks", a.CreateTask).Methods("POST")

	a.Router.HandleFunc("/api/tasks/{id:[0-9]+}", a.DeleteTask).Methods("DELETE")

	a.Router.HandleFunc("/api/tasks/{id:[0-9]+}", a.ToggleTaskStatus).Methods("PUT")
}

func (a *App) RunServer() {
	log.Printf("\nServer starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
