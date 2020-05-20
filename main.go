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
	"github.com/gorilla/handlers"

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

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	credentialsOk := handlers.AllowCredentials()
	// accesscontroll = handlers.Access
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	CorsMiddleware := handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)

	router := mux.NewRouter().StrictSlash(true)
	router.Use(loggingMiddleware)

	// fs := http.FileServer(http.Dir("./static/static"))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/signup", auth.Signup).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	userRouter := router.PathPrefix("/api").Subrouter()

	userRouter.Use(jwthelper.Jwthelper)
	userRouter.Use(CorsMiddleware)

	userRouter.HandleFunc("/add", api.Addroute).Methods("POST")
	userRouter.HandleFunc("/offlinesyncadd", api.OfflinesyncAddroute).Methods("POST")
	userRouter.HandleFunc("/getall", api.Getall).Methods("GET")
	userRouter.HandleFunc("/gettrashall", api.GetTrashall).Methods("GET")
	userRouter.HandleFunc("/update", api.Modify).Methods("POST")
	userRouter.HandleFunc("/changestate", api.Changenotestate).Methods("POST")

	// contextedMux := AddContext(router)

	// Serve static files

	buildHandler := http.FileServer(http.Dir("./build"))
	router.PathPrefix(`/{rest:[a-zA-Z0-9=\-\/]+}`).Handler(buildHandler)

	// staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("/build/static")))
	// router.PathPrefix("/static/").Handler(staticHandler)
	// router.PathPrefix("/").Handler(http.FileServer(http.Dir("")))

	// Serve index page on all unhandled routes
	// router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./build/index.html")
	// })

	// router.Handle("*", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(GetPort(), CorsMiddleware(router)))

}
