package main

import (	
	"log"
	"database/sql"
    "time"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)
type Record struct {
    Id    int
    Meal  string
    MealId  int
    MealName  string
    MealOptions []Food
    Insulin  string
    InsulinId  int
    InsulinName  string
    InsulinOptions []Food
    gbm float64 
    gam float64
    dose float64 
    created time.Time
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