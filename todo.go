package main

import (
	"html"
	"html/template"
	"net/http"
	"strings"
)

// セッションIDをキーとして、ToDoリストを保持する
var todoLists = make(map[string][]string)

// セッションIDに紐づくToDoリストを取得する
func getTodoList(sessionId string) []string {
	todos, ok := todoLists[sessionId]
	if !ok {
		todos = []string{}
		todoLists[sessionId] = todos
	}
	return todos
}

// ToDoリストを返す。
func handleTodo(w http.ResponseWriter, r *http.Request) {
	sessionId, err := ensureSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	todos := getTodoList(sessionId)

	t, _ := template.ParseFiles("templates/todo.html")
	t.Execute(w, todos)
}

// セッション上のToDoリストにTodoを追加する
func handleAdd(w http.ResponseWriter, r *http.Request) {
	sessionId, err := ensureSession(w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	todos := getTodoList(sessionId)

	r.ParseForm()
	todo := strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))
	if todo != "" {
		todoLists[sessionId] = append(todos, todo)
	}
	http.Redirect(w, r, "/todo", 303)
}
