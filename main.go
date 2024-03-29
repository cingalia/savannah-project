package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	db *sql.DB
)

const (
	username = "sandbox"                                                          //Your Africa's Talking Username
	apiKey   = "86208402466e939a2ed4c971b20dd84a7ad7674237fead0c6e8ba5e4f82e7152" //Production or Sandbox API Key
	env      = "Sandbox"                                                          // Choose either Sandbox or Production
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
	router.POST("/items", createItem)

	router.Run("localhost:8088")
}

func getItems(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name, description, price, created_at FROM items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []item
	for rows.Next() {
		var a item
		err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Price, &a.Created_At)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, items)
}

func createItem(c *gin.Context) {
	var newItem itemViewModel
	if err := c.BindJSON(&newItem); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO items (name, description, price) VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(newItem.Name, newItem.Description, newItem.Price); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, newItem)
}

func createOrder(c *gin.Context) {
	var newOrder orderViewModel
	var summarybuilder strings.Builder
	summarybuilder.WriteString("You order:")

	if err := c.BindJSON(&newOrder); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// get items ids
	itemIds := strings.Split(newOrder.Items_Ids, ",")
	//get items information from database
	for i := 0; i < len(itemIds); i++ {
		var currentItem item
		row := db.QueryRow("SELECT id, name, description, price, created_at FROM items WHERE id = $1", itemIds[i])
		err := row.Scan(&currentItem.ID, &currentItem.Name, &currentItem.Description, &currentItem.Price, &currentItem.Created_At)
		switch err {
		case sql.ErrNoRows:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Item not found"})
			return
		case nil:
			// current item information
			log.Println(currentItem)
			summarybuilder.WriteString(" " + strconv.Itoa(i+1) + " - " + currentItem.Name + " Price: Kshs " + strconv.Itoa(currentItem.Price))
			newOrder.Total_Price = newOrder.Total_Price + currentItem.Price
		default:
			log.Fatal(err)
		}
	}
	newOrder.Summary = summarybuilder.String()
	//Insert into db
	stmt, err := db.Prepare("INSERT INTO orders (customer_id, items_ids, summary, total_price) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(newOrder.Customer_Id, newOrder.Items_Ids, newOrder.Summary, newOrder.Total_Price); err != nil {
		log.Fatal(err)
	}

	//get customers number
	var phone string
	row := db.QueryRow("SELECT phone FROM customers WHERE id = $1", newOrder.Customer_Id)
	err = row.Scan(&phone)
	switch err {
	case sql.ErrNoRows:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "phone not found"})
		return
	case nil:
		// current item information
		log.Println(phone)
		//Send SMS to customer
		client := &http.Client{}
		var data = strings.NewReader(`username=sandbox&to=%2B` + phone + `&message=Hello%20World!&from=23370`)
		req, err := http.NewRequest("POST", "https://api.sandbox.africastalking.com/version1/messaging", data)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("apiKey", "86208402466e939a2ed4c971b20dd84a7ad7674237fead0c6e8ba5e4f82e7152")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", bodyText)
	default:
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, newOrder)
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
	Customer_Id int    `json:"customer_id"`
	Items_Ids   string `json:"Items_ids"`
	Summary     string `json:"summary"`
	Total_Price int    `json:"total_price"`
	Created_At  string `json:"created_at"`
}

type orderViewModel struct {
	Customer_Id int    `json:"customer_id"`
	Items_Ids   string `json:"Items_ids"`
	Summary     string `json:"summary"`
	Total_Price int    `json:"total_price"`
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
