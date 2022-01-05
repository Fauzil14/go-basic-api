package main

import (
	"fmt"
	"log"
	"net/http"
	"quiz-3/functions"

	"github.com/julienschmidt/httprouter"
)

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && (user == "admin" && password == "password") || (user == "editor" && password == "secret") || (user == "trainer" && password == "rahasia") {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func main() {
	router := httprouter.New()

	// soal 2
	router.GET("/bangun-datar/segitiga-sama-sisi", functions.GetFunction)
	router.GET("/bangun-datar/persegi", functions.GetFunction)
	router.GET("/bangun-datar/persegi-panjang", functions.GetFunction)
	router.GET("/bangun-datar/lingkaran", functions.GetFunction)
	router.GET("/bangun-datar/jajar-genjang", functions.GetFunction)

	// soal 3
	router.GET("/categories", functions.GetCategory)
	router.POST("/categories", BasicAuth(functions.PostCategory))
	router.PUT("/categories/:id", BasicAuth(functions.UpdateCategory))
	router.DELETE("/categories/:id", BasicAuth(functions.DeleteCategory))
	// router.GET("/categories/:id/books", functions.GetBooksFromCategory)
	
	router.GET("/books", functions.GetBook)
	router.POST("/books", BasicAuth(functions.PostBook))
	router.PUT("/books/:id", BasicAuth(functions.UpdateBook))
	router.DELETE("/books/:id", BasicAuth(functions.DeleteBook))
	
	fmt.Println("Server is runnin at port 8003")
	log.Fatal(http.ListenAndServe(":8003", router))
}
