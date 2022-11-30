package myhandlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/martoranam/sql_db"
)

// type sqlinput struct {
// 	sql_db.Task
// }

// type sqlreturn struct {
// 	sql_db.Task
// }

var SqlHandlerDB *sql_db.Database

type AllTasksPage struct {
	Title    string
	AllTasks []sql_db.Task
}

func getalltodos(w http.ResponseWriter, r *http.Request) {

	alltasks := sql_db.GetAllTasks(SqlHandlerDB.Db)

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := AllTasksPage{Title: "Displaying All Tasks:", AllTasks: alltasks}
	t, _ := template.ParseFiles("html/alltaskstemplate.html")
	fmt.Println(r)
	fmt.Println(t.Execute(w, p))
}

func gettodosbyid(w http.ResponseWriter, r *http.Request) {

}

func completetodosbyid(w http.ResponseWriter, r *http.Request) {

}

func addtodo(w http.ResponseWriter, r *http.Request) {

}
