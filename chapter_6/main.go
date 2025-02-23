package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	sessionManager *HttpSessionManager
	accountManager *UserAccountManager
	templates      *template.Template

	ErrMethodNotAllowed = errors.New("method not allowed")
)

func main() {
	sessionManager = NewHttpSessionManager()
	accountManager = NewUserAccountManager()

	// テンプレート読み込み
	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/create-user-account", handleCreateUserAccount)
	http.HandleFunc("/new-user-account", handleNewUserAccount)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/todo", handleTodo)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/favicon.ico", handleNotFound)
	http.HandleFunc("/", handleRoot)

	port := getPortNumber()
	fmt.Printf("listening port : %d\n", port)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}

// ログイン画面へリダイレクトする
func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

// HTTPリクエストが指定したメソッドかチェックする
func checkMethod(w http.ResponseWriter, r *http.Request, method string) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return ErrMethodNotAllowed
	}
	return nil
}

// エラー出力
func writeInternalServerError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("500 Internal Server Error: %s", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(msg))
}

// ログイン済みかチェックする
func isAuthenticated(w http.ResponseWriter, r *http.Request, session *HttpSession) bool {
	if session.UserAccount != nil {
		return true
	}

	page := LoginPageData{}
	page.ErrorMessage = "未ログインです。"
	session.PageData = page

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return false
}
