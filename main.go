package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Expense struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Date        string  `json:"date"` // Date format can be YYYY-MM-DD
}

var expenses []Expense
var mu sync.Mutex
var nextID = 1

// Create a new expense
func createExpense(w http.ResponseWriter, r *http.Request) {
	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	expense.ID = nextID
	nextID++
	mu.Lock()
	expenses = append(expenses, expense)
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}

// Get all expenses
func getAllExpenses(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// Get an expense by ID
func getExpenseByID(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	id := mux.Vars(r)["id"]
	for _, expense := range expenses {
		if id == strconv.Itoa(expense.ID) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expense)
			return
		}
	}
	http.NotFound(w, r)
}

// Update an expense
func updateExpense(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	id := mux.Vars(r)["id"]
	for i, expense := range expenses {
		if id == strconv.Itoa(expense.ID) {
			if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			expense.ID = expense.ID // Preserve the ID
			expenses[i] = expense
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expense)
			return
		}
	}
	http.NotFound(w, r)
}

// Delete an expense
func deleteExpense(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	id := mux.Vars(r)["id"]
	for i, expense := range expenses {
		if id == strconv.Itoa(expense.ID) {
			expenses = append(expenses[:i], expenses[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/expenses", createExpense).Methods("POST")
	router.HandleFunc("/expenses", getAllExpenses).Methods("GET")
	router.HandleFunc("/expenses/{id}", getExpenseByID).Methods("GET")
	router.HandleFunc("/expenses/{id}", updateExpense).Methods("PUT")
	router.HandleFunc("/expenses/{id}", deleteExpense).Methods("DELETE")

	http.ListenAndServe(":3690", router)
}
