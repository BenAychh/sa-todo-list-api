package main

import (
	"fmt"
	"log"
	"os"
)

const (
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbName     = "DB_NAME"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbSslMode  = "DB_SSL_MODE"
	serverPort = "SERVER_PORT"
)

func main() {
	todoApp := TodoApp{}
	todoApp.Initialize(
		getEnvVariable(dbHost),
		getEnvVariable(dbPort),
		getEnvVariable(dbName),
		getEnvVariable(dbUser),
		getEnvVariable(dbPassword),
		getEnvVariable(dbSslMode),
	)
	address := fmt.Sprintf(":%s", getEnvVariable(serverPort))
	todoApp.Start(address)
}

func getEnvVariable(name string) string {
	variable, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Environment variable %s is not set", name)
	}
	return variable
}
