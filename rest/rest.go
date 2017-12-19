package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"

)

/*
RouterObject provides a pointer to the underlying mux object for status checks
 */
type RouterObject struct {
	Mux *mux.Router
}

var router *RouterObject

func init() {
	router = &RouterObject{}
	go func() {
		log.Println("Starting HTTP Server...")
		log.Fatal(http.ListenAndServe(":8001", getRouter()))
	}()
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)



		}
	}()
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/inwards/cutoff", initializeInwards)

	router.Mux = r
	return r
}

/*
Router starts the router service
 */
func Router() (*RouterObject) {
	return router
}
