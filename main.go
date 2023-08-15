package main

/*
	This is the main entry point of the project. We will call the Router function to get all the api's call


*/

import (
	"fmt"
	"todo-list/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	

	fmt.Println("Starting the  server on the port 3000...")
	
	log.Fatal(http.ListenAndServe(":3000", r))
}
