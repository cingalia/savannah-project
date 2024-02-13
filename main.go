package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	prefix string = "/api/v1" // API prefix
	db     *sql.DB
)

func main() {

	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/savannahdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/customers", getCustomers)
	router.POST("/customers", createCustomer)

	router.Run("localhost:8080")
}

func getCustomers(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, firstname, lastname, phone, email FROM customers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var customers []customer
	for rows.Next() {
		var a customer
		err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Phone, &a.Email)
		if err != nil {
			log.Fatal(err)
		}
		customers = append(customers, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, customers)
}

type customer struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}
