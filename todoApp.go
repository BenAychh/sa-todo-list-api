package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
)

// TodoApp is the Todo API
type TodoApp struct {
	router *chi.Mux
	dB     *sql.DB
}

// Initialize must be called before start to setup the database
func (app *TodoApp) Initialize(host, port, dbname, user, password, sslmode string) {
	connectionInformation := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

	var err error

	app.dB, err = sql.Open("postgres", connectionInformation)
	if err != nil {
		panic(err)
	}
	err = app.dB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to DB")

	app.router = chi.NewRouter()
	app.router.Use(render.SetContentType(render.ContentTypeJSON))
	app.setupRoutes()
}

// Start the actual todo API
func (app *TodoApp) Start(address string) {
	log.Println("Starting HTTP Server")
	log.Fatal(http.ListenAndServe(address, app.router))
}

func (app *TodoApp) setupRoutes() {
	app.router.Route("/v1", func(router chi.Router) {
		router.Get("/", app.getAllTodos)
	})
}

func (app *TodoApp) getAllTodos(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, 200, []string{})
}

type response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

func sendJSON(w http.ResponseWriter, code int, payload interface{}) {
	r := response{
		Status: "ok",
		Data:   payload,
		Error:  nil,
	}
	jsonResponse, error := json.Marshal(r)

	if error != nil {
		fmt.Println("Error converting to json")
		fmt.Println(error)
		sendError(w, http.StatusInternalServerError, error)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, error = w.Write(jsonResponse)
		if error != nil {
			fmt.Println("Error sending json")
			fmt.Println(error)
		}
	}
}

func sendError(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	r := response{
		Status: "error",
		Data:   nil,
		Error:  payload,
	}
	jsonResponse, error := json.Marshal(r)

	if error != nil {
		fmt.Println("Error converting error to json")
		fmt.Println(error)
		r = response{
			Status: "error",
			Data:   nil,
			Error:  nil,
		}
		jsonResponse, error = json.Marshal(r)
		if error != nil {
			fmt.Println("Error converting basic response to json")
			fmt.Println(error)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, error = w.Write(jsonResponse)
	if error != nil {
		fmt.Println("Error sending error")
		fmt.Println(error)
	}
}
