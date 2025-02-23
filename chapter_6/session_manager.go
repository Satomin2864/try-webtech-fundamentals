package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

// エラー定義
var (
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidSession  = errors.New("invalid session id")
)

// セッションを管理する構造体
type HttpSessionManager struct {
	// セッションIDをキーとして、セッション情報保持する
	sessions map[string]*HttpSession
}

func NewHttpSessionManager() *HttpSessionManager {
	mgr := &HttpSessionManager{
		sessions: make(map[string]*HttpSession),
	}
	return mgr
}

// セッション開始。CookieにセッションIDを書き込む
func (m *HttpSessionManager) StartSession(w http.ResponseWriter) (*HttpSession, error) {
	// 新しいセッションIDを生成する
	sessionId, err := m.makeSessionId()
	if err != nil {
		return nil, err
	}

	// セッション情報を生成する
	log.Printf("start session : %s", sessionId)
	session := NewHttpSession(sessionId, 10*time.Minute)
	m.sessions[sessionId] = session
	session.SetCookie(w)

	return session, nil
}

// セッションIDを生成する
func (m *HttpSessionManager) makeSessionId() (string, error) {
	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes), nil
}

// セッションを削除する
func (m *HttpSessionManager) RevokeSession(w http.ResponseWriter, sessionId string) {
	// セッション情報を削除
	delete(m.sessions, sessionId)
	log.Printf("session revoked : %s", sessionId)

	if w == nil {
		return
	}
	// MaxAgeに負の値を設定すると、すぐに該当のクッキーは廃棄される
	cookie := &http.Cookie{
		Name:    CookieNameSessionId,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, cookie)
}

// セッションが存在するかチェックする
// 存在しない場合にはログイン画面日リダイレクト
func checkSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	// CookieのセッションIDに紐づくセッション情報を取得
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		// セッション情報を取得できたら終了
		return session, nil
	}
	orgErr := err

	//  セッションが有効期限切れ、または不正な場合にはセッションを作り直す
	log.Printf("session check failed: %s", err.Error())
	session, err = sessionManager.StartSession(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// Refererヘッダの有無で他画面からの遷移か判定
	// アプリケーショントップのURLに直接アクセスした際にはセッションが存在しないのが正常であるため、エラーを表示しないための措置
	if r.Referer() != "" {
		page := LoginPageData{}
		page.ErrorMessage = "セッションが不正です。"
		session.PageData = page
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil, orgErr
}

// cookieから有効なセッション情報を取得する
func (m *HttpSessionManager) GetValidSession(r *http.Request) (*HttpSession, error) {
	c, err := r.Cookie(CookieNameSessionId)
	// CookieにセッションIDが存在しない場合
	if err == http.ErrNoCookie {
		return nil, ErrSessionNotFound
	}

	// Cookieにセッション情報が存在している場合
	if err == nil {
		// セッション情報を取得し返す
		sessionId := c.Value
		session, err := m.getSession(sessionId)
		return session, err
	}
	return nil, err
}

// セッションIDに紐づくセッション情報を返す
func (m *HttpSessionManager) getSession(sessionId string) (*HttpSession, error) {
	if session, exists := m.sessions[sessionId]; exists {
		// セッション情報の有効期限チェック
		if time.Now().After(session.Expires) {
			// 期限切れの場合にはエラーを返す
			delete(m.sessions, sessionId)
			return nil, ErrSessionExpired
		}
		return session, nil
	} else {
		return nil, ErrSessionNotFound
	}
}

// セッションが開始されていることを保証する
// 存在していない場合には新しく発行
func ensureSession(w http.ResponseWriter, r *http.Request) (*HttpSession, error) {
	session, err := sessionManager.GetValidSession(r)
	if err == nil {
		return session, nil
	}

	// セッションが存在しないか不正な場合は新しく開始
	log.Printf("session check failed: %s", err.Error())
	session, err = sessionManager.StartSession(w)
	if err != nil {
		writeInternalServerError(w, err)
		return nil, err
	}
	return session, err
}
