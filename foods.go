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
    GroupId  int
    GroupOptions []FoodGroup
    Name string
}
func InitFoodsTable(db *sql.DB) {
    log.Println("Init Foods")
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Foods ( " +
		" id SERIAL PRIMARY KEY, "+
		" grp_id integer references Foods_Groups, "+
		" name varchar(255) NOT NULL " +
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListFoods(w http.ResponseWriter, r *http.Request){
    db := dbConn()
    log.Println("List Foods")
    selDB, err := db.Query("SELECT A.id, A.name, B.id, B.name as grp_name "+
    " FROM Foods A left join Foods_Groups B "+
    " on A.grp_id = B.id ORDER BY a.id DESC")
    if err != nil {
        panic(err.Error())
	}
	food := Food{}
    res := []Food{}
	for selDB.Next() {
        var id, groupid int
        var group, name string
        err = selDB.Scan(&id, &name, &groupid, &group)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
        food.Group = group
        food.GroupId = groupid
        res = append(res, food)
	}
	tmpl.ExecuteTemplate(w, "ListFoods", res)
	defer db.Close()
}

func ShowFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Show Food")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT a.id, a.name, b.id, b.name as grp_name FROM Foods a"+
    " left join Foods_Groups b on a.grp_id = b.id WHERE a.id=$1 ORDER BY a.id DESC", nId)
    if err != nil {
        panic(err.Error())
    }
    food := Food{}
    for selDB.Next() {
        var id, groupid int
        var name, group string
        err = selDB.Scan(&id, &name, &groupid, &group)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
        food.Group = group
        food.GroupId = groupid
    }
    tmpl.ExecuteTemplate(w, "ShowFood", food)
    defer db.Close()
}

func EditFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Edit Food")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT id, name, grp_id FROM Foods WHERE id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    food := Food{}
    for selDB.Next() {
        var id, groupid  int
        var name string
        err = selDB.Scan(&id, &name, &groupid)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
        food.GroupId = groupid
    }
    selGroupsDB, err := db.Query("SELECT id, name FROM Foods_Groups")
    if err != nil {
        panic(err.Error())
    }
    foodGroup := FoodGroup{}
    groups := []FoodGroup{}
    for selGroupsDB.Next() {
        var id int
        var name string
        err = selGroupsDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        foodGroup.Id = id
        foodGroup.Name = name
        if food.GroupId == id {
            foodGroup.Selected = true
        } else {
            foodGroup.Selected = false
        }
        groups = append(groups, foodGroup)
    }
    food.GroupOptions = groups
    tmpl.ExecuteTemplate(w, "EditFood", food)
    defer db.Close()
}

func InsertFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Insert Food")
    if r.Method == "POST" {
        name := r.FormValue("name")
        group := r.FormValue("group")
		sqlStatement := "INSERT INTO Foods(name,grp_id) VALUES ($1,$2) RETURNING id"
		id := 0
		err := db.QueryRow(sqlStatement, name, group).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | Name: " + name + " | Group: " + group)
    }
    defer db.Close()
    http.Redirect(w, r, "/listFoods", 301)
}

func UpdateFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Update Food")
    if r.Method == "POST" {
        name := r.FormValue("name")
        group := r.FormValue("group")
		id := r.FormValue("uid")
		sqlStatement := "UPDATE Foods SET name=$1, grp_id=$2 WHERE id=$3"
		updtForm, err := db.Prepare(sqlStatement)
        if err != nil {
            panic(err.Error())
		}    
		updtForm.Exec(name, group, id)    
        log.Println("UPDATE: Id: " + id +" | Name: " + name + " | Group: " + group)
    }
    defer db.Close()
    http.Redirect(w, r, "/listFoods", 301)
}

func DeleteFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Delete Food")
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Foods WHERE id=$1")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE: Id: " + id)
    defer db.Close()
    http.Redirect(w, r, "/listFoods", 301)
}

func NewFood(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("New Food")
    selDB, err := db.Query("SELECT id, name FROM Foods_Groups")
    if err != nil {
        panic(err.Error())
    }
    foodGroup := FoodGroup{}
    groups := []FoodGroup{}
    for selDB.Next() {
        var id int
        var name string
        err = selDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        foodGroup.Id = id
        foodGroup.Name = name
        groups = append(groups, foodGroup)
    }    
	tmpl.ExecuteTemplate(w, "NewFood", groups)
}