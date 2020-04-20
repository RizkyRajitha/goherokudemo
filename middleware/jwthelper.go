package jwthelper

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

type Claims struct {
	UserId string `json:"userid"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

func Jwthelper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		tk := r.Header.Get("Authorization")

		if tk == "" {

			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return

		}

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tk, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		println("middleware *************************")
		println(claims.UserId)
		// println(err)

		context.Set(r, "Userid", claims.UserId)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		println(tk)

		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
