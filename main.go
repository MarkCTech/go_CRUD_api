package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/martoranam/go_site/myhandlers"
)

func main() {
	const port = 5000
	router := http.HandlerFunc(myhandlers.Serve)

	fmt.Printf("listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
