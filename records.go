package main

import (	
	"log"
	"net/http"
	"database/sql"
    "strconv"
    "time"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type Record struct {
    Id    int
    Meal  string
    MealId  int
    MealName  string
    MealOptions []Meal
    Insulin  string
    InsulinId  int
    InsulinName  string
    InsulinOptions []Insulin
    Gbm float64 
    Gam float64
    Dose float64 
    Created time.Time
}
func InitRecordsTable(db *sql.DB) {
    log.Println("Init Records") 
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Records ( " +
		" id SERIAL PRIMARY KEY, "+
		" meal_id integer references Meals, "+
        " insulin_id integer references Insulins, "+
        " gbm DOUBLE PRECISION NOT NULL, " +
        " gam DOUBLE PRECISION NOT NULL, " + 
        " dose DOUBLE PRECISION NOT NULL, " +
        " creation_date DOUBLE PRECISION NOT NULL " +  
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListRecords(w http.ResponseWriter, r *http.Request){
    db := dbConn()
    log.Println("List Records")    
    selDB, err := db.Query("SELECT " + 
    " A.id, A.meal_id, B.name as meal_name, A.insulin_id, D.name as insulin_name, A.gbm, A.gam, A.dose, A.creation_date " +
    " FROM public.records A " +
    " LEFT JOIN meals B on A.meal_id = B.ID " +
    " LEFT JOIN insulins D on A.insulin_id = D.id " +
    " ORDER BY id DESC")
    if err != nil {         
        panic(err.Error())
    }
	record := Record{} 
    res := []Record{} 
	for selDB.Next() {
        var id, mealid, insulinid int
        var mealname, insulinname string
        var gbm, gam, dose float64        
        var created time.Time
        err = selDB.Scan(&id, &mealid, &mealname, &insulinid, &insulinname, &gbm, &gam, &dose)
        if err != nil {
            panic(err.Error())
        }
        record.Id = id
        record.MealId = mealid
        record.MealName = mealname
        record.InsulinId = insulinid
        record.InsulinName = insulinname
        record.Gbm = gbm
        record.Gam = gam
        record.Dose = dose
        res = append(res, record)
    }
	tmpl.ExecuteTemplate(w, "ListRecords", res)
	defer db.Close()
}

func ShowRecord(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Show Record")
    nId := r.URL.Query().Get("id")
    sqlStatement := "SELECT " + 
    " A.id, A.meal_id, B.name as meal_name, A.insulin_id, D.name as insulin_name, A.gbm, A.gam, A.dose, A.creation_date " +
    " FROM public.records A " +
    " LEFT JOIN meals B on A.meal_id = B.ID " +
    " LEFT JOIN insulins D on A.insulin_id = D.id " +
    " ORDER BY id DESC WHERE a.id = $1"
    log.Println(sqlStatement)
    selDB, err := db.Query(sqlStatement, nId)
    if err != nil {
        panic(err.Error())
    }
    record := Record{}
    for selDB.Next() {
        var id, mealid, insulinid int
        var mealname, insulinname string
        var gbm, gam, dose float64        
        var created time.Time
        err = selDB.Scan(&id, &mealid, &mealname, &insulinid, &insulinname, &gbm, &gam, &dose)
        if err != nil {
            panic(err.Error())
        }
        record.Id = id
        record.MealId = mealid
        record.MealName = mealname
        record.InsulinId = insulinid
        record.InsulinName = insulinname
        record.Gbm = gbm
        record.Gam = gam
        record.Dose = dose
    }
    tmpl.ExecuteTemplate(w, "ShowRecord", record)
    defer db.Close()
}

func EditRecord(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Edit Record")
    nId := r.URL.Query().Get("id")
    sqlStatement := "SELECT " + 
    " A.id, A.meal_id, B.name as meal_name, A.insulin_id, D.name as insulin_name, A.gbm, A.gam, A.dose, A.creation_date " +
    " FROM public.records A " +
    " LEFT JOIN meals B on A.meal_id = B.ID " +
    " LEFT JOIN insulins D on A.insulin_id = D.id " +
    " ORDER BY id DESC WHERE a.id = $1"
    log.Println(sqlStatement)
    selDB, err := db.Query(sqlStatement, nId)
    if err != nil {
        panic(err.Error())
    }
    record := Record{}
    for selDB.Next() {
        var id, mealid, insulinid int
        var mealname, insulinname string
        var gbm, gam, dose float64        
        var created time.Time
        err = selDB.Scan(&id, &mealid, &mealname, &insulinid, &insulinname, &gbm, &gam, &dose)
        if err != nil {
            panic(err.Error())
        }
        record.Id = id
        record.MealId = mealid
        record.MealName = mealname
        record.InsulinId = insulinid
        record.InsulinName = insulinname
        record.Gbm = gbm
        record.Gam = gam
        record.Dose = dose
    }
    selMealsDB, err := db.Query("SELECT id, name FROM Meals")    
    if err != nil {
        panic(err.Error())
    }
    meal := Meal{}
    meals := []Meal{}
    for selMealsDB.Next() {
        var id int
        var name string
        err = selMealsDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        meal.Id = id
        meal.Name = name
        if record.MealId == id {
            meal.Selected = true
        } else {
            meal.Selected = false
        }
        meals = append(meals, meal)
    }
    record.MealOptions = meals
    selInsulinsDB, err := db.Query("SELECT id, name FROM Inslulins")
    if err != nil {
        panic(err.Error())
    }
    insulin := Unit{}
    insulins := []Unit{}
    for selInsulinsDB.Next() {
        var id int
        var name string
        err = selInsulinsDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        insulin.Id = id
        insulin.Name = name
        if record.UnitId  == id {
            insulin.Selected = true
        } else {
            insulin.Selected = false
        }
        insulins = append(insulins, insulin)
    }
    record.InsulinOptions = insulins
    tmpl.ExecuteTemplate(w, "EditRecord", record)
    defer db.Close()
}

func InsertRecord(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Insert Record")
    if r.Method == "POST" {
        foodid := r.FormValue("foodid")
        unitid := r.FormValue("unitid")
        quantity := r.FormValue("quantity")
        CHO := r.FormValue("CHO")
        id := 0
        log.Println("INSERT: FoodId: " + foodid + " | UnitId: " + unitid + " | Quantity: " + quantity + " | CHO: " + CHO)
        sqlStatement := "INSERT INTO Records(food_id, unit_id, quantity, CHO) VALUES ($1,$2,$3,$4) RETURNING id"
		err := db.QueryRow(sqlStatement,foodid,unitid,quantity,CHO).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | FoodId: " + foodid + " | UnitId: " + unitid + " | Quantity: " + quantity + " | CHO: " + CHO)
    }
    defer db.Close()
    http.Redirect(w, r, "/listRecords", 301)
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Update Record")
    if r.Method == "POST" {
        foodid := r.FormValue("foodid")
        unitid := r.FormValue("unit")
        quantity := r.FormValue("quantity")
        CHO := r.FormValue("CHO")
		id := r.FormValue("uid")
		sqlStatement := "UPDATE Records SET food_id=$1, unit_id=$2, quantity=$3, cho=$4 WHERE id=$5"
		updtForm, err := db.Prepare(sqlStatement)
        if err != nil {
            panic(err.Error())
		}    
		updtForm.Exec(foodid, unitid, quantity, CHO, id)
        log.Println("UPDATE: Id: " + id +" | FoodId: " + foodid + " | UnitId: " + unitid + " | Quantity: " + quantity + " | CHO: " + CHO)
    }
    defer db.Close()
    http.Redirect(w, r, "/listRecords", 301)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Delete Record")
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Records WHERE id=$1")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE: Id: " + id)
    defer db.Close()
    http.Redirect(w, r, "/listRecords", 301)
}

func NewRecord(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("New Record")
    record := Record{}
    selFoodsDB, err := db.Query("SELECT id, name FROM Foods")
    if err != nil {
        panic(err.Error())
    }
    food := Food{}
    foods := []Food{}
    for selFoodsDB.Next() {
        var id int
        var name string
        err = selFoodsDB.Scan(&id, &name)
        if err != nil {
            panic(err.Error())
        }
        food.Id = id
        food.Name = name
        foods = append(foods, food)
    }
    record.FoodOptions = foods
    selUnitsDB, err := db.Query("SELECT id, symbol, description FROM Units")
    if err != nil {
        panic(err.Error())
    }
    unit := Unit{}
    units := []Unit{}
    for selUnitsDB.Next() {
        var id int
        var symbol, description string
        err = selUnitsDB.Scan(&id, &symbol, &description)
        if err != nil {
            panic(err.Error())
        }
        unit.Id = id
        unit.Symbol = symbol
        unit.Description = description
        units = append(units, unit)
    }
    record.UnitOptions = units
	tmpl.ExecuteTemplate(w, "NewRecord", record)
}