package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"frontendmasters.com/reelingit/data"
	"frontendmasters.com/reelingit/handlers"
	"frontendmasters.com/reelingit/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func initializeLogger() *logger.Logger {
	ll, err := logger.NewLogger("movie.log")
	// ll.Error("hello from logging system", nil)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	defer ll.Close()
	return ll
}

func main() {

	// log init
	logInstance := initializeLogger()

	// environmental variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file was available")
	}

	//connect to the database
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	//init movie repo
	movieRepo, err := data.NewMovieRepository(db, logInstance)
	if err != nil {
		log.Fatalf("failed to initialize the repository: %v", err)
	}
	// handler init
	movieHandler := handlers.NewMovieHandler(movieRepo, logInstance)

	http.HandleFunc("/api/movies/top/", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/random/", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/movies/search/", movieHandler.SearchMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovie) // api/movie/140
	http.HandleFunc("/api/genres/", movieHandler.GetGenres)

	catchAllClientRoutesHandler := func(w http.ResponseWriter, r *http.Request) {
		// render index.html
		fmt.Println("serving from the new handlefunc")
		http.ServeFile(w, r, "./public/index.html")
	}

	http.HandleFunc("/movies", catchAllClientRoutesHandler)
	http.HandleFunc("/movies/", catchAllClientRoutesHandler)
	http.HandleFunc("/account/", catchAllClientRoutesHandler)

	// handler for static files
	http.Handle("/", http.FileServer(http.Dir("public")))

	fmt.Println("Serving the files")

	const addr = ":8080"
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server has failed: %v", err)
		logInstance.Error("server failed", err)
	}
}
