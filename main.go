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
	"github.com/RizkyRajitha/goherokudemo/websocket"
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

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id := params["id"]

	println(id)

	ws, err := websocket.Upgrade(w, r)

	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	println(r.RequestURI)

	client := &websocket.Client{
		ID:   id,
		Conn: ws,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()

	// websocket.Writer(ws)
	// websocket.Reader(ws)

}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,

// 	// We'll need to check the origin of our connection
// 	// this will allow us to make requests from our React
// 	// development server to here.
// 	// For now, we'll do no checking and just allow any connection
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// // define a reader which will listen for
// // new messages being sent to our WebSocket
// // endpoint
// func reader(conn *websocket.Conn) {
// 	for {
// 		// read in a message
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		// print out that message for clarity
// 		fmt.Println(string(p))

// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			log.Println(err)
// 			return
// 		}

// 	}
// }

// // define our WebSocket endpoint
// func serveWs(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r.Host)

// 	// upgrade this connection to a WebSocket
// 	// connection
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	// listen indefinitely for new messages coming
// 	// through on our WebSocket connection
// 	reader(ws)
// }

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

	pool := websocket.NewPool()
	go pool.Start()

	router.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	userRouter := router.PathPrefix("/api").Subrouter()

	userRouter.Use(jwthelper.Jwthelper)
	userRouter.Use(CorsMiddleware)

	userRouter.HandleFunc("/add", api.Addroute).Methods("POST")
	userRouter.HandleFunc("/offlinesyncadd", api.OfflinesyncAddroute).Methods("POST")
	userRouter.HandleFunc("/getall", api.Getall).Methods("GET")
	userRouter.HandleFunc("/gettrashall", api.GetTrashall).Methods("GET")
	userRouter.HandleFunc("/update", api.Modify).Methods("POST")
	userRouter.HandleFunc("/changestate", api.Changenotestate).Methods("POST")
	userRouter.HandleFunc("/changepinnedstate", api.Changenotepinnned).Methods("POST")

	// contextedMux := AddContext(router)

	//  wshander :=  http.Handle("/ws", serveWs)
	// Serve static files

	buildHandler := http.FileServer(http.Dir("./build"))
	router.PathPrefix("/").Handler(buildHandler)

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
