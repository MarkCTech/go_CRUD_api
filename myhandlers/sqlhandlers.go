package myhandlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	// "regexp"
	"strconv"

	"github.com/martoranam/sql_db"
)

var SqlHandlerDB *sql_db.Database
var InputTask *sql_db.Task
var ReturnedTask *sql_db.Task

// var validPath = regexp.MustCompile("^/(todos)/([0-9]+)$")

// func getField(r *http.Request, index int) string {
// 	fields := r.Context().Value(ctxKey{}).([]string)
// 	return fields[index]
// }

type AllTasksPage struct {
	Title    string
	AllTasks []sql_db.Task
}

func getTodos(w http.ResponseWriter, r *http.Request) { //8
	defer r.Body.Close()
	path := r.URL.Path
	segments := strings.Split(path, "/")
	InputTask.Id = 0
	var allTasks []sql_db.Task
	if len(segments) > 2 {
		InputTask.Id, _ = strconv.Atoi(segments[2])
		allTasks = sql_db.GetTaskbyId(SqlHandlerDB.Db, InputTask)
	}
	if InputTask.Id == 0 {
		allTasks = sql_db.GetAllTasks(SqlHandlerDB.Db)
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := AllTasksPage{Title: "Displaying Tasks:", AllTasks: allTasks}
	t, _ := template.ParseFiles("html/alltaskstemplate.html")
	fmt.Println(r)
	fmt.Println(t.Execute(w, p))
}

type IdTasksPage struct {
	AllTasksPage
}

// func getTodosByTitle(w http.ResponseWriter, r *http.Request) {
// 	allTasks := sql_db.GetTaskbyTitle(SqlHandlerDB.Db, InputTask)

// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
// 	p := AllTasksPage{Title: "Displaying tasks by Title:", AllTasks: allTasks}
// 	t, _ := template.ParseFiles("html/taskbytitletemplate.html")
// 	fmt.Println(r)
// 	fmt.Println(t.Execute(w, p))
// }

func completeTodosById(w http.ResponseWriter, r *http.Request) {
	sql_db.CompleteTask(SqlHandlerDB.Db, InputTask)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	sql_db.AddTask(SqlHandlerDB.Db, InputTask)
}
