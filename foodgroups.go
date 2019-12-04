package main

import (	
	"log"
	"net/http"
	"database/sql"
	"strconv"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type FoodGroup struct {
    Id    int
    Name  string
    Selected bool
}
func InitFoodsGroupsTable(db *sql.DB) {
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Foods_Groups ( " +
		" id SERIAL PRIMARY KEY, "+
		" name varchar(20) NOT NULL " +
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListFoodsGroups(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	log.Println("Index")
	selDB, err := db.Query("SELECT * FROM Foods_Groups ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	group := FoodGroup{}
    groups := []FoodGroup{}
	for selDB.Next() {
		var id int
      var name string
      err = selDB.Scan(&id, &name)
      if err != nil {
         panic(err.Error())
      }
      group.Id = id
      group.Name = name      
      groups = append(groups, group)
}
	tmpl.ExecuteTemplate(w, "ListFoodsGroups", groups) 
	defer db.Close()
}

func ShowFoodGroup(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, name FROM Foods_Groups WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    food := Food{}
    for selDB.Next() {
        var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
    }
    tmpl.ExecuteTemplate(w, "ShowFood", food)
    defer db.Close()
}

func EditFoodGroup(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, name FROM Foods_Groups WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    food := Food{}
    for selDB.Next() {
        var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
    }
    tmpl.ExecuteTemplate(w, "EditFoodGroup", food)
    defer db.Close()
}

func InsertFoodGroup(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
      name := r.FormValue("name")
		sqlStatement := "INSERT INTO Foods_Groups(name) VALUES ($1) RETURNING id"
		id := 0
		err := db.QueryRow(sqlStatement, name).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | Name: " + name)
    }
    defer db.Close()
    http.Redirect(w, r, "/listF", 301)
}

func UpdateFoodGroup(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
      name := r.FormValue("name")
		id := r.FormValue("uid")
		sqlStatement := "UPDATE Foods SET name=$1 WHERE id=$2"
		updtForm, err := db.Prepare(sqlStatement)
        if err != nil {
            panic(err.Error())
		}    
		updtForm.Exec(name, id)    
      log.Println("UPDATE: Id: " + id +" | Name: " + name ) 
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func DeleteFoodGroup(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Foods_Groups WHERE id=$1")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE: Id: " + id)
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func NewFoodGroup(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "NewFoodGroup", nil)
}