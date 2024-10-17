package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Employee struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	TotalPosts int    `json:"total_posts"`
}

func initDB() {
	var err error
	db, err = sql.Open("postgres", "postgres://fiapi:fiapi@localhost/testing")
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(15)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("could not connect to database: %v\n", err)
	}
	fmt.Println("connected to the database successfully!")
}

func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(db.Stats()); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func decorateEmployee(e *Employee) error {
	randomTime := time.Duration(rand.Intn(50)) * time.Millisecond
	time.Sleep(randomTime)
	e.TotalPosts = rand.Intn(1000)
	return nil
}

func GetEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT
            id,
            name,
            email
        FROM employee
    `)

	if err != nil {
		http.Error(w, fmt.Sprintf("error querying database: %v", err), http.StatusInternalServerError)
		return
	}

	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Email); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("error during row iteration: %v", err), http.StatusInternalServerError)
		return
	}
	rows.Close()

	for i := range employees {
		e := employees[i]

		decorateEmployee(&e)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		log.Printf("Error encoding response: %v", err) // Log the error instead of using fmt.Sprintf
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func main() {
	initDB()

	http.HandleFunc("/info", GetInfoHandler)
	http.HandleFunc("/employees", GetEmployeesHandler)

	port := ":8080"
	fmt.Printf("starting server on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server failed to start: %v\n", err)
	}
}
