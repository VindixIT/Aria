package main

import (	
	"log"
	"net/http"
	"database/sql"
	"strconv"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type Meal struct {
    Id    int
    Name  string
    Selected bool
}
func InitMealsTable(db *sql.DB) {
    log.Println("Init Meals")
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Meals ( " +
		" id SERIAL PRIMARY KEY, "+
		" name varchar(20) NOT NULL " +
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListMeals(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	log.Println("List Meals")
	selDB, err := db.Query("SELECT * FROM Meals ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	meal := Meal{}
    meals := []Meal{}
	for selDB.Next() {
		var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        meal.Id = id
        meal.Name = name      
        meals = append(meals, meal)
    }
	tmpl.ExecuteTemplate(w, "ListMeals", meals) 
	defer db.Close()
}

func ShowMeal(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Show Meal")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, name FROM Meals WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    meal := Meal{}
    for selDB.Next() {
        var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        meal.Id = id
        meal.Name = name
    }
    tmpl.ExecuteTemplate(w, "ShowMeal", meal)
    defer db.Close()
}

func EditMeal(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Edit Meal")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, name FROM Meals WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    meal := Meal{}
    for selDB.Next() {
        var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        meal.Id = id
        meal.Name = name
    }
    tmpl.ExecuteTemplate(w, "EditMeal", meal)
    defer db.Close()
}

func InsertMeal(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Insert Meal")
    if r.Method == "POST" {
        name := r.FormValue("name")
		sqlStatement := "INSERT INTO Meals(name) VALUES ($1) RETURNING id"
		id := 0
		err := db.QueryRow(sqlStatement, name).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | Name: " + name)
    }
    defer db.Close()
    http.Redirect(w, r, "/listMeals", 301)
}

func UpdateMeal(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Update Meal")
    if r.Method == "POST" {
        name := r.FormValue("name")
		id := r.FormValue("uid")
		sqlStatement := "UPDATE Meals SET name=$1 WHERE id=$2"
		updtForm, err := db.Prepare(sqlStatement)
        if err != nil {
            panic(err.Error())
		}    
		updtForm.Exec(name, id)    
        log.Println("UPDATE: Id: " + id +" | Name: " + name ) 
    }
    defer db.Close()
    http.Redirect(w, r, "/listMeals", 301)
}

func DeleteMeal(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Delete Meal")
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Meals WHERE id=$1")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE: Id: " + id)
    defer db.Close()
    http.Redirect(w, r, "/listMeals", 301)
}

func NewMeal(w http.ResponseWriter, r *http.Request) {
    log.Println("New Meal")
	tmpl.ExecuteTemplate(w, "NewMeal", nil)
}