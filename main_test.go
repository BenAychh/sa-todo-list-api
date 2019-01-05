package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	testDbHost     = "TEST_DB_HOST"
	testDbPort     = "TEST_DB_PORT"
	testDbName     = "TEST_DB_NAME"
	testDbUser     = "TEST_DB_USER"
	testDbPassword = "TEST_DB_PASSWORD"
)

var todoApp TodoApp

func TestMain(m *testing.M) {
	todoApp = TodoApp{}
	todoApp.Initialize(
		getEnvVariableWithDefault(testDbHost, "localhost"),
		getEnvVariableWithDefault(testDbPort, "5432"),
		getEnvVariableWithDefault(testDbName, "sa_todo_list_test"),
		getEnvVariableWithDefault(testDbUser, "sa_todo_list_tester"),
		getEnvVariableWithDefault(testDbPassword, "testing123"),
		"disable",
	)

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/v1", nil)
	response := executeRequest(req)

	checkCode(t, http.StatusOK, response.Code)

	checkData(t, []interface{}{}, response.Body.Bytes())
}

func TestCreateTodo(t *testing.T) {
	clearTable()

	payload := []byte(`{"description": "todo one"}`)

	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	expected := map[string]interface{}{}
	expected["id"] = 1
	expected["description"] = "todo one"
	expected["complete"] = false

	checkCode(t, http.StatusCreated, response.Code)

	checkData(t, expected, response.Body.Bytes())
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	todoApp.router.ServeHTTP(recorder, req)
	return recorder
}

func checkCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Fatalf("Expected the response code %d but got %d instead", expected, actual)
	}
}

func checkData(t *testing.T, expected interface{}, actual []byte) {
	var body map[string]interface{}
	json.Unmarshal(actual, &body)

	if body["error"] != nil {
		t.Fatalf("Expected body.error to be nil but got %+v instead", body["error"])
	}

	if body["status"] != "ok" {
		t.Fatalf("Expected body.status to be 'OK' but got %s instead", body["status"])
	}

	expectedBytes, _ := json.Marshal(expected)

	actualBytes, _ := json.Marshal(body["data"])

	if string(actualBytes) != string(expectedBytes) {
		t.Fatalf("Expected body.data to be %+v but got %+v instead", expected, body["data"])
	}
}

func getEnvVariableWithDefault(name string, def string) string {
	variable, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return variable
}

func ensureTableExists() {
	_, err := todoApp.dB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	todoApp.dB.Exec("TRUNCATE todos RESTART IDENTITY;")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS public.todos
(
    id serial NOT NULL,
    description character varying NOT NULL,
    complete boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
)`
