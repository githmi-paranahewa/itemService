package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// "golang.org/x/tools/go/analysis/passes/appends"
)

type Item struct {
	ID         string
	Name       string
	Price      float64
	Quantity   int
	OrderItems []OrderItem `gorm:"foreignKey:ItemID"`
}

type OrderItem struct {
	ItemID   string
	Quantity int
	OrderID  string `gorm:"foreignKey:OrderID"`
}

type Order struct {
	Items  []OrderItem
	Total  float64
	Status string
}

var items []Item

func GetItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	json.NewEncoder(w).Encode(items)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for index, instance := range items {
		if instance.ID == params["itemId"] {
			items = append(items[:index], items[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(items)
}

func GetItemById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for _, instance := range items {
		if instance.ID == params["itemId"] {
			json.NewEncoder(w).Encode(instance)
			return
		}
	}
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	var item Item
	_ = json.NewDecoder(r.Body).Decode(&item)
	item.ID = strconv.Itoa(rand.Intn(100000000))
	items = append(items, item)
	json.NewEncoder(w).Encode(item)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for index, instance := range items {
		if instance.ID == params["itemId"] {
			items = append(items[:index], items[index+1:]...)
			var item Item
			_ = json.NewDecoder(r.Body).Decode(&item)
			item.ID = params["itemId"]
			items = append(items, item)
			json.NewEncoder(w).Encode(item)
		}
	}

}

func main() {
	r := mux.NewRouter()

	items = append(items, Item{ID: "1", Name: "Book", Price: 300, Quantity: 10})
	items = append(items, Item{ID: "2", Name: "Pen", Price: 40, Quantity: 20})
	r.HandleFunc("/item", AddItem).Methods("POST")
	r.HandleFunc("/item", GetItem).Methods("GET")
	r.HandleFunc("/item/{itemId}", GetItemById).Methods("GET")
	r.HandleFunc("/item/{itemId}", UpdateItem).Methods("PUT")
	r.HandleFunc("/item/{itemId}", DeleteItem).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:9010", r))
}
