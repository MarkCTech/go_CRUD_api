package myhandlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/martoranam/sql_db"
)

var SqlHandlerDB *sql_db.Database
var InputTask *sql_db.Task
var ReturnedTask *sql_db.Task

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
	sql_db.GetTaskbyTitle(SqlHandlerDB.Db, InputTask)
}

func completetodosbyid(w http.ResponseWriter, r *http.Request) {
	sql_db.CompleteTask(SqlHandlerDB.Db, InputTask)
}

func addtodo(w http.ResponseWriter, r *http.Request) {
	sql_db.AddTask(SqlHandlerDB.Db, InputTask)
}
