package main

import (	
	"log"
	"net/http"
	"database/sql"
	"strconv"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type Insulin struct {
    Id    int
    Name  string
    Selected bool
}
func InitInsulinsTable(db *sql.DB) {
    log.Println("Init Insulins")
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Insulins ( " +
		" id SERIAL PRIMARY KEY, "+
		" name varchar(20) NOT NULL " +
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListInsulins(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	log.Println("List Insulins")
	selDB, err := db.Query("SELECT * FROM Insulins ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	insulin := Insulin{}
    insulins := []Insulin{}
	for selDB.Next() {
		var id int
      var name string
      err = selDB.Scan(&id, &name)
      if err != nil {
         panic(err.Error())
      }
      insulin.Id = id
      insulin.Name = name      
      insulins = append(insulins, insulin)
}
	tmpl.ExecuteTemplate(w, "ListInsulins", insulins) 
	defer db.Close()
}

func ShowInsulin(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Show Insulin")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, name FROM Insulins WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    insulin := Insulin{}
    for selDB.Next() {
        var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        insulin.Id = id
        insulin.Name = name
    }
    tmpl.ExecuteTemplate(w, "ShowInsulin", insulin)
    defer db.Close()
}

func EditInsulin(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Edit Insulin")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, name FROM Insulins WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    insulin := Insulin{}
    for selDB.Next() {
        var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        insulin.Id = id
        insulin.Name = name
    }
    tmpl.ExecuteTemplate(w, "EditInsulin", insulin)
    defer db.Close()
}

func InsertInsulin(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Insert Insulin")
    if r.Method == "POST" {
        name := r.FormValue("name")
		sqlStatement := "INSERT INTO Insulins(name) VALUES ($1) RETURNING id"
		id := 0
		err := db.QueryRow(sqlStatement, name).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | Name: " + name)
    }
    defer db.Close()
    http.Redirect(w, r, "/listInsulins", 301)
}

func UpdateInsulin(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Update Insulin")
    if r.Method == "POST" {
        name := r.FormValue("name")
        id := r.FormValue("uid")
        sqlStatement := "UPDATE Insulins SET name=$1 WHERE id=$2"
        updtForm, err := db.Prepare(sqlStatement)
        if err != nil {
            panic(err.Error())
        }    
        updtForm.Exec(name, id)    
        log.Println("UPDATE: Id: " + id +" | Name: " + name ) 
    }
    defer db.Close()
    http.Redirect(w, r, "/listInsulins", 301)
}

func DeleteInsulin(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Delete Insulin")
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Insulins WHERE id=$1")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE: Id: " + id)
    defer db.Close()
    http.Redirect(w, r, "/listInsulins", 301)
}

func NewInsulin(w http.ResponseWriter, r *http.Request) {
    log.Println("New Insulin")
	tmpl.ExecuteTemplate(w, "NewInsulin", nil)
}