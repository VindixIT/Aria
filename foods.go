package main

import (	
	"log"
	"net/http"
	"database/sql"
	"strconv"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type Food struct {
    Id    int
    Group  string
    Name string
}
func InitFoodsTable(db *sql.DB) {
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS foods ( " +
		" id SERIAL PRIMARY KEY, "+
		" grp varchar(20) NOT NULL, "+
		" name varchar(255) NOT NULL " +
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListFoods(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	log.Println("Index")
	selDB, err := db.Query("SELECT * FROM Foods ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	food := Food{}
    res := []Food{}
	for selDB.Next() {
		var id int
        var group, name string
        err = selDB.Scan(&id, &group, &name)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
        food.Group = group
        res = append(res, food)
	}
	tmpl.ExecuteTemplate(w, "ListFoods", res)
	defer db.Close()
}

func ShowFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Foods WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    food := Food{}
    for selDB.Next() {
        var id int
        var name, group string
        err = selDB.Scan(&id, &name, &group)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
        food.Group = group
    }
    tmpl.ExecuteTemplate(w, "ShowFood", food)
    defer db.Close()
}


func InsertFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        group := r.FormValue("group")
		sqlStatement := "INSERT INTO Foods(name,grp) VALUES ($1,$2) RETURNING id"
		id := 0
		err := db.QueryRow(sqlStatement, name, group).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | Name: " + name + " | Group: " + group)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func NewFood(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "NewFood", nil)
}