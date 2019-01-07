package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
)

// TodoApp is the Todo API
type TodoApp struct {
	router *chi.Mux
	dB     *sql.DB
}

// Initialize must be called before start to setup the database
func (app *TodoApp) Initialize(host, port, dbname, user, password, sslmode string, corsOrigins []string) {
	connectionInformation := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

	var err error

	app.dB, err = sql.Open("postgres", connectionInformation)

	if err != nil {
		panic(err)
	}

	err = app.dB.Ping()

	for i := 1; i <= 10 && err != nil; i++ {
		fmt.Printf("Error connecting to database try %d of 10\n", i)
		if i == 10 {
			fmt.Println("Could not connect to database after 10 tries, giving up")
			panic(err)
		} else {
			time.Sleep(3 * time.Second)
		}
		err = app.dB.Ping()
	}

	_, err = app.dB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB")

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	app.router = chi.NewRouter()
	app.router.Use(cors.Handler)
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
		router.Post("/", app.createTodo)
		router.Patch("/{todoID}", app.patchTodo)
		router.Delete("/{todoID}", app.deleteTodo)
	})
}

func (app *TodoApp) getAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := getTodos(app.dB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Println("Error getting todos")
		log.Println(err)
	} else {
		sendJSON(w, 200, todos)
	}
}

func (app *TodoApp) createTodo(w http.ResponseWriter, r *http.Request) {
	var t todo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Malformed payload")
		return
	}

	err = t.create(app.dB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, t)

	closeError := r.Body.Close()
	if closeError != nil {
		log.Println("Error closing body in createTodo")
		log.Println(closeError)
	}
}

func (app *TodoApp) patchTodo(w http.ResponseWriter, r *http.Request) {
	todoIDString := chi.URLParam(r, "todoID")
	todoID, err := strconv.Atoi(todoIDString)
	if err != nil {
		sendError(w, http.StatusNotFound, "Not Found")
		return
	}

	t := todo{ID: todoID}
	err = t.get(app.dB)

	if err != nil {
		sendError(w, http.StatusNotFound, "Not Found")
		return
	}

	payload := map[string]string{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Malformed payload")
		return
	}

	mergeErrorDetails := mergeMapAndTodo(&t, payload)

	if mergeErrorDetails != nil {
		sendError(w, mergeErrorDetails.code, mergeErrorDetails.message)
		return
	}

	err = t.update(app.dB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, t)

	closeError := r.Body.Close()
	if closeError != nil {
		log.Println("Error closing body in createTodo")
		log.Println(closeError)
	}
}

func (app *TodoApp) deleteTodo(w http.ResponseWriter, r *http.Request) {
	todoIDString := chi.URLParam(r, "todoID")
	todoID, err := strconv.Atoi(todoIDString)
	if err != nil {
		sendError(w, http.StatusNotFound, "Not Found")
		return
	}

	t := todo{ID: todoID}
	err = t.delete(app.dB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendJSON(w, http.StatusOK, nil)
}

type errorDetails struct {
	code    int
	message string
}

func mergeMapAndTodo(t *todo, m map[string]string) *errorDetails {
	for key, value := range m {
		switch key {
		case "description":
			t.Description = value
		case "complete":
			b, bErr := strconv.ParseBool(value)
			if bErr != nil {
				return &errorDetails{
					code:    http.StatusBadRequest,
					message: "Invalid value for complete, must be true or false",
				}
			}
			t.Complete = b
		default:
			return &errorDetails{
				code:    http.StatusBadRequest,
				message: "Only the keys description and complete are allowed",
			}
		}
	}
	return nil
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS public.todos
(
    id serial NOT NULL,
    description character varying NOT NULL,
    complete boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
)`
