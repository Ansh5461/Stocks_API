//5:57
package main

import (
	"Stocks_API/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	//we are calling Router function inside router folder, which will start the server
	r := router.Router()
	fmt.Println("Starting server at port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
