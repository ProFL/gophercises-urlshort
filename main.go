package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ProFL/gophercises-urlshort/handlers"
	"github.com/ProFL/gophercises-urlshort/repositories"
	_ "github.com/mattn/go-sqlite3"
)

var yamlFilePath string
var jsonFilePath string
var dataSourceName string

func init() {
	flag.StringVar(&yamlFilePath, "yaml", "", "yaml file redirect mappings")
	flag.StringVar(&jsonFilePath, "json", "", "json file redirect mappings")
	flag.StringVar(&dataSourceName, "db", "./sample.db?cache=shared&mode=rw", "database connection string")
	flag.Parse()
}

func connectToDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Panic("Failed to connect to the database", err.Error())
	}
	return db
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := make([]byte, 0)
	if yamlFilePath != "" {
		var err error
		log.Println("Reading YAML file...")
		yaml, err = os.ReadFile(yamlFilePath)
		if err != nil {
			log.Panic("Failed to read YAML file", err.Error())
		}
	}

	yamlHandler, err := handlers.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json := []byte("[]")
	if jsonFilePath != "" {
		var err error
		log.Println("Reading JSON file...")
		json, err = os.ReadFile(jsonFilePath)
		if err != nil {
			log.Panic("Failed to read JSON file", err.Error())
		}
	}

	jsonHandler, err := handlers.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	db := connectToDatabase()
	defer db.Close()

	redirectRepo := repositories.RedirectRepository{Db: db}
	redirectRepo.CreateTableIfNotExists()
	redirectRepo.Seed()

	databaseHandler := handlers.DatabaseHandler(&redirectRepo, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", databaseHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
