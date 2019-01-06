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

	testEnsureTableExists()

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

func TestCreateTodoBadPayload(t *testing.T) {
	// Setup
	clearTable()

	// Execute
	payload := []byte(`"description"="todo one"`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkCode(t, http.StatusBadRequest, response.Code)
	checkError(t, "Malformed payload", response.Body.Bytes())
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

func TestUpdateTodoNotFound(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	payload = []byte(`{"complete": "true"}`)
	req, _ = http.NewRequest("PATCH", "/v1/2", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusNotFound, response.Code)
	checkError(t, "Not Found", response.Body.Bytes())
}

func TestUpdateTodoNonNumericId(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	payload = []byte(`{"complete": "true"}`)
	req, _ = http.NewRequest("PATCH", "/v1/abc", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusNotFound, response.Code)
	checkError(t, "Not Found", response.Body.Bytes())
}

func TestUpdateTodoMalformedJSON(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	payload = []byte(`"complete"="true"`)
	req, _ = http.NewRequest("PATCH", "/v1/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusBadRequest, response.Code)
	checkError(t, "Malformed payload", response.Body.Bytes())
}

func TestUpdateTodoBadKey(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	payload = []byte(`{"BadKey": "abcdef"}`)
	req, _ = http.NewRequest("PATCH", "/v1/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusBadRequest, response.Code)
	checkError(t, "Only the keys description and complete are allowed", response.Body.Bytes())
}

func TestUpdateTodoNonBooleanComplete(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	payload = []byte(`{"complete": "abcdef"}`)
	req, _ = http.NewRequest("PATCH", "/v1/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusBadRequest, response.Code)
	checkError(t, "Invalid value for complete, must be true or false", response.Body.Bytes())
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
	checkCode(t, http.StatusOK, response.Code)
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

func TestDeletionNonNumericId(t *testing.T) {
	// Setup
	clearTable()

	payload := []byte(`{"description": "todo one"}`)
	req, _ := http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	payload = []byte(`{"description": "todo two"}`)
	req, _ = http.NewRequest("POST", "/v1", bytes.NewBuffer(payload))
	_ = executeRequest(req)

	// Execute
	req, _ = http.NewRequest("DELETE", "/v1/abc", nil)
	response := executeRequest(req)

	// Test
	checkCode(t, http.StatusNotFound, response.Code)
	checkError(t, "Not Found", response.Body.Bytes())
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
		t.Fatalf("Expected body.status to be 'ok' but got %s instead", body["status"])
	}

	expectedBytes, _ := json.Marshal(expected)

	actualBytes, _ := json.Marshal(body["data"])

	if string(actualBytes) != string(expectedBytes) {
		t.Fatalf("Expected body.data to be %+v but got %+v instead", expected, body["data"])
	}
}

func checkError(t *testing.T, expected string, actual []byte) {
	var body map[string]interface{}
	json.Unmarshal(actual, &body)

	if body["data"] != nil {
		t.Fatalf("Expected body.data to be nil but got %+v instead", body["error"])
	}

	if body["status"] != "error" {
		t.Fatalf("Expected body.status to be 'error' but got %s instead", body["status"])
	}

	if body["error"] != expected {
		t.Fatalf("Expected body.error to be %+v but got %+v instead", expected, body["error"])
	}
}

func getEnvVariableWithDefault(name string, def string) string {
	variable, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return variable
}

func testEnsureTableExists() {
	_, err := todoApp.dB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	todoApp.dB.Exec("TRUNCATE todos RESTART IDENTITY;")
}

func todoMap(id int, description string, complete bool) map[string]interface{} {
	return map[string]interface{}{
		"id":          id,
		"description": description,
		"complete":    complete,
	}
}
