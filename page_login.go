package main

import (
	"log"
	"net/http"
)

type LoginPageData struct {
	UserId       string
	ErrorMessage string
}

// ログインに関するリクエスト情報
func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := ensureSession(w, r)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		showLogin(w, r, session)
		return
	case http.MethodPost:
		login(w, r, session)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// ログイン画面を表示する
// TODO: 第２引数の使ってない気がする？確認して連絡
func showLogin(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	var pageData LoginPageData
	if p, ok := session.PageData.(LoginPageData); ok {
		pageData = p
	} else {
		pageData = LoginPageData{}
	}
	templates.ExecuteTemplate(w, "login.html", pageData)
	session.ClearPageData()
}

// ログイン処理
func login(w http.ResponseWriter, r *http.Request, session *HttpSession) {
	// POSTパラメータ処理
	r.ParseForm()
	userId := r.Form.Get("userId")
	password := r.Form.Get("password")

	// 認証
	log.Printf("login attempt: %s\n", userId)
	account, err := accountManager.Authenticate(userId, password)
	if err != nil {
		log.Printf("login failed: %s\n", err)
		session.PageData = LoginPageData{
			ErrorMessage: "ユーザIDまたはパスワードが違います",
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// ログイン成功
	session.UserAccount = account

	log.Printf("login succeed: %s\n", account.Id)
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
	return

}
