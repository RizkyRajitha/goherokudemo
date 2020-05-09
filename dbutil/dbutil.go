package dbutil

import (
	// "fmt"
	"log"
	"os"

	// "github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Notes struct {
	Id       string `json:"id" gorm:"primary_key"`
	UserId   string `json:"userid"`
	Note     string `json:"note"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Active   bool   `json:"active"`
}

type User struct {
	Email   string `json:"email" gorm:"primary_key"`
	UserId  string `json:"userid"`
	Hash    string `json:"hash"`
	Created string `json:"created"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

var DBcon *gorm.DB
var err error

func ConnectDB() {

	var hostip = os.Getenv("HOSTIP")
	var dbuser = os.Getenv("DBUSER")
	var dbpassword = os.Getenv("DBPASSWORD")

	var connstrprod = "host= " + hostip + " port=5432 user=" + dbuser + " dbname=" + dbuser + " password=" + dbpassword

	// var connstrdev = "host= " + "localhost" + " port=5432 user=" + "superroot" + " dbname=" + "gotest" + " password=" + "123"

	DBcon, err = gorm.Open("postgres", connstrprod)
	// println(connstrdev)
	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function

	if err != nil {
		println(err)
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
		DBcon.AutoMigrate(&Notes{})
		DBcon.AutoMigrate(&User{})

		// handleRequests()
	}

}
