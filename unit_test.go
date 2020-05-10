package main

import (
	"bytes"
	// "database/sql"
	// "fmt"
	// "os"

	// "context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// "github.com/RizkyRajitha/goherokudemo/api"
	"github.com/RizkyRajitha/goherokudemo/api"
	"github.com/RizkyRajitha/goherokudemo/auth"
	"github.com/RizkyRajitha/goherokudemo/dbutil"
	// "github.com/dgrijalva/jwt-go"
)

var Token = ""

func TestSignup(t *testing.T) {

	var jsonStr = []byte(`{"email":"test123" , "name" : "123" , "password":"123" }`)


	dbutil.ConnectDB()

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.Signup)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"msg":"success"}` + "\n"
	print((expected))
	print(rr.Body.String())

	println(rr.Body.String() == (expected))
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// initDb()

}

// var (
// 	dbhost = os.Getenv("HOSTIP")
// 	dbport = "5432"
// 	dbuser = os.Getenv("DBUSER")
// 	dbpass = os.Getenv("DBPASSWORD")
// 	dbname = os.Getenv("DBUSER")
// )

// var db *sql.DB

// func initDb() {

// 	var err error
// 	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		dbhost, dbport,
// 		dbuser, dbpass, dbname)

// 	println("host=%v port=%v user=%v "+
// 		"password=%v dbname=%v sslmode=disable",
// 		dbhost, dbport,
// 		dbuser, dbpass, dbname)

// 	db, err = sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Successfully connected!")
// }

func TestSignupDuplicate(t *testing.T) {

	var jsonStr = []byte(`{"email":"test123" , "name" : "123" , "password":"123" }`)

	// req, err := http.NewRequest("PUT", "/entry", bytes.NewBuffer(jsonStr))

	dbutil.ConnectDB()

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.Signup)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	// Check the response body is what we expect.
	expected := `{"msg":"Duplicate_email"}` + "\n"
	print((expected))
	print(rr.Body.String())

	println(rr.Body.String() == (expected))
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLogin(t *testing.T) {

	var jsonStr = []byte(`{"email":"test123" , "password":"123" }`)

	// req, err := http.NewRequest("PUT", "/entry", bytes.NewBuffer(jsonStr))

	dbutil.ConnectDB()

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// ctx := context.Background()
	// ctx = context.WithValue(ctx, "app.auth.token", "abc123")

	handler := http.HandlerFunc(auth.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"msg":"Duplicate_email"}` + "\n"
	// print((expected))
	// print(rr.Body.String())

	type tokens struct {
		Token string `json:"token"`
	}

	var tk tokens

	json.Unmarshal(rr.Body.Bytes(), &tk)

	print(tk.Token)

	Token = tk.Token

	// println(rr.Body.String() == (expected))
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

func TestLoginInvaliduser(t *testing.T) {

	var jsonStr = []byte(`{"email":"test1234" , "password":"123" }`)

	// req, err := http.NewRequest("PUT", "/entry", bytes.NewBuffer(jsonStr))

	dbutil.ConnectDB()

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Check the response body is what we expect.
	expected := `{"msg":"Invalid_email"}` + "\n"
	print((expected))
	print(rr.Body.String())

	println(rr.Body.String() == (expected))
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLoginInvalidPassword(t *testing.T) {

	var jsonStr = []byte(`{"email":"test123" , "password":"1234" }`)

	// req, err := http.NewRequest("PUT", "/entry", bytes.NewBuffer(jsonStr))

	dbutil.ConnectDB()

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is what we expect.
	expected := `{"msg":"Invalid_password"}` + "\n"
	print((expected))
	print(rr.Body.String())

	println(rr.Body.String() == (expected))
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// func NewContextWithRequestID(ctx context.Context, r *http.Request) context.Context {
//     return context.WithValue(ctx, "reqId", "1234")
// }

// func AddContextWithRequestID(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         var ctx = context.Background()
//         ctx = NewContextWithRequestID(ctx, r)
//         next.ServeHTTP(w, r.WithContext(ctx))
//     })
// }

// func TestGetallNotes(t *testing.T) {

// 	// var jsonStr = []byte(`{"email":"test123" , "password":"1234" }`)

// 	// req, err := http.NewRequest("PUT", "/entry", bytes.NewBuffer(jsonStr))

// 	dbutil.ConnectDB()

// 	req, err := http.NewRequest("GET", "/api/getall", nil)
// 	// println(Token)
// 	// req.Header.Set("name", "value")
// 	req.Header.Set("Authorization", Token)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()

// 	handler := http.HandlerFunc( (api.Getall))

// 	handler.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	expected := `{"msg":"Invalid_password"}` + "\n"
// 	print((expected))
// 	print(rr.Body.String())

// 	println(rr.Body.String() == (expected))
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
