package main

import (
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
    "log"
	"os"
	"text/template"   
)

var tmpl = template.Must(template.ParseGlob("form/*"))

func InitDB(db * sql.DB){
	InitUnitsTable(db) 		
	InitInsulinsTable(db) 	
	InitMealsTable(db) 
	InitFoodsGroupsTable(db)
	InitFoodsTable(db)
	log.Println("InitDB Sucesso")
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	return db
}

func main(){
	
	log.Println("Server started on: http://127.0.0.1:5000")
	
	database := dbConn()

	log.Println("database")

	InitDB(database)

	http.HandleFunc("/", ListFoods)
	http.HandleFunc("/listFoods", ListFoods)
	http.HandleFunc("/listFoodsGroups", ListFoodsGroups)
	http.HandleFunc("/listMeals", ListMeals)
	http.HandleFunc("/listInsulins", ListInsulins)
	http.HandleFunc("/listUnits", ListUnits)
	http.HandleFunc("/newFood", NewFood)
	http.HandleFunc("/showFood", ShowFood)
	http.HandleFunc("/editFood", EditFood)
	http.HandleFunc("/insertFood", InsertFood)
	http.HandleFunc("/updateFood", UpdateFood)
	http.HandleFunc("/deleteFood", DeleteFood)
	http.HandleFunc("/newFoodGroup", NewFoodGroup)
	http.HandleFunc("/showFoodGroup", ShowFoodGroup)
	http.HandleFunc("/editFoodGroup", EditFoodGroup)
	http.HandleFunc("/insertFoodGroup", InsertFoodGroup)
	http.HandleFunc("/updateFoodGroup", UpdateFoodGroup)
	http.HandleFunc("/deleteFoodGroup", DeleteFoodGroup)
	http.HandleFunc("/newMeal", NewMeal)
	http.HandleFunc("/showMeal", ShowMeal)
	http.HandleFunc("/editMeal", EditMeal)
	http.HandleFunc("/insertMeal", InsertMeal)
	http.HandleFunc("/updateMeal", UpdateMeal)
	http.HandleFunc("/deleteMeal", DeleteMeal)
	http.HandleFunc("/newInsulin", NewInsulin)
	http.HandleFunc("/showInsulin", ShowInsulin)
	http.HandleFunc("/editInsulin", EditInsulin)
	http.HandleFunc("/insertInsulin", InsertInsulin)
	http.HandleFunc("/updateInsulin", UpdateInsulin)
	http.HandleFunc("/deleteInsulin", DeleteInsulin)
	http.HandleFunc("/newUnit", NewUnit)
	http.HandleFunc("/showUnit", ShowUnit)
	http.HandleFunc("/editUnit", EditUnit)
	http.HandleFunc("/insertUnit", InsertUnit)
	http.HandleFunc("/updateUnit", UpdateUnit)
	http.HandleFunc("/deleteUnit", DeleteUnit)

	http.ListenAndServe(":5000", nil)
	
	defer database.Close()

}



