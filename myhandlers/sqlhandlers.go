package myhandlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	// "regexp"

	"github.com/martoranam/sql_db"
)

var SqlHandlerDB *sql_db.Database
var InputUrlTask *sql_db.Task
var ReturnedTask *sql_db.Task

type TasksPage struct {
	Title    string
	AllTasks []sql_db.Task
}

// func getAllTodos(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()

// 	allTasks := sql_db.GetAllTasks(SqlHandlerDB.Db)

// 	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
// 	p := TasksPage{Title: "Displaying All Tasks:", AllTasks: allTasks}
// 	t, _ := template.ParseFiles("html/alltaskstemplate.html")
// 	fmt.Println(r)
// 	fmt.Println(t.Execute(w, p))
// }

func getUrlTodos(w http.ResponseWriter, r *http.Request) { //8
	defer r.Body.Close()
	path := r.URL.Path
	segments := strings.Split(path, "/")

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	var allTasks []sql_db.Task

	// If url longer than 2 segments:
	//
	if len(segments) > 2 {
		urlIndex, err := strconv.Atoi(segments[2])
		if err != nil {
			InputUrlTask.Title = segments[2]
			allTasks = sql_db.GetTaskbyTitle(SqlHandlerDB.Db, InputUrlTask)
			p := TasksPage{Title: "Displaying Tasks by Title:", AllTasks: allTasks}
			t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
			fmt.Println(r)
			fmt.Println(t.Execute(w, p))
			return
		}

		if InputUrlTask.Id != urlIndex {
			InputUrlTask.Id = urlIndex
			allTasks = sql_db.GetTaskbyId(SqlHandlerDB.Db, InputUrlTask)
			p := TasksPage{Title: "Displaying Tasks by ID:", AllTasks: allTasks}
			t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
			fmt.Println(r)
			fmt.Println(t.Execute(w, p))
			InputUrlTask.Id = 0
			return
		}
	}
	allTasks = sql_db.GetAllTasks(SqlHandlerDB.Db)
	p := TasksPage{Title: "Displaying All Tasks:", AllTasks: allTasks}
	t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
	fmt.Println(r)
	fmt.Println(t.Execute(w, p))
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
	sql_db.CompleteTask(SqlHandlerDB.Db, InputUrlTask)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	sql_db.AddTask(SqlHandlerDB.Db, InputUrlTask)
}
