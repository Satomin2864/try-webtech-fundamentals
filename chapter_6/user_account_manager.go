package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"regexp"
	"time"
)

var (
	// エラー定義一覧
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidUserIdFormat = errors.New("invalid user id format")
	ErrLoginFailed         = errors.New("login failed")

	// アカウント名チェック用の正規表現フォーマット
	RegexAccountId = regexp.MustCompile(`^[A-Za-z0-9_.+@-]{1,32}$`)
)

// ユーザアカウントを管理する構造体
type UserAccountManager struct {
	userAccounts map[string]*UserAccount
	location     *time.Location
}

// ユーザアカウントマネージャーの生成関数
func NewUserAccountManager() *UserAccountManager {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	m := &UserAccountManager{
		userAccounts: make(map[string]*UserAccount),
		location:     jst,
	}
	return m
}

// ユーザIDの形式を検証する
func (m *UserAccountManager) ValidateUserId(userId string) bool {
	return RegexAccountId.MatchString(userId)
}

// 新しいユーザアカウント作成する
func (m *UserAccountManager) NewUserAccount(userId string, password string) (*UserAccount, error) {
	// ユーザー名のフォーマットが不適切な場合にはエラーを返す
	if !m.ValidateUserId(userId) {
		return nil, ErrInvalidUserIdFormat
	}
	// 重複しているアカウントを作ると困るので、エラーを返す
	_, exists := m.userAccounts[userId]
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// アカウント生成
	expires := time.Now().In(m.location).Add(time.Minute * UserAccountLimitInMinute)
	account := NewUserAccount(userId, password, expires)

	// アカウントリストを登録
	// DBがあるシステムだと、このタイミングでDB登録とか認証確認メール送りそう
	m.userAccounts[userId] = account
	log.Printf("user account created: %s\n", account.Id)
	return account, nil
}

// ユーザアカウントを取得する
func (m UserAccountManager) GetUserAccount(userId string) (*UserAccount, bool) {
	a, exists := m.userAccounts[userId]
	return a, exists
}

// ユーザアカウントを認証する]
func (m *UserAccountManager) Authenticate(userId string, password string) (*UserAccount, error) {
	// アカウントの存在チェック
	account, exists := m.GetUserAccount(userId)
	if !exists {
		return nil, ErrLoginFailed
	}

	// パスワードチェック
	err := bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(password))
	if err != nil {
		return nil, ErrLoginFailed
	}
	return account, nil
}

// ランダムなパスワードを生成する.
func MakePassword() string {
	password := make([]byte, PasswordLength)
	for i := 0; i < PasswordLength; i++ {
		password[i] = PasswordChars[rand.Intn(len(PasswordChars))]
	}
	return string(password)
}
