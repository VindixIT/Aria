package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

type Record struct {
	Id             int
	Meal           string
	MealId         int
	MealName       string
	MealOptions    []Meal
	Insulin        string
	InsulinId      int
	InsulinName    string
	InsulinOptions []Insulin
	Gbm            float64
	Gam            float64
	Dose           float64
	CHO            float64
	Created        time.Time
	Items          []Item
}

func InitSessionCookies(rw http.ResponseWriter, request *http.Request, cookieName string) {
	session, _ := store.Get(request, cookieName)
	session.Options.MaxAge = -1
	sessions.Save(request, rw)
}

func StoreRecordInSession(w http.ResponseWriter, r *http.Request) {
	log.Println("Store Record in Session")
	if r.Method == "POST" {
		mealid, _ := strconv.Atoi(r.FormValue("mealid"))
		insulinid, _ := strconv.Atoi(r.FormValue("insulinid"))
		gbm, _ := strconv.ParseFloat(r.FormValue("gbm"), 64)
		gam, _ := strconv.ParseFloat(r.FormValue("gam"), 64)
		dose, _ := strconv.ParseFloat(r.FormValue("dose"), 64)
		session, _ := store.Get(r, "mysession")
		// sessionRecord := session.Values["myRecord"]
		newRecord := Record{
			MealId:    mealid,
			InsulinId: insulinid,
			Gbm:       gbm,
			Gam:       gam,
			Dose:      dose,
		}
		bytesRecord, _ := json.Marshal(newRecord)
		session.Values["myRecord"] = string(bytesRecord)
	}
	sessions.Save(r, w)
	http.Redirect(w, r, "/newRecord", 301)

}

func InitRecordsTable(db *sql.DB) {
	log.Println("Init Records")
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Records ( " +
			" id SERIAL PRIMARY KEY, " +
			" meal_id integer references Meals, " +
			" insulin_id integer references Insulins, " +
			" gbm DOUBLE PRECISION NOT NULL, " +
			" gam DOUBLE PRECISION NOT NULL, " +
			" dose DOUBLE PRECISION NOT NULL, " +
			" creation_date TIMESTAMP without time zone NOT NULL " +
			" )"); err != nil {
		log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListRecords(w http.ResponseWriter, r *http.Request) {
	InitSessionCookies(w, r, "mysession")
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
		err = selDB.Scan(&id, &mealid, &mealname, &insulinid, &insulinname, &gbm, &gam, &dose, &created)
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
		record.Created = created
		res = append(res, record)
	}
	tmpl.ExecuteTemplate(w, "ListRecords", res)
	defer db.Close()
}

func ShowRecord(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("Show Record")
	nId := r.URL.Query().Get("id")
	session, _ := store.Get(r, "mysession")
	sqlStatement := "SELECT " +
		" A.id, A.meal_id, B.name as meal_name, A.insulin_id, D.name as insulin_name, A.gbm, A.gam, A.dose, A.creation_date " +
		" FROM public.records A " +
		" LEFT JOIN meals B on A.meal_id = B.ID " +
		" LEFT JOIN insulins D on A.insulin_id = D.id " +
		" WHERE a.id = $1 ORDER BY id DESC "
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
		err = selDB.Scan(&id, &mealid, &mealname, &insulinid, &insulinname, &gbm, &gam, &dose, &created)
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
		record.Created = created
	}
	sqlStatement = "SELECT " +
		"	a.id," +
		"	a.food_id," +
		"	b.name as food_name," +
		"	a.unit_id," +
		"	c.symbol as unit_symbol," +
		"	a.quantity," +
		"	a.cho " +
		" from " +
		"	items a" +
		" left join " +
		"	foods b " +
		" on a.food_id = b.id" +
		" left join " +
		"	units c " +
		" on a.unit_id = c.id" +
		" where record_id = $1"
	log.Println(sqlStatement)
	selDB, err = db.Query(sqlStatement, nId)
	if err != nil {
		panic(err.Error())
	}
	item := Item{}
	for selDB.Next() {
		var id, foodid, unitid int
		var foodName, unitSymbol string
		var quantity, CHO float64
		err = selDB.Scan(&id, &foodid, &foodName, &unitid, &unitSymbol, &quantity, &CHO)
		if err != nil {
			panic(err.Error())
		}
		item.Id = strconv.Itoa(id)
		item.FoodId = foodid
		item.UnitId = unitid
		item.FoodName = foodName
		item.UnitSymbol = unitSymbol
		item.Quantity = quantity
		item.CHO = CHO
		record.Items = append(record.Items, item)
	}
	bytesItems, _ := json.Marshal(record.Items)
	session.Values["myitems"] = string(bytesItems)
	sessions.Save(r, w)
	tmpl.ExecuteTemplate(w, "ShowRecord", record)
	defer db.Close()
}

