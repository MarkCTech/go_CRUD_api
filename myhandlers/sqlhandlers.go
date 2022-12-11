package myhandlers

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/martoranam/sql_db"
	"github.com/rs/xid"
)

var (
	Database *sql_db.Database
	list     []sql_db.Task
	mtx      sync.RWMutex
	once     sync.Once
)

// type todosPage struct {
// 	Title       string
// 	AddNewTitle string
// 	AllTasks    []sql_db.Task
// }

func init() {
	once.Do(initialiseList)
}

func initialiseList() {
	list = []sql_db.Task{}
}

func GetAllTodos(c *gin.Context) {
	list = sql_db.GetAllTasks(Database.Db)

	c.HTML(http.StatusOK, "alltaskstemplate.html", gin.H{
		"Title":        "Displaying All Tasks:",
		"AllTasks":     list,
		"AddNewTitle:": "Add a new Task here: ",
	})
}

func GetTodobyId(c *gin.Context) {
	id := c.Param("id")
	list = sql_db.GetTaskbyId(Database.Db, id)

	c.HTML(http.StatusOK, "tasksbyurltemplate.html", gin.H{
		"Title":        "Displaying Tasks by URL index:",
		"AllTasks":     list,
		"AddNewTitle:": "Add a new Task here: ",
	})
}

func PostNewTodoFromDetail(c *gin.Context) {

}

func newTodo(title string) sql_db.Task {
	return sql_db.Task{
		Id:        xid.New().String(),
		Title:     title,
		Completed: false,
	}
}

func AddTodo(c *gin.Context) {
	parsedTitle := c.PostForm("inputTitle")
	parsedComplete := c.PostForm("inputComplete")

	toBeTodo := newTodo(parsedTitle)
	if parsedComplete != "" {
		boolFromStr, err := strconv.ParseBool(parsedComplete)
		if err != nil {
			panic(err.Error)
		}
		toBeTodo.Completed = boolFromStr
	}
	mtx.Lock()
	list = append(list, toBeTodo)
	mtx.Unlock()
	sql_db.AddTask(Database.Db, &toBeTodo)
	GetAllTodos(c)
}

func CompletebyId(c *gin.Context) {
	parsedId := c.PostForm("inputId")
	parsedComplete := c.PostForm("markComplete")
	intId, err := strconv.Atoi(parsedId)
	if err != nil {
		panic(err.Error())
	}
	statusTodo := new(sql_db.Task)
	if parsedComplete != "" {
		boolFromStr, err := strconv.ParseBool(parsedComplete)
		if err != nil {
			panic(err.Error)
		}
		statusTodo.Id = parsedId
		statusTodo.Completed = boolFromStr
	}
	// Updates local list and database
	setTodoCompleteByLocation(intId)
	sql_db.CompleteTask(Database.Db, statusTodo)
}

// Delete will remove a Todo from the Todo list
func DeletebyId(c *gin.Context) {
	id := c.Param("id")
	deleteTodo := new(sql_db.Task)
	deleteTodo.Id = id

	location, err := findTodoLocation(id)
	if err != nil {
		panic(err.Error())
	}
	sql_db.DeleteTask(Database.Db, deleteTodo)
	removeElementByLocation(location)
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	for i, t := range list {
		if isMatchingID(t.Id, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find LOCAL todo based on id")
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
