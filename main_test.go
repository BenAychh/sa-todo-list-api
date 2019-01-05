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
	// Setup
	clearTable()

	// Execute
	req, _ := http.NewRequest("GET", "/v1", nil)
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusOK, response.Code)
	checkData(t, []interface{}{}, response.Body.Bytes())
}

func TestCreateTodo(t *testing.T) {
	// Setup
	clearTable()

	// Execute
	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	expected := todoMap(1, "todo one", false)

	checkCode(t, http.StatusCreated, response.Code)
	checkData(t, expected, response.Body.Bytes())
}

func TestListOfTodos(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	payload = []byte(`{"description": "todo two"}`)
	req, _ = http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	req, _ = http.NewRequest("GET", "/v1", nil)
	response := executeRequest(req)

	// Test
	expected := []map[string]interface{}{
		todoMap(1, "todo one", false),
		todoMap(2, "todo two", false),
	}

	checkCode(t, http.StatusOK, response.Code)
	checkData(t, expected, response.Body.Bytes())
}

func TestUpdateCompleted(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	payload = []byte(`{"complete": "true"}`)
	req, _ = http.NewRequest("PATCH", "/v1/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	expected := todoMap(1, "todo one", true)

	checkCode(t, http.StatusOK, response.Code)
	checkData(t, expected, response.Body.Bytes())
}

func TestUpdateDescription(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	payload = []byte(`{"description": "new description"}`)
	req, _ = http.NewRequest("PATCH", "/v1/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	expected := todoMap(1, "new description", false)

	checkCode(t, http.StatusOK, response.Code)
	checkData(t, expected, response.Body.Bytes())
}

func TestDeletionCorrectResponse(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	payload = []byte(`{"description": "todo two"}`)
	req, _ = http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	req, _ = http.NewRequest("DELETE", "/v1/1", nil)
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusNoContent, response.Code)
	checkData(t, nil, response.Body.Bytes())
}

func TestDeletionActuallyDelete(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	payload = []byte(`{"description": "todo two"}`)
	req, _ = http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	req, _ = http.NewRequest("DELETE", "/v1/1", nil)
	_ = executeRequest(req)

	// Execute
	req, _ = http.NewRequest("GET", "/v1", nil)
	response := executeRequest(req)

	// Test
	expected := []map[string]interface{}{
		todoMap(2, "todo two", false),
	}

	// Test
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

func todoMap(id int, description string, complete bool) map[string]interface{} {
	m := map[string]interface{}{}
	m["id"] = id
	m["description"] = description
	m["complete"] = complete
	return m
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS public.todos
(
    id serial NOT NULL,
    description character varying NOT NULL,
    complete boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
)`
