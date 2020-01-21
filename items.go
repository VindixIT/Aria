package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

var store = sessions.NewCookieStore([]byte("mysession"))

type Item struct {
	Id          string  `json:"id"`
	RecordId    string  `json:"recordid"`
	Food        string  `json:"food"`
	FoodId      int     `json:"foodid"`
	FoodName    string  `json:"foodname"`
	FoodOptions []Food  `json:"foodoptions"`
	Unit        string  `json:"unit"`
	UnitId      int     `json:"unitid"`
	UnitSymbol  string  `json:"unitsymbol"`
	UnitOptions []Unit  `json:"unitoptions"`
	Quantity    float64 `json:"quantity"`
	CHO         float64 `json:"cho"`
}

func InsertItem(rw http.ResponseWriter, request *http.Request) {
	log.Println("*** Insert Item ***")
	if request.Method == "POST" {
		tmp := "tmp"
		foodid, _ := strconv.Atoi(request.FormValue("foodid"))
		foodName := request.FormValue("foodName")
		unitid, _ := strconv.Atoi(request.FormValue("unitid"))
		unitSymbol := request.FormValue("unitSymbol")
		quantity, _ := strconv.ParseFloat(request.FormValue("quantity"), 64)
		CHO, _ := strconv.ParseFloat(request.FormValue("CHO"), 64)
		recordid := request.FormValue("recordid")
		log.Println("FoodName: " + foodName + " | UnitSymbol: " + unitSymbol)
		log.Println("Create in SESSION: RecordId: " + recordid + " | FoodId: " + fmt.Sprint(foodid) + " | UnitId: " + fmt.Sprint(unitid) + " | Quantity: " + fmt.Sprint(quantity) + " | CHO: " + fmt.Sprint(CHO))
		session, _ := store.Get(request, "mysession")
		sessionItem := session.Values["myitems"]
		newItem := Item{
			FoodId:     foodid,
			FoodName:   foodName,
			UnitSymbol: unitSymbol,
			UnitId:     unitid,
			Quantity:   quantity,
			CHO:        CHO,
			RecordId:   recordid,
		}
		myItems := []Item{}
		if sessionItem == nil {
			newItem.Id = "0"
			myItems = append(myItems, newItem)
		} else {
			strItems := session.Values["myitems"].(string)
			json.Unmarshal([]byte(strItems), &myItems)
			newItem.Id = tmp + strconv.Itoa(len(myItems))
			myItems = append(myItems, newItem)
		}
		for index := range myItems {
			item := myItems[index]
			log.Println("Id: " + item.Id)
			log.Println("FoodId: " + strconv.Itoa(item.FoodId))
			log.Println("FoodName: " + item.FoodName)
			log.Println("UnitId: " + strconv.Itoa(item.UnitId))
			log.Println("UnitSymbol: " + item.UnitSymbol)
			log.Println("Quantity: " + fmt.Sprintf("%f", item.Quantity))
			log.Println("CHO: " + fmt.Sprintf("%f", item.CHO))
			log.Println("RecordId: " + item.RecordId)
		}
		bytesItems, _ := json.Marshal(myItems)
		session.Values["myitems"] = string(bytesItems)
		sessions.Save(request, rw)
		tmpl.ExecuteTemplate(rw, "CloseWindow", nil)
	}
}

func InitItemsTable(db *sql.DB) {
	log.Println("Init Items")
	if _, err := db.Exec(
		" CREATE TABLE IF NOT EXISTS Items ( " +
			" id SERIAL PRIMARY KEY, " +
			" record_id integer references Records, " +
			" food_id integer references Foods, " +
			" unit_id integer references Units, " +
			" quantity DOUBLE PRECISION NOT NULL, " +
			" CHO DOUBLE PRECISION NOT NULL " +
			" )"); err != nil {
		log.Fatalf("Error creating database: %q", err)
		return
	}
}

func ListItems(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("List Items")
	selDB, err := db.Query("SELECT " +
		" A.id, B.id, B.name AS food_name, C.id, C.symbol AS unit_symbol, A.quantity, A.CHO " +
		" FROM Items A left join Foods B " +
		" on A.food_id = B.id left join Units C on A.unit_id = C.id ORDER BY a.id DESC")
	if err != nil {
		panic(err.Error())
	}
	item := Item{}
	res := []Item{}
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
		res = append(res, item)
	}
	tmpl.ExecuteTemplate(w, "ListItems", res)
	defer db.Close()
}

func ShowItem(w http.ResponseWriter, r *http.Request) {
	log.Println("*** Show Item ***")
	nId := r.URL.Query().Get("id")
	session, _ := store.Get(r, "mysession")
	item := Item{}
	encontrado := false
	if session != nil {
		log.Println("session")
		if session.Values["myitems"] != nil {
			log.Println("myitems")
			strItems := session.Values["myitems"].(string)
			log.Println("strItems: " + strItems)
			myItems := []Item{}
			json.Unmarshal([]byte(strItems), &myItems)
			for i := range myItems {
				if myItems[i].Id == nId {
					encontrado = true
					item = myItems[i]
					log.Println(myItems[i])
					break
				}
			}
		}
	}
	if !encontrado {
		db := dbConn()
		log.Println("Show Item")
		sqlStatement := "SELECT " +
			" A.id, B.id, B.name AS food_name, C.id, C.symbol AS unit_symbol, A.quantity, A.CHO " +
			" FROM Items A left join Foods B " +
			" on A.food_id = B.id left join Units C on A.unit_id = C.id WHERE a.id = $1"
		log.Println(sqlStatement)
		selDB, err := db.Query(sqlStatement, nId)
		if err != nil {
			panic(err.Error())
		}

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
		}
		defer db.Close()
	}
	tmpl.ExecuteTemplate(w, "ShowItem", item)

}

