package main

import (
	"context"
	"encoding/json"
	"os"

	// "fmt"
	// "log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2/clientcredentials"
	// "golang.org/x/tools/go/analysis/passes/appends"
)

type Item struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
	// OrderItems []OrderItem `gorm:"foreignKey:ItemID"`
}

// type OrderItem struct {
// 	ItemID   string
// 	Quantity int
// 	OrderID  string `gorm:"foreignKey:OrderID"`
// }

// type Order struct {
// 	Items  []OrderItem
// 	Total  float64
// 	Status string
// }

var clientCredsConfig = clientcredentials.Config{
	ClientID:     "CLIENT_ID",
	ClientSecret: "CLIENT_SECRET",
	TokenURL:     "TOKEN_URL",
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
	client := clientCredsConfig.Client(context.Background())
	serviceURL := os.Getenv("SERVICE_URL")
	// clientID := os.Getenv("CONSUMER_KEY")
	// clientSecret := os.Getenv("CONSUMER_SECRET")
	// tokenURL := os.Getenv("TOKEN_URL")

	items = append(items, Item{ID: "1", Name: "Book", Price: 300, Quantity: 10})
	items = append(items, Item{ID: "2", Name: "Pen", Price: 40, Quantity: 20})
	r.HandleFunc("/item", AddItem).Methods("POST")
	r.HandleFunc("/item", GetItem).Methods("GET")
	r.HandleFunc("/item/{itemId}", GetItemById).Methods("GET")
	r.HandleFunc("/item/{itemId}", UpdateItem).Methods("PUT")
	r.HandleFunc("/item/{itemId}", DeleteItem).Methods("DELETE")
	http.Handle("/", authenticateMiddleware(client, serviceURL)(r))

	http.ListenAndServe(":9010", nil)

	// log.Fatal(http.ListenAndServe(":9010", r))
}
func authenticateMiddleware(client *http.Client, serviceURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Authenticate the request using the provided client
			_, err := client.Get(serviceURL)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
