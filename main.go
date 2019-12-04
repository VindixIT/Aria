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
	InitFoodsGroupsTable(db) // a ordem é importante.
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

	http.HandleFunc("/", ListFoodsGroups)
	http.HandleFunc("/listF", ListFoods)
	http.HandleFunc("/listFG", ListFoodsGroups)
	http.HandleFunc("/newF", NewFood)
	http.HandleFunc("/showF", ShowFood)
	http.HandleFunc("/editF", EditFood)
	http.HandleFunc("/insertF", InsertFood)
	http.HandleFunc("/updateF", UpdateFood)
	http.HandleFunc("/deleteF", DeleteFood)
	http.HandleFunc("/newFG", NewFoodGroup)
	http.HandleFunc("/showFG", ShowFoodGroup)
	http.HandleFunc("/editFG", EditFoodGroup)
	http.HandleFunc("/insertFG", InsertFoodGroup)
	http.HandleFunc("/updateFG", UpdateFoodGroup)
	http.HandleFunc("/deleteFG", DeleteFoodGroup)

	http.ListenAndServe(":5000", nil)
	
	defer database.Close()

}



