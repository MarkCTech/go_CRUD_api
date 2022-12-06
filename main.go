package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/martoranam/go_site/myhandlers"
	"github.com/martoranam/sql_db"
)

func main() {
	myhandlers.Database = sql_db.Dbstart("todosdb")
	const port = 5000
	router := http.HandlerFunc(myhandlers.Serve)
	fmt.Printf("listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
