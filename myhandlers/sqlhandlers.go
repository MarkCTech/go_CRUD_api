package myhandlers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/martoranam/sql_db"
	"github.com/rs/xid"
)

var (
	Database *sql_db.Database
	list     []sql_db.Task
	mtx      sync.RWMutex
	once     sync.Once
)

type todosPage struct {
	Title    string
	AllTasks []sql_db.Task
}

func init() {
	once.Do(initialiseList)
}

func initialiseList() {
	list = []sql_db.Task{}
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	list = sql_db.GetAllTasks(Database.Db)

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := todosPage{Title: "Displaying All Tasks:", AllTasks: list}
	t, _ := template.ParseFiles("html/alltaskstemplate.html")
	err := t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}

func GetTodobyId(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	segments := strings.Split(path, "/")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	defer r.Body.Close()

	urlIntIndex, err := strconv.ParseInt(segments[2], 10, 8)
	switch {
	case err == nil:
		list = sql_db.GetTaskbyId(Database.Db, segments[2])
	case urlIntIndex == 0:
		list = sql_db.GetTaskbyTitle(Database.Db, segments[2])
	default:
		GetAllTodos(w, r)
	}
	p := todosPage{Title: "Displaying Tasks by URL index:", AllTasks: list}
	t, _ := template.ParseFiles("html/tasksbyurltemplate.html")
	err = t.Execute(w, p)
	if err != nil {
		panic(err.Error())
	}
}

// func convertHTTPBodyToTodo(httpBody io.ReadCloser) (sql_db.Task, int, error) {
// 	body, err := io.ReadAll(httpBody)
// 	if err != nil {
// 		return sql_db.Task{}, http.StatusInternalServerError, err
// 	}
// 	defer httpBody.Close()
// 	return convertJSONBodyToTodo(body)
// }

// func convertJSONBodyToTodo(jsonBody []byte) (sql_db.Task, int, error) {
// 	var todoItem sql_db.Task
// 	err := json.Unmarshal(jsonBody, &todoItem)
// 	if err != nil {
// 		return sql_db.Task{}, http.StatusBadRequest, err
// 	}
// 	return todoItem, http.StatusOK, nil
// }

func newTodo(title string) sql_db.Task {
	return sql_db.Task{
		Id:        xid.New().String(),
		Title:     title,
		Completed: false,
	}
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("inputTitle")
	parsedTodo := newTodo(title)
	parsedComplete := r.FormValue("inputComplete")
	if parsedComplete != "" {
		boolFromStr, err := strconv.ParseBool(parsedComplete)
		if err != nil {
			panic(err.Error)
		}
		parsedTodo.Completed = boolFromStr
	}
	mtx.Lock()
	list = append(list, parsedTodo)
	mtx.Unlock()
	sql_db.AddTask(Database.Db, &parsedTodo)
	GetAllTodos(w, r)
}

func CompletebyId(w http.ResponseWriter, r *http.Request) {
	statusId := r.FormValue("inputTitle")
	statusTodo := sql_db.Task{Id: statusId}
	intId, err := strconv.Atoi(statusId)
	if err != nil {
		panic(err.Error())
	}
	// Updates local list and database
	setTodoCompleteByLocation(intId)
	sql_db.CompleteTask(Database.Db, &statusTodo)
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, t := range list {
		if isMatchingID(t.Id, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find todo based on id")
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Completed = true
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}

// Delete will remove a Todo from the Todo list
func DeletebyId(id string) error {
	location, err := findTodoLocation(id)
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}
