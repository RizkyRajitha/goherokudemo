package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/RizkyRajitha/goherokudemo/dbutil"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/context"
)

func Homeroute(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")

}

func Addroute(w http.ResponseWriter, r *http.Request) {

	// get the body of our POST request
	// return the string response containing the request body

	reqBody, _ := ioutil.ReadAll(r.Body)

	userid12121 := context.Get(r, "Userid")
	println("add new note")
	println("*************")
	println(userid12121.(string))
	uerid := userid12121.(string)

	var booking dbutil.Notes
	noteid := uuid.NewV4().String() //uuid.Must(uuid.NewV4()).String()
	booking.Id = noteid             //uuid.UUID()
	booking.UserId = uerid
	booking.Updated = time.Now().Format(time.RFC3339)
	booking.Created = time.Now().Format(time.RFC3339)
	booking.Active = true
	json.Unmarshal(reqBody, &booking)

	fmt.Println(booking)
	dbutil.DBcon.Create(&booking)
	fmt.Println("Endpoint Hit: Creating New Booking")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)

}

func OfflinesyncAddroute(w http.ResponseWriter, r *http.Request) {

	// get the body of our POST request
	// return the string response containing the request body

	reqBody, _ := ioutil.ReadAll(r.Body)

	userid12121 := context.Get(r, "Userid")
	println("add new note")
	println("*************")
	println(userid12121.(string))
	uerid := userid12121.(string)

	var booking dbutil.Notes
	noteid := uuid.NewV4().String()
	booking.Id = noteid
	booking.UserId = uerid
	// booking.Updated = time.Now().Format(time.RFC3339)
	// booking.Created = time.Now().Format(time.RFC3339)
	json.Unmarshal(reqBody, &booking)

	fmt.Println(booking)
	dbutil.DBcon.Create(&booking)
	fmt.Println("Endpoint Hit: Offline sync New Booking")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)

}

func Modify(w http.ResponseWriter, r *http.Request) {

	type modistrut struct {
		Note   string `json:"note"`
		NoteID string `json:"noteid"`
		Title  string `json:"title"`
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var modnotes modistrut

	var notes dbutil.Notes

	uptime := time.Now().Format(time.RFC3339)
	json.Unmarshal(reqBody, &modnotes)

	fmt.Println(modnotes)
	println(uptime)

	userid12121 := context.Get(r, "Userid")
	println(userid12121.(string))
	uerid := userid12121.(string)

	var upnote dbutil.Notes

	dbutil.DBcon.Where("id = ?", modnotes.NoteID).Find(&upnote)

	println(uerid)
	println("ll")
	println(upnote.UserId)

	if uerid != upnote.UserId {
		type errdb struct {
			Msg string `json:"msg"`
		}

		println("duplicate email")
		w.WriteHeader(http.StatusForbidden)
		var payload errdb
		payload.Msg = "unauthorized"
		json.NewEncoder(w).Encode(payload)

	} else {
		dbutil.DBcon.Model(&notes).Where("id = ?", modnotes.NoteID).Update("note", modnotes.Note)
		dbutil.DBcon.Model(&notes).Where("id = ?", modnotes.NoteID).Update("updated", uptime)
		dbutil.DBcon.Model(&notes).Where("id = ?", modnotes.NoteID).Update("title", modnotes.Title)

		type payload struct {
			Msg string `json:"msg"`
		}

		var payloadsend payload
		payloadsend.Msg = "success"
		fmt.Println("Endpoint Hit: update notes")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payloadsend)
	}

	// get the body of our POST request
	// return the string response containing the request body
	// reqBody, _ := ioutil.ReadAll(r.Body)
	// var booking dbutil.Notes
	// booking.Id = uuid.UUID()
	// booking.Created = time.Now().Format(time.RFC3339)
	// json.Unmarshal(reqBody, &booking)

	// fmt.Println(booking)
	// dbutil.DBcon.Create(&booking)

}

func Getall(w http.ResponseWriter, r *http.Request) {

	// get the body of our POST request
	// return the string response containing the request body

	userid12121 := context.Get(r, "Userid")
	println(userid12121.(string))
	uerid := userid12121.(string)
	// db.Where("name = ?", "jinzhu").First(&user)
	notes := []dbutil.Notes{}
	// dbutil.DBcon.Find(&notes).Where("userid = ?", uerid)
	dbutil.DBcon.Where("user_id = ? AND active = ?", uerid, true).Order("updated desc").Find(&notes)

	var user dbutil.User

	dbutil.DBcon.Where("user_id = ? ", uerid).Find(&user)

	type payload struct {
		Username string         `json:"username"`
		Notes    []dbutil.Notes `json:"notes"`
	}

	var payloadsend payload

	payloadsend.Notes = notes
	payloadsend.Username = user.Name

	fmt.Println("Endpoint Hit: Get all notes")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payloadsend)

}

func Changenotestate(w http.ResponseWriter, r *http.Request) {

	type modistrut struct {
		NoteID string `json:"noteid"`
		State  bool   `json:"state"`
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var modnotes modistrut

	var notes dbutil.Notes

	// uptime := time.Now().Format(time.RFC3339)
	json.Unmarshal(reqBody, &modnotes)

	fmt.Println(modnotes)
	// println(uptime)

	userid12121 := context.Get(r, "Userid")
	// println(userid12121.(string))
	uerid := userid12121.(string)

	var upnote dbutil.Notes

	dbutil.DBcon.Where("id = ?", modnotes.NoteID).Find(&upnote)

	println(uerid)
	println("ll")
	println(upnote.UserId)

	if uerid != upnote.UserId {
		type errdb struct {
			Msg string `json:"msg"`
		}

		println("invalid user")
		w.WriteHeader(http.StatusForbidden)
		var payload errdb
		payload.Msg = "unauthorized"
		json.NewEncoder(w).Encode(payload)

	} else {
		dbutil.DBcon.Model(&notes).Where("id = ?", modnotes.NoteID).Update("active", modnotes.State)

		type payload struct {
			Msg string `json:"msg"`
		}

		var payloadsend payload
		payloadsend.Msg = "success"
		fmt.Println("Endpoint Hit: update notes")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payloadsend)
	}

}


func GetTrashall(w http.ResponseWriter, r *http.Request) {

	// get the body of our POST request
	// return the string response containing the request body

	userid12121 := context.Get(r, "Userid")
	println(userid12121.(string))
	uerid := userid12121.(string)
	// db.Where("name = ?", "jinzhu").First(&user)
	notes := []dbutil.Notes{}
	// dbutil.DBcon.Find(&notes).Where("userid = ?", uerid)
	dbutil.DBcon.Where("user_id = ? AND active = ?", uerid, false).Order("updated desc").Find(&notes)

	var user dbutil.User

	dbutil.DBcon.Where("user_id = ? ", uerid).Find(&user)

	type payload struct {
		Username string         `json:"username"`
		Notes    []dbutil.Notes `json:"notes"`
	}

	var payloadsend payload

	payloadsend.Notes = notes
	payloadsend.Username = user.Name

	fmt.Println("Endpoint Hit: Get all notes")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payloadsend)

}