func EditItem(w http.ResponseWriter, r *http.Request) {
	log.Println("*** Edit Item ***")
	db := dbConn()
	nId := r.URL.Query().Get("id")
	session, _ := store.Get(r, "mysession")
	item := Item{} // ao abrir, a função é EditRecord
	if session != nil {
		log.Println("session")
		if session.Values["myitems"] != nil {
			log.Println("myitems")
			strItems := session.Values["myitems"].(string)
			log.Println("strItems: " + strItems)
			myItems := []Item{}
			json.Unmarshal([]byte(strItems), &myItems)
			for i := range myItems {
				if myItems[i].Id == nId {
					item = myItems[i]
					log.Println(myItems[i])
					break
				}
			}
		}
	} else {
		selDB, err := db.Query("SELECT "+
			" A.id, B.id, C.id, A.quantity, A.CHO, A.record_id "+
			" FROM Items A left join Foods B "+
			" on A.food_id = B.id left join Units C on A.unit_id = C.id WHERE A.id=$1", nId)
		if err != nil {
			panic(err.Error())
		}
		item := Item{}
		for selDB.Next() {
			var id, foodid, unitid int
			var recordid, foodName, unitSymbol string
			var quantity, CHO float64
			err = selDB.Scan(&id, &foodid, &unitid, &quantity, &CHO, &recordid)
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
			item.RecordId = recordid
		}
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
		if item.FoodId == id {
			food.Selected = true
		} else {
			food.Selected = false
		}
		foods = append(foods, food)
	}
	item.FoodOptions = foods
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
		if item.UnitId == id {
			unit.Selected = true
		} else {
			unit.Selected = false
		}
		units = append(units, unit)
	}
	item.UnitOptions = units
	tmpl.ExecuteTemplate(w, "EditItem", item)
	defer db.Close()
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	log.Println("*** Update Item ***")
	if r.Method == "POST" {
		foodid := r.FormValue("foodid")
		unitid := r.FormValue("unit")
		foodName := r.FormValue("foodName")
		unitSymbol := r.FormValue("unitSymbol")
		quantity := r.FormValue("quantity")
		CHO := r.FormValue("CHO")
		id := r.FormValue("uid")
		session, _ := store.Get(r, "mysession")
		encontrado := false
		updatedItems := []Item{}
		editedItem := Item{}
		editedItem.Id = id
		editedItem.FoodId, _ = strconv.Atoi(foodid)
		editedItem.UnitId, _ = strconv.Atoi(unitid)
		editedItem.Quantity, _ = strconv.ParseFloat(quantity, 64)
		editedItem.CHO, _ = strconv.ParseFloat(CHO, 64)
		editedItem.FoodName = foodName
		editedItem.UnitSymbol = unitSymbol
		if session != nil {
			log.Println("session")
			if session.Values["myitems"] != nil {
				log.Println("myitems")
				strItems := session.Values["myitems"].(string)
				log.Println("strItems: " + strItems)
				myItems := []Item{}
				json.Unmarshal([]byte(strItems), &myItems)
				for index := range myItems {
					myItem := myItems[index]
					if myItem.Id != id {
						updatedItems = append(updatedItems, myItem)
					} else {
						encontrado = true
						updatedItems = append(updatedItems, editedItem)
					}
				}
				bytesItems, _ := json.Marshal(updatedItems)
				session.Values["myitems"] = string(bytesItems)
				sessions.Save(r, w)
			}
		}
		if !encontrado { // OLHA SO, AQUI NAO ERA PARA A GENTE DAR UPDATE NO BANCO. SÓ O RECORD QUE TEM ESSE PODER. VOU MIJAR
			db := dbConn()
			sqlStatement := "UPDATE Items SET food_id=$1, unit_id=$2, quantity=$3, cho=$4 WHERE id=$5"
			updtForm, err := db.Prepare(sqlStatement)
			if err != nil {
				panic(err.Error())
			}
			updtForm.Exec(foodid, unitid, quantity, CHO, id)
			log.Println("UPDATE: Id: " + id + " | FoodId: " + foodid + " | UnitId: " + unitid + " | Quantity: " + quantity + " | CHO: " + CHO)
			defer db.Close()
		}
	}
	tmpl.ExecuteTemplate(w, "CloseWindow", nil)
}

func RemoveItem(w http.ResponseWriter, r *http.Request) {
	log.Println("*** Remove Item ***")
	id := r.URL.Query().Get("id")
	log.Println("Delete id: " + id)
	myItems := []Item{}
	updatedItems := []Item{}
	session, _ := store.Get(r, "mysession")
	sessionItem := session.Values["myitems"]
	if sessionItem != nil {
		strItems := session.Values["myitems"].(string)
		json.Unmarshal([]byte(strItems), &myItems)
		for index := range myItems {
			myItem := myItems[index]
			if myItem.Id != id {
				updatedItems = append(updatedItems, myItem)
			}
		}
		bytesItems, _ := json.Marshal(updatedItems)
		session.Values["myitems"] = string(bytesItems)
		sessions.Save(r, w)
	}
	if !strings.HasPrefix(id, "tmp") {
		db := dbConn()
		delForm, err := db.Prepare("DELETE FROM Items WHERE id=$1")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(id)
		log.Println("DELETE: Id: " + id)
		defer db.Close()
	}
	http.Redirect(w, r, "ReloadWindow", 200)
}

func NewItem(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	log.Println("New Item")
	recordId := r.URL.Query().Get("recordid")
	log.Println("recordid: " + recordId)
	item := Item{}
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
	item.FoodOptions = foods
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
	item.UnitOptions = units
	item.RecordId = recordId
	tmpl.ExecuteTemplate(w, "NewItem", item)
}
