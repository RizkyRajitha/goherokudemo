package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/RizkyRajitha/goherokudemo/dbutil"
	uuid "github.com/satori/go.uuid"



	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Claims struct {
	UserId string `json:"userid"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(creds)

	var user dbutil.User

	if dbutil.DBcon.Where("email = ?", creds.Email).First(&user).RecordNotFound() {
		// record not found
		fmt.Println("no user found")
		type nosuerErr struct {
			Msg string `json:"msg"`
		}

		println("duplicate email")
		w.WriteHeader(http.StatusNotFound)
		var payload nosuerErr
		payload.Msg = "Invalid_email"
		json.NewEncoder(w).Encode(payload)

	} else {
		fmt.Println(user)

		var hasherr error

		hasherr = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(creds.Password))
		// fmt.Println("comp - ")

		if hasherr != nil {
			fmt.Println(hasherr)

			type hashErr struct {
				Msg string `json:"msg"`
			}

			println("duplicate email")
			w.WriteHeader(http.StatusUnauthorized)
			var payload hashErr
			payload.Msg = "Invalid_password"
			json.NewEncoder(w).Encode(payload)
		} else {
			fmt.Println("all good sir eh !!!111")

			expirationTime := time.Now().Add(600 * time.Minute)
			// Create the JWT claims, which includes the username and expiry time
			claims := &Claims{
				UserId: user.UserId,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}
			// Declare the token with the algorithm used for signing, and the claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// Create the JWT string
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				// If there is an error in creating the JWT return an internal server error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Println(tokenString)

			type tokenPayload struct {
				Token string `json:"token"`
			}

			var payload tokenPayload
			payload.Token = tokenString
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)

		}

	}

}

type Userreg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Creating New user")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user dbutil.User
	var userdummy Userreg
	userid := uuid.NewV4().String() //uuid.Must(uuid.NewV4()).String()
	user.UserId = userid        //uuid.Must(uuid.NewV4()).String()
	user.Created = time.Now().Format(time.RFC3339)

	json.Unmarshal(reqBody, &userdummy)
	password := []byte(userdummy.Password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	} else {
		user.Hash = string(hashedPassword)
		user.Name = userdummy.Name
		json.Unmarshal(reqBody, &user)

		// var dberr error

		fmt.Println(user)

		if dberr := dbutil.DBcon.Create(&user).Error; dberr != nil {
			// error handling...
			fmt.Println(dberr)
			type errdb struct {
				Msg string `json:"msg"`
			}

			println("duplicate email")
			w.WriteHeader(http.StatusForbidden)
			var payload errdb
			payload.Msg = "Duplicate_email"
			json.NewEncoder(w).Encode(payload)
		} else {
			println("success")
			type success struct {
				Msg string `json:"msg"`
			}

			w.WriteHeader(http.StatusOK)
			var payload success
			payload.Msg = "success"
			json.NewEncoder(w).Encode(payload)
		}

	}

}
