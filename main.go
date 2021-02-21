package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Db *sql.DB

type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Age       int    `json:"age"`
}

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:Root@123@tcp(localhost:3306)/sakila")
	if err != nil {
		panic(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAll(w http.ResponseWriter, r *http.Request) {
	results, err := Db.Query("SELECT * FROM persons")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	psList := []Person{}
	for results.Next() {
		ps := Person{}
		err = results.Scan(&ps.ID, &ps.FirstName, &ps.LastName, &ps.Age)
		if err != nil {
			panic(err.Error())
		}
		psList = append(psList, ps)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonValue, _ := json.Marshal(psList)
	w.Write(jsonValue)
}

func returnSingle(w http.ResponseWriter, r *http.Request) {
	ps := Person{}
	vars := mux.Vars(r)
	key := vars["id"]
	err := Db.QueryRow("SELECT * from persons where id = ?", key).Scan(&ps.ID, &ps.FirstName, &ps.LastName, &ps.Age)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonValue, _ := json.Marshal(ps)
	w.Write(jsonValue)
}

func returnDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	_, err := Db.Exec("Delete from persons where id = ?", key)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	w.Write([]byte("Deletion is successful"))
}

func returnInsert(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	ps := Person{}
	json.Unmarshal(reqBody, &ps)
	_, err := Db.Exec("INSERT into persons(ID, FirstName, LastName, Age) VALUES(?, ?, ?, ?)", ps.ID, ps.FirstName, ps.LastName, ps.Age)

	if err != nil {
		panic(err.Error())
	}
	w.Write([]byte("instertion is succesful"))
}

func returnUpdate(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	ps := Person{}
	json.Unmarshal(reqBody, &ps)
	_, err := Db.Exec("update persons set FirstName = ?, LastName = ?, Age= ? where ID = ?", ps.FirstName, ps.LastName, ps.Age, ps.ID)

	if err != nil {
		panic(err.Error())
	}
	w.Write([]byte("updation is succesful"))
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/person/{id}", returnSingle)
	myRouter.HandleFunc("/persons", returnAll)
	myRouter.HandleFunc("/delete/{id}", returnDelete)
	myRouter.HandleFunc("/person-entry", returnInsert)
	myRouter.HandleFunc("/person-update", returnUpdate)
	log.Fatal(http.ListenAndServe(":8080", myRouter))

	defer Db.Close()
}
