package myhandlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/martoranam/hello_world/helloworlder"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

var routes = []route{
	newRoute("GET", "/", home),
	newRoute("GET", "/contact", contact),
	newRoute("GET", "/helloworld", helloworld),
	newRoute("GET", "/todos", getAllTodos),
	newRoute("GET", "/todos/([0-9]+)", getTodosById),
	//newRoute("GET", "/todos/:title", getTodosByTitle),
	//newRoute("PATCH", "/todos/:id", completeTodosById),
	newRoute("POST", "/todos", addTodo),
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type ctxKey struct{}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type homePage struct {
	Title   string
	HomeMsg string
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := homePage{Title: "Welcome to my Website!", HomeMsg: "This is the HomePage"}
	t, _ := template.ParseFiles("html/hometemplate.html")
	fmt.Println(r)
	fmt.Println(t.Execute(w, p))
}

type contactPage struct {
	Title    string
	Contacts string
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := contactPage{Title: "Contact me via: ", Contacts: "mark_tech@hotmail.com"}
	t, _ := template.ParseFiles("html/contactstemplate.html")
	fmt.Println(r)
	fmt.Println(t.Execute(w, p))
}

type helloworldPage struct {
	Title   string
	HWprint string
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	var printable string
	helloworlder.Update(&printable)

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p := helloworldPage{Title: "This page should print a message below!", HWprint: printable}
	t, _ := template.ParseFiles("html/helloworldtemplate.html")
	fmt.Println(r)
	fmt.Println(t.Execute(w, p))
}
