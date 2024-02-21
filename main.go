package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	prefix string = "/api/v1" // API prefix
	db     *sql.DB
)

func main() {

	var err error
	// db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/savannahdb?sslmode=disable")
	db, err = sql.Open("postgres", "dbname=savannahdb user=postgres password=Test12345 host=localhost sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/customers", getCustomers)
	router.POST("/customers", createCustomer)
	router.POST("/register", registerCustomer)
	router.POST("/login", loginCustomer)
	router.GET("/orders", getOrders)
	router.POST("/orders", createOrder)
	router.GET("/items", getItems)
	router.POST("/items", createItems)

	router.Run("localhost:8088")
}

func getItems(c *gin.Context) {

}

func createItems(c *gin.Context) {

}

func createOrder(c *gin.Context) {

}

func getOrders(c *gin.Context) {
	// c.Header("Content-Type", "application/json")

	// rows, err := db.Query("SELECT id, customer_id, item, description, price, created_at FROM orders")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// var orders []order
	// for rows.Next() {
	// 	var a order
	// 	err := rows.Scan(&a.ID, &a.Customer_Id, &a.Item, &a.Description, &a.Price, &a.Created_At)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	orders = append(orders, a)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// c.IndentedJSON(http.StatusOK, orders)
}

// returns a list of customers from the database
func getCustomers(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, firstname, lastname, phone, email, created_at FROM customers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var customers []customerViewModel
	for rows.Next() {
		var a customerViewModel
		err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Phone, &a.Email, &a.Created_At)
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
	ID         string `json:"id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Created_At string `json:"created_at"`
	Last_Login string `json:"last_login"`
}

type order struct {
	ID          string `json:"id"`
	Customer_Id string `json:"customer_id"`
	Items       []item `json:"Items"`
	Created_At  string `json:"created_at"`
}

type item struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Created_At  string `json:"created_at"`
}

type itemViewModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
}

type customerReadModel struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type customerViewModel struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Created_At string `json:"created_at"`
	Last_Login string `json:"last_login"`
}

type login struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func loginCustomer(c *gin.Context) {
	var loginCustomer login
	var password string
	// var validation string
	if err := c.BindJSON(&loginCustomer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	//get user by phone
	row := db.QueryRow("SELECT phone, password FROM customers WHERE phone = $1", loginCustomer.Phone)
	err := row.Scan(&loginCustomer.Phone, &password)
	switch err {
	case sql.ErrNoRows:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User Not Found"})
		return
	case nil:
		// check password
		log.Println(password)
		log.Println(loginCustomer.Password)
		validation := bcrypt.CompareHashAndPassword([]byte(password), []byte(loginCustomer.Password))
		if validation == nil {
			c.JSON(http.StatusOK, "Login Ok")
			return
		} else {
			c.JSON(http.StatusForbidden, "Password do not match")
			return
		}

	default:
		log.Fatal(err)
	}
}

func registerCustomer(c *gin.Context) {
	var awesomeCustomer customerReadModel
	if err := c.BindJSON(&awesomeCustomer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	//hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(awesomeCustomer.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO customers (firstname, lastname, phone, password, email) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(awesomeCustomer.FirstName, awesomeCustomer.LastName, awesomeCustomer.Phone, passwordHash, awesomeCustomer.Email); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, awesomeCustomer)
}

func createCustomer(c *gin.Context) {

	var awesomeCustomer customerReadModel
	if err := c.BindJSON(&awesomeCustomer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO customers (firstname, lastname, phone, password, email) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(awesomeCustomer.FirstName, awesomeCustomer.LastName, awesomeCustomer.Phone, awesomeCustomer.Password, awesomeCustomer.Email); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, awesomeCustomer)
}
