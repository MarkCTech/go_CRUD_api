package myhandlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/martoranam/sql_db"
)

var SqlHandlerDB *sql_db.Database
var InputUrlTask *sql_db.Task
var ReturnedTask *sql_db.Task

type TasksPage struct {
	Title    string
	AllTasks []sql_db.Task
}

func getAllTodos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	allTasks := sql_db.GetAllTasks(SqlHandlerDB.Db)

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := TasksPage{Title: "Displaying All Tasks:", AllTasks: allTasks}
	t, _ := template.ParseFiles("html/alltaskstemplate.html")
	err := t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}

func getUrlTodos(w http.ResponseWriter, r *http.Request) { //8
	defer r.Body.Close()
	path := r.URL.Path
	segments := strings.Split(path, "/")

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	var allTasks []sql_db.Task

	// If url longer than 2 segments:
	// Check if trailing slash
	if len(segments) > 2 {
		if segments[2] == "" {
			getAllTodos(w, r)
			return
		}
		// If second segment is not assignable to an integer:
		// sql query for []Task by string Title
		urlIndex, err := strconv.Atoi(segments[2])
		if err != nil {
			InputUrlTask.Title = segments[2]
			allTasks = sql_db.GetTaskbyTitle(SqlHandlerDB.Db, InputUrlTask)
			p := TasksPage{Title: "Displaying Tasks by Title:", AllTasks: allTasks}
			t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
			err := t.Execute(w, p)
			if err != nil {
				panic(err)
			}
			return
		}
		// If urlIndex was changed by the int conversion of the second segment:
		// sql query by int Id
		if InputUrlTask.Id != urlIndex {
			InputUrlTask.Id = urlIndex
			allTasks = sql_db.GetTaskbyId(SqlHandlerDB.Db, InputUrlTask)
			p := TasksPage{Title: "Displaying Tasks by ID:", AllTasks: allTasks}
			t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
			err := t.Execute(w, p)
			if err != nil {
				panic(err)
			}
			InputUrlTask.Id = 0
			return
		}
	}
	getAllTodos(w, r)
}

func completeTodosById(w http.ResponseWriter, r *http.Request) {
	sql_db.CompleteTask(SqlHandlerDB.Db, InputUrlTask)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	sql_db.AddTask(SqlHandlerDB.Db, InputUrlTask)
}
