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
var tasksPage TasksPage

type TasksPage struct {
	Title    string
	AllTasks []sql_db.Task
}

func getAllTodos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tasksPage.AllTasks = sql_db.GetAllTasks(SqlHandlerDB.Db)

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := TasksPage{Title: "Displaying All Tasks:", AllTasks: tasksPage.AllTasks}
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
			tasksPage.AllTasks = sql_db.GetTaskbyTitle(SqlHandlerDB.Db, InputUrlTask)
			p := TasksPage{Title: "Displaying Tasks by Title:", AllTasks: tasksPage.AllTasks}
			t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
			err := t.Execute(w, p)
			if err != nil {
				panic(err)
			}
			return
		}
		// If urlIndex was successfully changed by the int conversion of the second segment:
		// sql query by int Id
		if InputUrlTask.Id != urlIndex {
			InputUrlTask.Id = urlIndex
			tasksPage.AllTasks = sql_db.GetTaskbyId(SqlHandlerDB.Db, InputUrlTask)
			p := TasksPage{Title: "Displaying Tasks by ID:", AllTasks: tasksPage.AllTasks}
			t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
			err := t.Execute(w, p)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	getAllTodos(w, r)
}

func completeTodosById(w http.ResponseWriter, r *http.Request) {
	sql_db.CompleteTask(SqlHandlerDB.Db, InputUrlTask)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	InputUrlTask.Title = r.FormValue("inputTitle")
	inputFormValue := r.FormValue("inputComplete")
	if inputFormValue != "" {
		boolFromStr, err := strconv.ParseBool(inputFormValue)
		if err != nil {
			panic(err.Error())
		}
		InputUrlTask.Completed = boolFromStr
	}
	sql_db.AddTask(SqlHandlerDB.Db, InputUrlTask)
	getAllTodos(w, r)
}
