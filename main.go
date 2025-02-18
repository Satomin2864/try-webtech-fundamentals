package main

import (
	"html/template"
	"log"
	"net/http"
)

// TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.
var todoList []string

func handleTodo(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/todo.html")
	t.Execute(w, todoList)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	todo := r.Form.Get("todo")
	todoList = append(todoList, todo)
	http.Redirect(w, r, "/todo", 303)
}

func main() {
	todoList = append(todoList, "顔を洗う", "朝食を食べる", "歯を磨く")
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/todo", handleTodo)
	http.HandleFunc("/add", handleAdd)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}
