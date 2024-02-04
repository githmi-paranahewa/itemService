package main

import (
	"encoding/json"
	"fmt"
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

// var clientCredsConfig = clientcredentials.Config{

// 	ClientID:     "CLIENT_ID",
// 	ClientSecret: "CLIENT_SECRET",
// 	TokenURL:     "TOKEN_URL",
// }

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
	fmt.Println("hi")
	for index, instance := range items {
		fmt.Println("hi index", index)
		fmt.Println("hi outside of ", instance.ID, " para ", params["itemId"])
		if instance.ID == params["itemId"] {
			fmt.Println("hi inside of ", instance.ID)
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

	items = append(items, Item{ID: "1", Name: "Book", Price: 300, Quantity: 10})
	items = append(items, Item{ID: "2", Name: "Pen", Price: 40, Quantity: 20})

	r.HandleFunc("/item", AddItem).Methods("POST")
	r.HandleFunc("/item", GetItem).Methods("GET")
	r.HandleFunc("/item/{itemId}", GetItemById).Methods("GET")
	r.HandleFunc("/item/{itemId}", UpdateItem).Methods("PUT")
	r.HandleFunc("/item/{itemId}", DeleteItem).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9010", r))
}
