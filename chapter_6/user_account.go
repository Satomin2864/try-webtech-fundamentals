package main

import (
	// パスワードをハッシュ化するためのライブラリ
	"golang.org/x/crypto/bcrypt"
	"time"
	_ "time/tzdata"
)

const UserAccountLimitInMinute = 60
const PasswordLength = 10
const PasswordChars = "23456789abcdefghijkmnpqrstuvwxyz"

// ユーザ情報を表現する構造体
type UserAccount struct {
	Id             string
	HashedPassword string
	Expires        time.Time
	ToDoList       []string
}

func NewUserAccount(userId string, plainPassword string, expires time.Time) *UserAccount {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), PasswordLength)
	account := &UserAccount{
		Id:             userId,
		HashedPassword: string(hashedPassword),
		Expires:        expires,
		ToDoList:       make([]string, 0, 10),
	}
	return account
}

func (u UserAccount) ExpiresText() string {
	return u.Expires.Format("2006/01/02 15:04:05")
}
