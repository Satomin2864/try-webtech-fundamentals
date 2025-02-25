// defer属性によってDOM構築が完了した後にsetup関数を呼び出す
setup();

// 編集しているToDo項目のinput要素
let currentEditingTodo;

// 各種イベントハンドラを設定する
function setup() {
    addEventListenerByQuery('input[type="text"].todo', "click", onClickTodoInput);
    addEventListenerByQuery('button.btn-cancel', "click", onClickCancelButton);
    addEventListenerByQuery('button.btn-save', "click", onClickSaveButton);
}

// 関数をイベントハンドラとして登録する
// @param {string} 登録先要素を指定するクエリ
// @param {string} イベント名
// @param {function} イベントハンドラとして登録する関数
function addEventListenerByQuery(query, eventName, callback) {
    const elements = document.querySelectorAll(query);
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener(eventName, callback);
    }
}

// ToDo項目の入力欄がクリックされた際のイベントハンドラ
// @param {Event} event イベントオブジェクト
function onClickTodoInput(event) {
    if (isNotNull(currentEditingTodo) && currentEditingTodo !== event.target) {
        cancelTodoEdit(currentEditingTodo);
    }

    enableTodoInput(event.target);
    currentEditingTodo = event.target;
}


// todo項目編集可能にする
// save/canleボタンを表示する
// @param {HTMLInputElement} todo項目のinput要素
function enableTodoInput(todoInput) {
    // キャンセル時に戻すため、編集前の内容を保存
    if (isNull(todoInput.dataset.originalValue)) {
        todoInput.dataset.originalValue = todoInput.value;
    }

    todoInput.readOnly = false;

    // save/cancelボタンを表示
    const todoEditorControl = getTodoEditorControl(todoInput);
    if (isNotNull(todoEditorControl)) {
        todoEditorControl.classList.remove("hidden");
    } else {
        console.error("TodoEditorControl not found.");
    }
}

// ToDo項目の入力欄を編集不可にする
// save/cancelボタンを非表示にする
// @param {HTMLInputElement} todo項目のinput要素
function disableTodoInput(todoInput) {
    todoInput.readOnly = true;

    // save/cancelボタンを非表示
    const todoEditorControl = getTodoEditorControl(todoInput);
    if (isNotNull(todoEditorControl)) {
        todoEditorControl.classList.add("hidden");
    } else {
        console.error("TodoEditorControl not found.");
    }
}
// todo項目に対応するsave/cancelボタンのコンテナ要素を取得する
// @param {HTMLInputElement} todo項目のinput要素
// @returns {HTMLElement} save/cancelボタンの親div要素、存在しない場合はnull
function getTodoEditorControl(todoInput) {
    const todoEditorControl = todoInput.nextElementSibling;
    if (todoEditorControl.classList.contains("todo-item-control")) {
        return todoEditorControl;
    }
    return null;
}

// キャンセルボタンがクリックされた際のイベントハンドラ
// @param {Event} イベントオブジェクト
function onClickCancelButton(event){
    const cancelBtn = event.target;
    const todoInput = cancelBtn.parentNode.previousElementSibling;
    cancelToDoEdit(todoInput);
}

// todo項目の編集をキャンセルする
// @param {HTMLInputElement} todo項目のinput要素
function cancelToDoEdit(todoInput){
    disableTodoInput(todoInput);
    // 編集要素を元に戻す
    todoInput.value = todoInput.dataset.originalValue;
    currentEditingTodo = null;
}

// saveボタンがクリックされた際のイベントハンドラ
// @param {Event} イベントオブジェクト
function onClickSaveButton(event){
    // ブラウザ上で押されたSaveボタンに対応する要素を取得
    const saveBtn = event.target;
    // saveボタンに対応するInput要素を取得
    const todoInput = saveBtn.parentNode.previousElementSibling;
    todoInput.blur();

    // リクエストの準備
    const request = {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        },
        body: `id=${todoInput.id}&todo=${todoInput.value}`
    }

    // Fetch APIを使ってサーバにリクエストを送信
    fetch("/edit", request)
        .then(()=>{
            // 編集が完了したらInput要素を編集不可にする
            disableTodoInput(todoInput);
            currentEditingTodo = null;
        })
        .catch((error)=>{
            console.error("Request failed", error);
        });
}


// 値が null or undefined か判定をする
// @param value 判定する値
// @returns {boolean} null または undefined の場合はtrue
function isNull(value){
    return value === null || typeof value === "undefined";
}

// 値が null or undefineddでないか判定をする
// @param value 判定する値
// @returns {boolean} null または undefined の場合はfalse
function isNotNull(value) {
    return !isNull(value);
}