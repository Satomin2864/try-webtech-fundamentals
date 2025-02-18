// 【注意】本ロジックは勉強のために実装しているが、本実装の時はセキュリティの観点からライブラリを利用し独自実装しないこと
package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"
)

const cookieSessionId = "sessionId"

// セッションの確認処理
// セッションが存在しなければ新しく発行する。
func ensureSession(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie(cookieSessionId)
	if err == http.ErrNoCookie {
		// CookieにセッションIDが入っていない場合は、新規発行して返す
		sessionId, err := startSession(w)
		return sessionId, err
	}
	if err == nil {
		// CookieにセッションIDが入っている場合はそれを返す。
		sessionId := c.Value
		return sessionId, nil
	}
	return "", err
}

// セッションを開始する
func startSession(w http.ResponseWriter) (string, error) {
	sessionId, err := makeSessionId()
	if err != nil {
		return "", err
	}

	cookie := &http.Cookie{
		Name:     cookieSessionId,
		Value:    sessionId,
		Expires:  time.Now().Add(600 * time.Second),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return sessionId, nil
}

// セッションIDを生成する
func makeSessionId() (string, error) {
	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes), nil
}
