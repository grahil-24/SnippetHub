package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"snippetbox.rahilganatra.net/internal/models"
)

// dependency injection so that handlers can use our custom loggers too.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {

	//addr variable can be used as an flag at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	//database command lind variable
	dsn := flag.String("dsn", "root:grahil11@/snippetbox?parseTime=true", "mysql DSN")

	flag.Parse()

	db, err := openDB(*dsn)
	fmt.Print("db opened\n")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	//creates new logger with 3 params
	//- first takes the destination
	//- second is the prefix
	//- third is extra info like date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}
	fmt.Println("app created\n")
	defer db.Close()

	mux := app.routes()

	infoLog.Print("Starting server at ", *addr)

	//make the server use our custom error logger
	srv := &http.Server{
		ErrorLog: errorLog,
		Handler:  mux,
		Addr:     *addr,
	}

	err_srv := srv.ListenAndServe()
	infoLog.Print("Server started at ", *addr)
	log.Fatal(err_srv)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
