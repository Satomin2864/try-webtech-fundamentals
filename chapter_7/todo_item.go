package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// ToDo項目を表す構造体
type ToDoItem struct {
	Id   string `json:"id"`
	Todo string `json:"todo"`
}

// 新しいToDoItemを生成する
func NewToDoItem(todo string) *ToDoItem {
	id := MakeTodoId(todo)
	return &ToDoItem{
		Id:   id,
		Todo: todo,
	}
}

// ToDo項目のIDを作成する
func MakeTodoId(todo string) string {
	timeBytes := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	hasher := md5.New()
	hasher.Write(timeBytes)
	hasher.Write([]byte(todo))
	return hex.EncodeToString(hasher.Sum(nil))
}
