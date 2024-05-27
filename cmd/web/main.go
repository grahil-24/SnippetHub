package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	//addr variable can be used as an flag at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	//tells go where static files are present
	fileServer := http.FileServer(http.Dir("./ui/static"))

	//creates new logger with 3 params
	//- first takes the destination
	//- second is the prefix
	//- third is extra info like date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//strips the uri to extract the file name.
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	//creating a new server
	mux := http.NewServeMux()
	//route handler
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	infoLog.Print("Starting server at ", *addr)

	//make the server use our custom error logger
	srv := &http.Server{
		ErrorLog: errorLog,
		Handler:  mux,
		Addr:     *addr,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)

}
