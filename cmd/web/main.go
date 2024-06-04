package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox.rahilganatra.net/internal/models"
	"time"
)

// dependency injection so that handlers can use our custom loggers too.
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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
	templateCache, err := newTemplatesCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//initialize a form decode
	formDecoder := form.NewDecoder()

	//initialize a new session
	sessionManager := scs.New()
	//set the store for the sessions as mysql database
	sessionManager.Store = mysqlstore.New(db)
	//session should expire after 12 hours of creation
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	fmt.Println("app created\n")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	infoLog.Print("Starting server at ", *addr)

	//Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		//these two are not as cpu intensive as other methods, so these will be used for better
		//performance
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	//make the server use our custom error logger
	srv := &http.Server{
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		Addr:      *addr,
		TLSConfig: tlsConfig,
		//add idle, read and write timeouts
		IdleTimeout: time.Minute,
		// This means that if the request headers or body are still being read 5 seconds after the request is first accepted, then
		//Go will close the underlying connection
		ReadTimeout: 5 * time.Second,
		//if data is written after writetimeout value from the time request was accepted, GO will close
		//the connection. WriteTimeout should always be greater than ReadTimeout
		WriteTimeout: 10 * time.Second,
	}

	//ListenAndServe starts an HTTP while TLS starts an HTTPS. Before creating an HTTPS server
	//we need to generate TLS certificates. We can do that by creating a TLS folder in our root project dir
	//and running the command go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
	//in the terminal. It will generate an public and private key stored in cert.pem and key.pem files respectively
	errSrv := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	infoLog.Print("Server started at ", *addr)
	log.Fatal(errSrv)
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
