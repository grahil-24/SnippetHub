package main

import "net/http"

//all custom middlewares  go here

// want this middleware to be executed on every request
func secureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//restricts where the resources like javascript, fonts, images, css can be loaded from
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		//sets if the page can be embedded in other sites. prevents clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")
		//to prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		//0 is to turn off XSS. Use the better alternative ie the Content-Security-Policy
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
