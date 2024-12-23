package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	CACHE_ON = "1"
)

// Container contains environment variables for the application, database, cache, token, and http server
type (
	Container struct {
		App   *App
		Redis *Redis
		DB    *DB
		HTTP  *HTTP
	}
	// App contains all the environment variables for the application
	App struct {
		Name      string
		Env       string
		Port      string
		IsCacheOn string
	}

	// Redis contains all the environment variables for the cache service
	Redis struct {
		Addr     string
		Password string
	}
	// Database contains all the environment variables for the database
	DB struct {
		Connection  string
		Host        string
		Port        string
		User        string
		Password    string
		ClusterName string
		DBName      string
	}
	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}
)

// New creates a new container instance
func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name:      os.Getenv("APP_NAME"),
		Env:       os.Getenv("APP_ENV"),
		Port:      os.Getenv("APP_PORT"),
		IsCacheOn: os.Getenv("IS_CACHE_ON"),
	}

	redis := &Redis{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	db := &DB{
		Connection:  os.Getenv("DB_CONNECTION"),
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		User:        os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		ClusterName: os.Getenv("DB_CLUSTER_NAME"),
		DBName:      os.Getenv("DB_NAME"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	isValid, errMsg := app.validate()
	if !isValid {
		panic(errMsg)
	}
	isValid, errMsg = redis.validate()
	if !isValid {
		panic(errMsg)
	}
	isValid, errMsg = db.validate()
	if !isValid {
		panic(errMsg)
	}

	return &Container{
		app,
		redis,
		db,
		http,
	}, nil
}

func (a App) validate() (isValid bool, errMessage string) {
	isValid = true
	errMessage = "invalid"
	switch {
	case strings.EqualFold(strings.TrimSpace(a.Name), ""):
		isValid = false
		errMessage = "Please provide APP_NAME"
	case strings.EqualFold(strings.TrimSpace(a.Env), ""):
		isValid = false
		errMessage = "Please provide APP_ENV"
	case strings.EqualFold(strings.TrimSpace(a.Port), ""):
		isValid = false
		errMessage = "Please provide APP_PORT"
	}

	return
}

func (db DB) validate() (isValid bool, errMessage string) {
	isValid = true
	errMessage = "invalid"
	switch {
	case strings.EqualFold(strings.TrimSpace(db.Connection), ""):
		isValid = false
		errMessage = "Please provide DB_CONNECTION"

	case strings.EqualFold(strings.TrimSpace(db.Host), ""):
		isValid = false
		errMessage = "Please provide DB_HOST"

	case strings.EqualFold(strings.TrimSpace(db.User), ""):
		isValid = false
		errMessage = "Please provide DB_USER"

	case strings.EqualFold(strings.TrimSpace(db.Password), ""):
		isValid = false
		errMessage = "Please provide DB_PASSWORD"

	case strings.EqualFold(strings.TrimSpace(db.ClusterName), ""):
		isValid = false
		errMessage = "Please provide DB_CLUSTER_NAME"

	case strings.EqualFold(strings.TrimSpace(db.DBName), ""):
		isValid = false
		errMessage = "Please provide DB_NAME"
	}

	return
}

func (rd Redis) validate() (isValid bool, errMessage string) {
	isValid = true
	errMessage = "invalid"
	switch {
	case strings.EqualFold(strings.TrimSpace(rd.Addr), ""):
		isValid = false
		errMessage = "Please provide REDIS_ADDR"
	}

	return
}
