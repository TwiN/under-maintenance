package main

import (
	"net/http"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Under maintenance")
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