func InsertRecord(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("Insert Record")
	recordid := 0
	if r.Method == "POST" {
		mealid := r.FormValue("mealid")
		insulinid := r.FormValue("insulinid")
		gbm := r.FormValue("gbm")
		gam := r.FormValue("gam")
		dose := r.FormValue("dose")
		created := time.Now()
		formatedTime := created.Format(time.RFC1123)
		log.Println("INSERT: MealId: " + mealid + " | InsulinId: " + insulinid + " | GBM: " + gbm + " | GAM: " + gam + " | Dose: " + dose + " | Created: " + formatedTime)
		sqlStatement := "INSERT INTO Records(meal_id, insulin_id, gbm, gam, dose, creation_date) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
		err := db.QueryRow(sqlStatement, mealid, insulinid, gbm, gam, dose, created).Scan(&recordid)
		if err != nil {
			panic(err.Error())
		}
	}
	session, _ := store.Get(r, "mysession")
	strItems := session.Values["myitems"].(string)
	myItems := []Item{}
	json.Unmarshal([]byte(strItems), &myItems)
	for index := range myItems {
		item := myItems[index]
		log.Println("FoodId: " + strconv.Itoa(item.FoodId))
		log.Println("FoodName: " + item.FoodName)
		log.Println("UnitId: " + strconv.Itoa(item.UnitId))
		log.Println("UnitSymbol: " + item.UnitSymbol)
		log.Println("Quantity: " + fmt.Sprintf("%f", item.Quantity))
		log.Println("CHO: " + fmt.Sprintf("%f", item.CHO))
		id := 0
		sqlStatement := "INSERT INTO Items(record_id, food_id, unit_id, quantity, cho) VALUES ($1,$2,$3,$4,$5) RETURNING id"
		err := db.QueryRow(sqlStatement, recordid, item.FoodId, item.UnitId, item.Quantity, item.CHO).Scan(&id)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
	http.Redirect(w, r, "/listRecords", 301)
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("Update Record")
	if r.Method == "POST" {
		mealid := r.FormValue("mealid")
		insulinid := r.FormValue("insulinid")
		gbm := r.FormValue("gbm")
		gam := r.FormValue("gam")
		dose := r.FormValue("dose")
		id := r.FormValue("uid")
		sqlStatement := "UPDATE Records SET meal_id=$1, insulin_id=$2, gbm=$3, gam=$4, dose=$5 WHERE id=$6"
		updtForm, err := db.Prepare(sqlStatement)
		if err != nil {
			panic(err.Error())
		}
		updtForm.Exec(mealid, insulinid, gbm, gam, dose, id)
		log.Println("UPDATE: Id: " + id + " | MealId: " + mealid + " | InsulinId: " + insulinid + " | GBM: " + gbm + " | GAM: " + gam + " | Dose: " + dose)
	}
	defer db.Close()
	http.Redirect(w, r, "/listRecords", 301)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("Delete Items")
	id := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Items WHERE record_id=$1")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)
	log.Println("Delete Record")
	delForm, err = db.Prepare("DELETE FROM Records WHERE id=$1")
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
	session, _ := store.Get(r, "mysession")
	sessionItem := session.Values["myitems"]
	sessionRecord := session.Values["myRecord"]
	if sessionRecord != nil {
		strRecord := session.Values["myRecord"].(string)
		json.Unmarshal([]byte(strRecord), &record)
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
	selInsulinsDB, err := db.Query("SELECT id, name FROM Insulins")
	if err != nil {
		panic(err.Error())
	}
	insulin := Insulin{}
	insulins := []Insulin{}
	for selInsulinsDB.Next() {
		var id int
		var name string
		err = selInsulinsDB.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		insulin.Id = id
		insulin.Name = name
		if record.InsulinId == id {
			insulin.Selected = true
		} else {
			insulin.Selected = false
		}
		insulins = append(insulins, insulin)
	}
	record.InsulinOptions = insulins
	myItems := []Item{}
	if sessionItem != nil {
		strItems := session.Values["myitems"].(string)
		json.Unmarshal([]byte(strItems), &myItems)
		record.Items = myItems
	}
	tmpl.ExecuteTemplate(w, "NewRecord", record)
}

func EditRecord(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("Edit Record")
	nId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	log.Println("id = ", nId)
	myItems := []Item{}
	record := Record{}
	session, _ := store.Get(r, "mysession")
	sessionItem := session.Values["myitems"]
	sessionRecord := session.Values["myRecord"]
	if sessionItem != nil {
		strItems := session.Values["myitems"].(string)
		json.Unmarshal([]byte(strItems), &myItems)
		record.Items = myItems
	}
	if sessionRecord != nil {
		strRecord := session.Values["myRecord"].(string)
		json.Unmarshal([]byte(strRecord), &record)

	} else {
		sqlStatement := "SELECT " +
			" A.id, A.meal_id, B.name as meal_name, A.insulin_id, D.name as insulin_name, A.gbm, A.gam, A.dose, A.creation_date " +
			" FROM public.records A " +
			" LEFT JOIN meals B on A.meal_id = B.ID " +
			" LEFT JOIN insulins D on A.insulin_id = D.id " +
			" WHERE a.id = $1"
		log.Println(sqlStatement)
		selDB, err := db.Query(sqlStatement, nId)
		if err != nil {
			panic(err.Error())
		}
		for selDB.Next() {
			var id, mealid, insulinid int
			var mealname, insulinname string
			var gbm, gam, dose float64
			var created time.Time
			err = selDB.Scan(&id, &mealid, &mealname, &insulinid, &insulinname, &gbm, &gam, &dose, &created)
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
			record.Created = created
		}
		sqlStatement = "SELECT " +
			"	a.id," +
			"	a.food_id," +
			"	b.name as food_name," +
			"	a.unit_id," +
			"	c.symbol as unit_symbol," +
			"	a.quantity," +
			"	a.cho " +
			" from " +
			"	items a" +
			" left join " +
			"	foods b " +
			" on a.food_id = b.id" +
			" left join " +
			"	units c " +
			" on a.unit_id = c.id" +
			" where record_id = $1"
		log.Println(sqlStatement)
		selDB, err = db.Query(sqlStatement, nId)
		if err != nil {
			panic(err.Error())
		} //megadeath
		item := Item{}
		for selDB.Next() {
			var id, foodid, unitid int
			var foodName, unitSymbol string
			var quantity, CHO float64
			err = selDB.Scan(&id, &foodid, &foodName, &unitid, &unitSymbol, &quantity, &CHO)
			if err != nil {
				panic(err.Error())
			}
			item.Id = strconv.Itoa(id)
			item.FoodId = foodid
			item.UnitId = unitid
			item.FoodName = foodName
			item.UnitSymbol = unitSymbol
			item.Quantity = quantity
			item.CHO = CHO
			record.Items = append(record.Items, item)
		}
		bytesItems, _ := json.Marshal(record.Items)
		session.Values["myitems"] = string(bytesItems)
		sessions.Save(r, w)
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
	selInsulinsDB, err := db.Query("SELECT id, name FROM Insulins")
	if err != nil {
		panic(err.Error())
	}
	insulin := Insulin{}
	insulins := []Insulin{}
	for selInsulinsDB.Next() {
		var id int
		var name string
		err = selInsulinsDB.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		insulin.Id = id
		insulin.Name = name
		if record.InsulinId == id {
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
