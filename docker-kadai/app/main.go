package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	ID   int
	Task string
	Done bool
}

var db *sql.DB

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=true", dbUser, dbPassword, dbHost, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/toggle", handleToggle)
	http.HandleFunc("/delete", handleDelete)

	fmt.Println("サーバーを8080番ポートで起動しています...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, task, done FROM todos ORDER BY id")
	if err != nil {
		http.Error(w, fmt.Sprintf("DBエラー：%v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Task, &t.Done); err != nil {
			http.Error(w, fmt.Sprintf("DBエラー：%v", err), http.StatusInternalServerError)
			return
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("DBエラー：%v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "テンプレートエラー", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, todos); err != nil {
		log.Printf("テンプレート出力エラー: %v", err)
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		task := r.FormValue("task")
		if task != "" {
			if _, err := db.Exec("INSERT INTO todos (task) VALUES (?)", task); err != nil {
				http.Error(w, fmt.Sprintf("DBエラー：%v", err), http.StatusInternalServerError)
				return
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if _, err := db.Exec("UPDATE todos SET done = NOT done WHERE id = ?", r.FormValue("id")); err != nil {
			http.Error(w, fmt.Sprintf("DBエラー：%v", err), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if _, err := db.Exec("DELETE FROM todos WHERE id = ?", r.FormValue("id")); err != nil {
			http.Error(w, fmt.Sprintf("DBエラー：%v", err), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>ToDoアプリ</title>
    <style>
        body { font-family: sans-serif; max-width: 600px; margin: 2rem auto; padding: 0 1rem; }
        .todo-item { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #ccc; padding: 10px 0; }
        .done { text-decoration: line-through; color: #888; }
        button { cursor: pointer; }
    </style>
</head>
<body>
    <h1>ToDoアプリ</h1>
    <form action="/add" method="POST">
        <input type="text" name="task" required placeholder="新しいタスクを入力">
        <button type="submit">追加</button>
    </form>
    <div id="list">
        {{range .}}
        <div class="todo-item">
            <span class="{{if .Done}}done{{end}}">{{.Task}}</span>
            <div>
                <form action="/toggle" method="POST" style="display:inline;">
                    <input type="hidden" name="id" value="{{.ID}}">
                    <button type="submit">{{if .Done}}戻す{{else}}完了{{end}}</button>
                </form>
                <form action="/delete" method="POST" style="display:inline;">
                    <input type="hidden" name="id" value="{{.ID}}">
                    <button type="submit">削除</button>
                </form>
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>
`
