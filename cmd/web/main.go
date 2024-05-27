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
	//-
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//strips the uri to extract the file name.
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	//creating a new server
	mux := http.NewServeMux()

	//route handler
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)

}
