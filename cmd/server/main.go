package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

)

var db *sql.DB

func main() {
	var err error 
	db, err = sql.Open("sqlite3", "questions.db")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
	)`)

	http.HandleFunc("/", hello)
	http.HandleFunc("/add", add)
	http.HandleFunc("/browse", browse)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, `
	<html>
		<body>
			<h1>wanna add a question?</h1>
			<form action="/add" method="get">
				<button type="submit">Add</button>
			</form>
		</body>
	</html>`)
}

func add(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		text := req.FormValue("question")
		if text != "" {
			_, err := db.Exec("INSERT INTO questions (text) VALUES (?)", text)
			if err != nil {
				http.Error(w, "DB insert error", 500)
				return
			}
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return
		}
	}

	tmpl := `
	<html>
		<body>
			<h1>Add Question</h1>
			<form method="POST" action="/add">
				<input type="text" name="question" placeholder="Enter question" required>
				<button type="submit">Submit</button>
			</form>
		</body>
	</html>`
	t := template.Must(template.New("add").Parse(tmpl))
	t.Execute(w, nil)
}

func browse(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT id, text FROM questions ORDER BY id")
	if err != nil {
		http.Error(w, "DB query error", 500)
		return
	}
	defer rows.Close()

	fmt.Fprintf(w, `
	<html>
		<body>
			<h1>All Questions</h1>
			<table border="1" cellpadding="5">
				<tr><th>ID</th><th>Question</th></tr>`)

	for rows.Next() {
		var id int
		var text string
		if err := rows.Scan(&id, &text); err != nil {
			http.Error(w, "DB scan error", 500)
			return
		}
		fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td></tr>", id, text)
	}

	fmt.Fprintf(w, `
			</table>
			<br>
			<a href="/">Back</a>
		</body>
	</html>`)
}
