package main

import (	
	"log"
	"net/http"
	"database/sql"
	"strconv"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type Measure struct {
    Id    int
    Food  string
    FoodId  int
    FoodName  string
    FoodOptions []Food
    Unit  string
    UnitId  int
    UnitSymbol string
    UnitOptions []Unit
    Quantity float64
    CHO float64 
}
func InitMeasuresTable(db *sql.DB) {
    log.Println("Init Measures") 
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Measures ( " +
		" id SERIAL PRIMARY KEY, "+
		" food_id integer references Foods, "+
		" unit_id integer references Units, "+
		" quantity DOUBLE PRECISION NOT NULL, " +
		" CHO DOUBLE PRECISION NOT NULL " +
		" )"); err != nil {
			log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListMeasures(w http.ResponseWriter, r *http.Request){
    db := dbConn()
    log.Println("List Measures")
    selDB, err := db.Query("SELECT "+
    " A.id, B.id, B.name AS food_name, C.id, C.symbol AS unit_symbol, A.quantity, A.CHO "+
    " FROM Measures A left join Foods B "+
    " on A.food_id = B.id left join Units C on A.unit_id = C.id ORDER BY a.id DESC")
    if err != nil {         
        panic(err.Error())
	}
	measure := Measure{}
    res := []Measure{}
	for selDB.Next() {
        var id, foodid, unitid int
        var foodName, unitSymbol string
        var quantity, CHO float64        
        err = selDB.Scan(&id, &foodid, &foodName, &unitid, &unitSymbol, &quantity, &CHO)
        if err != nil {
            panic(err.Error())
        }
        measure.Id = id
        measure.FoodId = foodid
        measure.UnitId = unitid
        measure.FoodName = foodName
        measure.UnitSymbol = unitSymbol
        measure.Quantity = quantity
        measure.CHO = CHO
        res = append(res, measure)
	}
	tmpl.ExecuteTemplate(w, "ListMeasures", res)
	defer db.Close()
}

func ShowMeasure(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Show Measure")
    nId := r.URL.Query().Get("id")
    sqlStatement := "SELECT "+
    " A.id, B.id, B.name AS food_name, C.id, C.symbol AS unit_symbol, A.quantity, A.CHO "+
    " FROM Measures A left join Foods B "+
    " on A.food_id = B.id left join Units C on A.unit_id = C.id WHERE a.id = $1"
    log.Println(sqlStatement)
    selDB, err := db.Query(sqlStatement, nId)
    if err != nil {
        panic(err.Error())
    }
    measure := Measure{}
    for selDB.Next() {
        var id, foodid, unitid int
        var foodName, unitSymbol string
        var quantity, CHO float64    
        err = selDB.Scan(&id, &foodid, &foodName, &unitid, &unitSymbol, &quantity, &CHO)
        if err != nil {
            panic(err.Error())
        }
        measure.Id = id
        measure.FoodId = foodid
        measure.UnitId = unitid
        measure.FoodName = foodName
        measure.UnitSymbol = unitSymbol
        measure.Quantity = quantity
        measure.CHO = CHO
    }
    tmpl.ExecuteTemplate(w, "ShowMeasure", measure)
    defer db.Close()
}

func EditMeasure(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Edit Measure")
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT "+
    " A.id, B.id, C.id, A.quantity, A.CHO "+
    " FROM Measures A left join Foods B "+
    " on A.food_id = B.id left join Units C on A.unit_id = C.id WHERE a.id=$1", nId)
    if err != nil {
        panic(err.Error())
    }
    measure := Measure{}
    for selDB.Next() {
        var id, foodid, unitid int
        var foodName, unitSymbol string
        var quantity, CHO float64  
        err = selDB.Scan(&id, &foodid, &unitid, &quantity, &CHO)
        if err != nil {
            panic(err.Error())
        }
        measure.Id = id
        measure.FoodId = foodid
        measure.UnitId = unitid
        measure.FoodName = foodName
        measure.UnitSymbol = unitSymbol
        measure.Quantity = quantity
        measure.CHO = CHO
    }
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
        if measure.FoodId == id {
            food.Selected = true
        } else {
            food.Selected = false
        }
        foods = append(foods, food)
    }
    measure.FoodOptions = foods
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
        if measure.UnitId  == id {
            unit.Selected = true
        } else {
            unit.Selected = false
        }
        units = append(units, unit)
    }
    measure.UnitOptions = units
    tmpl.ExecuteTemplate(w, "EditMeasure", measure)
    defer db.Close()
}

func InsertMeasure(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Insert Measure")
    if r.Method == "POST" {
        foodid := r.FormValue("foodid")
        unitid := r.FormValue("unitid")
        quantity := r.FormValue("quantity")
        CHO := r.FormValue("CHO")
        id := 0
        log.Println("INSERT: FoodId: " + foodid + " | UnitId: " + unitid + " | Quantity: " + quantity + " | CHO: " + CHO)
        sqlStatement := "INSERT INTO Measures(food_id, unit_id, quantity, CHO) VALUES ($1,$2,$3,$4) RETURNING id"
		err := db.QueryRow(sqlStatement,foodid,unitid,quantity,CHO).Scan(&id)
        if err != nil {
            panic(err.Error())
        }        
        log.Println("INSERT: Id: " + strconv.Itoa(id) +" | FoodId: " + foodid + " | UnitId: " + unitid + " | Quantity: " + quantity + " | CHO: " + CHO)
    }
    defer db.Close()
    http.Redirect(w, r, "/listMeasures", 301)
}

func UpdateMeasure(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Update Measure")
    if r.Method == "POST" {
        foodid := r.FormValue("food")
        unitid := r.FormValue("unit")
        quantity := r.FormValue("quantity")
        CHO := r.FormValue("cho")
		id := r.FormValue("uid")
		sqlStatement := "UPDATE Measures SET food_id=$1, unit_id=$2, quantity=$3, cho=$4 WHERE id=$5"
		updtForm, err := db.Prepare(sqlStatement)
        if err != nil {
            panic(err.Error())
		}    
		updtForm.Exec(foodid, unitid, quantity, CHO, id)
        log.Println("UPDATE: Id: " + id +" | FoodId: " + foodid + " | UnitId: " + unitid + " | Quantity: " + quantity + " | CHO: " + CHO)
    }
    defer db.Close()
    http.Redirect(w, r, "/listMeasures", 301)
}

func DeleteMeasure(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("Delete Measure")
    id := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Measures WHERE id=$1")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(id)
    log.Println("DELETE: Id: " + id)
    defer db.Close()
    http.Redirect(w, r, "/listMeasures", 301)
}

func NewMeasure(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    log.Println("New Measure")
    measure := Measure{}
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
    measure.FoodOptions = foods
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
    measure.UnitOptions = units
	tmpl.ExecuteTemplate(w, "NewMeasure", measure)
}