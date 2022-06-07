package middleware

import (
	"Stocks_API/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Response struct {
	ID      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "stocksdb"
)

func createConnection() *sql.DB {
	/*err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file, please fix it")
	}*/

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	//db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalf("Error in opening database %v", err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres")

	return db

}

//-----------------------------------------------------------------------------------------------------------------------------------

//we need pointer to request because request is something we will receive from postman or user
func GetAllStock(w http.ResponseWriter, r *http.Request) {

	//we dont have to send anything to this function
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Error while getting all the stocks %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func GetStock(w http.ResponseWriter, r *http.Request) {

	//we will be provided with a parameter value and we need to get that
	params := mux.Vars(r)

	//we have id in string form, lets convert it to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal("Unable to convert string to int")
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Error while getting stocks. %v", err)
	}

	//our response, stock, is encoded in w
	json.NewEncoder(w).Encode(stock)
}

func CreateStock(w http.ResponseWriter, r *http.Request) {

	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatal("Unable to decode the request body")
	}

	insertID := insertStock(stock)
	insertid := int(insertID)
	msg := Response{
		ID:      insertid,
		Message: "Stock created successfully",
	}

	//here we are passing w with a value msg, by encoding it to json, because Golang does not understand json
	json.NewEncoder(w).Encode(msg)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {

	//we will be provided with ID
	params := mux.Vars(r)

	//convert gotten id to number
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error while converting to int %v", err)
	}

	deleteRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock data deleted. Total rows affected %v", deleteRows)

	res := Response{
		ID:      id,
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func UpdateStock(w http.ResponseWriter, r *http.Request) {

	//here we will send the id for which stocks value to update
	//so firstly catch the ID
	params := mux.Vars(r)

	//params is in string format
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error while converting id to int %v", err)
	}

	//now we need something to hold all the stocks data
	var stock models.Stock

	//also, we have entered new data to be updated in r, and we need to decode that to use that
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Error while decoding the information from user to stocks %v", err)
	}

	//now in stock we have the decoded information of what we want to send in the database, let's send that data in database
	updateRows := updateStock(int64(id), stock)

	m := fmt.Sprintf("Stocks updated successfully, %v rows affected", updateRows)

	msg := Response{
		ID:      id,
		Message: m,
	}

	json.NewEncoder(w).Encode(msg)

}

//-----------------------------------------------------------------------------------------------------------------------------

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES($1, $2, $3) RETURNING stockid`

	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to insert the record from final function %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id

}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid = $1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Printf("No rows were returned ")
		return stock, nil

	case nil:
		return stock, nil

	default:
		log.Fatalf("Error while returning specific stock %v", err)
	}
	return stock, err

}

func getAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock

	sqlStatement := `SELECT * FROM stocks`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Got error while selecting all stocks from database %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		//we will be receiving all the nodes, appending them and returning them
		err = rows.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan rows %v", err)
		}

		stocks = append(stocks, stock)

	}

	return stocks, err

}

func updateStock(id int64, stock models.Stock) int64 {

	db := createConnection()
	defer db.Close()

	//this is going to be a little diferent, here we will update the values

	sqlStatement := `UPDATE stocks SET name = $2, price = $3, company = $4 where stockid = $1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Error while updating value %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while getting rows affected")
	}

	fmt.Printf("Total rows affected = %v", rowsAffected)

	return rowsAffected
}

func deleteStock(id int64) int64 {

	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid = $1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Error while deleting the record %v", err)
	}

	rows, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while getting rows affected %v", err)
	}

	fmt.Printf("Rows affected while deleting = %v", rows)

	return rows

}

//{
//    "name":"Tsl",
//    "price":125,
//    "company":"Tesla"
//}
