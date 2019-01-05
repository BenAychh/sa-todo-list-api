package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
