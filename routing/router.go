package routing

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {

	STATIC_DIR := "/resources/"

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix(STATIC_DIR).
		Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	//following this -> router.Methods("GET").Path("/hello").Name("hello").HandlerFunc(Hello)
	//router.Methods("GET").Path("/hello").Name("hello").HandlerFunc(Hello)

	for _, route := range routes {

		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(Logger(route.HandlerFunc, route.Name))

	}

	return router
}
