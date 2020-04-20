package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RizkyRajitha/goherokudemo/api"
	"github.com/RizkyRajitha/goherokudemo/auth"
	"github.com/RizkyRajitha/goherokudemo/dbutil"
	jwthelper "github.com/RizkyRajitha/goherokudemo/middleware"

	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	fmt.Println("hello")
	dbutil.ConnectDB()

	log.Println("Starting development server at http://127.0.0.1" + GetPort())

	router := mux.NewRouter().StrictSlash(true)
	router.Use(loggingMiddleware)
	router.HandleFunc("/", api.Homeroute).Methods("GET")
	router.HandleFunc("/signup", auth.Signup).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	userRouter := router.PathPrefix("/api").Subrouter()

	userRouter.Use(jwthelper.Jwthelper)

	userRouter.HandleFunc("/add", api.Addroute).Methods("POST")
	userRouter.HandleFunc("/offlinesyncadd", api.OfflinesyncAddroute).Methods("POST")
	userRouter.HandleFunc("/getall", api.Getall).Methods("GET")
	userRouter.HandleFunc("/update", api.Modify).Methods("POST")

	// contextedMux := AddContext(router)

	log.Fatal(http.ListenAndServe(GetPort(), router))

}
