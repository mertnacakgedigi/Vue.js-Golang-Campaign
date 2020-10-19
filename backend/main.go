package main

import (
	"fmt"
	"backend/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("Starting server on the port 8000...")
	
}