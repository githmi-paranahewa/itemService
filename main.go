package main

import (
	"encoding/json"
	"log"

	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Item struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}

type Order struct {
	ID    string
	Total int
}

var items []Item
var orders []Order

func GetItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	json.NewEncoder(w).Encode(items)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	json.NewEncoder(w).Encode(orders)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "pkgication/json")
	params := mux.Vars(r)
	for index, instance := range items {
		if instance.ID == params["itemId"] {
			items = append(items[:index], items[index+1:]...)
			json.NewEncoder(w).Encode("Successfully Deleted the Item")
			return
		}
	}

	http.NotFound(w, r)
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
	http.NotFound(w, r)
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
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, instance := range items {
		if instance.ID == params["itemId"] {
			var updatedItem Item
			err := json.NewDecoder(r.Body).Decode(&updatedItem)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedItem.ID = params["itemId"]
			items[index] = updatedItem

			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	r := mux.NewRouter()
	r2 := mux.NewRouter()

	items = append(items, Item{ID: "1", Name: "Book", Price: 300, Quantity: 10})
	items = append(items, Item{ID: "2", Name: "Pen", Price: 40, Quantity: 20})

	orders = append(orders, Order{ID: "1", Total: 1500})

	r.HandleFunc("/item", AddItem).Methods("POST")
	r.HandleFunc("/item", GetItem).Methods("GET")
	r.HandleFunc("/item/{itemId}", GetItemById).Methods("GET")
	r.HandleFunc("/item/{itemId}", UpdateItem).Methods("PUT")
	r.HandleFunc("/item/{itemId}", DeleteItem).Methods("DELETE")

	r2.HandleFunc("/order", GetOrder).Methods("GET")
	// log.Fatal(http.ListenAndServe(":9010", r))

	// 2nd endpoint

	// log.Fatal(http.ListenAndServe(":9090", r2))
	go func() {
		log.Fatal(http.ListenAndServe(":9010", r))
	}()

	// Start HTTP server for the second endpoint
	go func() {
		log.Fatal(http.ListenAndServe(":9090", r2))
	}()

	// Block main goroutine to keep servers running
	select {}
}
