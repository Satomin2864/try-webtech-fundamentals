// 【注意】本ロジックは勉強のために実装しているが、本実装の時はセキュリティの観点からライブラリを利用し独自実装しないこと
package main

import (
	"net/http"
	"time"
)

const CookieNameSessionId = "sessionId"

// セッション情報を保持する構造体
type HttpSession struct {
	SessionId   string
	Expires     time.Time
	PageData    any
	UserAccount *UserAccount
}

// 新しいセッションを生成する
func NewHttpSession(sessionId string, validityTime time.Duration) *HttpSession {
	session := &HttpSession{
		SessionId: sessionId,
		Expires:   time.Now().Add(validityTime),
		PageData:  "",
	}
	return session
}

// ページデータを削除する
func (s *HttpSession) ClearPageData() {
	s.PageData = ""
}

// セッションIDをCookieへ書き込む
func (s *HttpSession) SetCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     CookieNameSessionId,
		Value:    s.SessionId,
		Expires:  s.Expires,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}
