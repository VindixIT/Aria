package main

import (	
	"log"
	"net/http"
	"database/sql"
	"strconv"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type Unit struct {
    Id    int
    Symbol  string
    Description  string
    Selected bool
}
func InitUnitsTable(db *sql.DB) {
    log.Println("Init Units")
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Units ( " +
		" id SERIAL PRIMARY KEY, "+
		" symbol varchar(5) NOT NULL, " +
		" description varchar(20) NOT NULL " +
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListUnits(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	log.Println("List Units")
	selDB, err := db.Query("SELECT id, symbol, description FROM Units ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	unit := Unit{}
    units := []Unit{}
	for selDB.Next() {
	    var id int
        var symbol, description string
        err = selDB.Scan(&id, &symbol, &description)
        if err != nil {
            panic(err.Error())
        }
        unit.Id = id
        unit.Symbol = symbol      
        unit.Description = description      
        units = append(units, unit)
    }
	tmpl.ExecuteTemplate(w, "ListUnits", units) 
	defer db.Close()
}

func ShowUnit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Show Unit")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, symbol, description FROM Units WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    unit := Unit{}
    for selDB.Next() {
        var id int
        var description, symbol string
        err = selDB.Scan(&id, &symbol, &description)
        if err != nil {
            panic(err.Error())
        }
        unit.Id = id
        unit.Description = description
        unit.Symbol = symbol
    }
    tmpl.ExecuteTemplate(w, "ShowUnit", unit)
    defer db.Close()
}

func EditUnit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Edit Unit")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, symbol, description FROM Units WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    unit := Unit{}
    for selDB.Next() {
        var id int
        var symbol, description string
        
        err = selDB.Scan(&id, &symbol, &description)
        if err != nil {
            panic(err.Error())
        }
        unit.Id = id
        unit.Symbol = symbol
        unit.Description = description
    }
    tmpl.ExecuteTemplate(w, "EditUnit", unit)
    defer db.Close()
}

func InsertUnit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Insert Unit")
    if r.Method == "POST" {
        symbol := r.FormValue("symbol")
        description := r.FormValue("description")
		sqlStatement := "INSERT INTO Units(symbol, description) VALUES ($1, $2) RETURNING id"
		id := 0
		err := db.QueryRow(sqlStatement, symbol, description).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | Symbol: " + symbol+" | Description: " + description)
    }
    defer db.Close()
    http.Redirect(w, r, "/listUnits", 301)
}

func UpdateUnit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Update Unit")
    if r.Method == "POST" {
        symbol := r.FormValue("symbol") 
        description := r.FormValue("description")
        id := r.FormValue("uid")        
		sqlStatement := "UPDATE Units SET symbol=$1, description=$2 WHERE id=$3"
		updtForm, err := db.Prepare(sqlStatement)
        if err != nil {
            panic(err.Error())
		}    
		updtForm.Exec(symbol, description, id)
        log.Println("UPDATE: Id: " + id +" | Symbol: " +symbol+" | Description: " + description ) 
    }
    defer db.Close()
    http.Redirect(w, r, "/listUnits", 301)
}

func DeleteUnit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Delete Unit")
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Units WHERE id=$1")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE: Id: " + id)
    defer db.Close()
    http.Redirect(w, r, "/listUnits", 301)
}

func NewUnit(w http.ResponseWriter, r *http.Request) {
    log.Println("New Unit")
	tmpl.ExecuteTemplate(w, "NewUnit", nil)
}