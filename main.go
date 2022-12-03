package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/martoranam/go_site/myhandlers"
	"github.com/martoranam/sql_db"
)

func main() {
	database := sql_db.Dbstart("todosdb")
	myhandlers.SqlHandlerDB = database
	myhandlers.InputUrlTask = &sql_db.Task{Title: "TESTENTRY"}
	sql_db.AddTask(database.Db, myhandlers.InputUrlTask)
	const port = 5000
	router := http.HandlerFunc(myhandlers.Serve)
	fmt.Printf("listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
