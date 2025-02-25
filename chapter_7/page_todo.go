package main

import (
	"html"
	"net/http"
	"strings"
)

type TodoPageData struct {
	UserId   string
	Expires  string
	ToDoList []*ToDoItem
}

func handleTodo(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	if !isAuthenticated(w, r, session) {
		return
	}

	pageData := TodoPageData{
		UserId:   session.UserAccount.Id,
		Expires:  session.UserAccount.ExpiresText(),
		ToDoList: session.UserAccount.ToDoList.Items,
	}

	templates.ExecuteTemplate(w, "todo.html", pageData)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	if !isAuthenticated(w, r, session) {
		return
	}

	r.ParseForm()
	todo := strings.TrimSpace(html.EscapeString(r.Form.Get("todo")))
	if todo != "" {
		session.UserAccount.ToDoList.Append(todo)
	}
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
	// POSTメソッドによるリクエストであることの確認
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// セッション情報の取得
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}
	// 認証チェック
	if !isAuthenticated(w, r, session) {
		return
	}

	// POSTパラメータを解析
	r.ParseForm()
	todoId := r.Form.Get("id")
	todo := r.Form.Get("todo")

	// todo項目を更新
	_, err = session.UserAccount.ToDoList.Update(todoId, todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスの返却
	w.WriteHeader(http.StatusOK)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := checkSession(w, r)
	logRequest(r, session)
	if err != nil {
		return
	}

	sessionManager.RevokeSession(w, session.SessionId)
	sessionManager.StartSession(w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
