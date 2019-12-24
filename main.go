package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	_ "github.com/lib/pq"
)

var tmpl = template.Must(template.ParseGlob("form/*"))

func InitDB(db *sql.DB) {
	InitUnitsTable(db)
	InitInsulinsTable(db)
	InitMealsTable(db)
	InitFoodsGroupsTable(db)
	InitFoodsTable(db)
	InitMeasuresTable(db)
	InitRecordsTable(db)
	InitItemsTable(db)
	
	log.Println("InitDB Sucesso")
}

func dbConn() (db *sql.DB) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	return db
}

func main() {

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
	http.HandleFunc("/listMeasures", ListMeasures)
	http.HandleFunc("/listItems", ListItems)
	http.HandleFunc("/listRecords", ListRecords)
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
	http.HandleFunc("/newMeasure", NewMeasure)
	http.HandleFunc("/showMeasure", ShowMeasure)
	http.HandleFunc("/editMeasure", EditMeasure)
	http.HandleFunc("/insertMeasure", InsertMeasure)
	http.HandleFunc("/updateMeasure", UpdateMeasure)
	http.HandleFunc("/deleteMeasure", DeleteMeasure)
	http.HandleFunc("/newItem", NewItem)
	http.HandleFunc("/showItem", ShowItem)
	http.HandleFunc("/editItem", EditItem)
	http.HandleFunc("/insertItem", InsertItem)
	http.HandleFunc("/updateItem", UpdateItem)
	http.HandleFunc("/deleteItem", DeleteItem)
	http.HandleFunc("/newRecord", NewRecord)
	http.HandleFunc("/showRecord", ShowRecord)
	http.HandleFunc("/editRecord", EditRecord)
	http.HandleFunc("/insertRecord", InsertRecord)
	http.HandleFunc("/updateRecord", UpdateRecord)
	http.HandleFunc("/deleteRecord", DeleteRecord)

	http.HandleFunc("/calculate", Calculate)
	http.HandleFunc("/storeRecordInSession", StoreRecordInSession)

	http.ListenAndServe(":5000", nil)

	defer database.Close()

}

func Calculate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("Calculate")
	foodId := r.URL.Query().Get("foodid")
	quantityItem, err := strconv.ParseFloat(r.URL.Query().Get("quantity"), 64)
	sqlStatement := "SELECT quantity, CHO from Measures WHERE food_id = $1"
	selDB, err := db.Query(sqlStatement, foodId)
	log.Println(sqlStatement + " - " + foodId)
	if err != nil {
		panic(err.Error())
	}
	for selDB.Next() {
		var quantity, CHO, CHOitem float64
		err = selDB.Scan(&quantity, &CHO)
		if err != nil {
			panic(err.Error())
		}
		CHOitem = CHO * quantityItem / quantity
		log.Println("CHO: " + strconv.FormatFloat(CHOitem, 'f', 6, 64))
		io.WriteString(w, fmt.Sprintf("%.2f", CHOitem))
	}
	defer db.Close()
}
