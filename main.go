package main

import (
//	"fmt"
//	"io"
	"log"
	"net/http"
	"hqrUtils"
//	"os"
)
func handleRequests() {
	http.HandleFunc("/", hqrUtils.ConvertXml)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {

	handleRequests()
}
