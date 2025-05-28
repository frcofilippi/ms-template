package main

type DatabaseConfiguration struct {
	connectionStr string
}

func NewDatabaseConfig() *DatabaseConfiguration {
	return &DatabaseConfiguration{
		connectionStr: readEnvValueAsString("DB_CONNECTION_STR", "postgres://pedimeapp:mysecretpwd@localhost/pedimedb?sslmode=disable"),
	}
}
