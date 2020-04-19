package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/eze-kiel/freeboard/database"
	"github.com/eze-kiel/freeboard/utils"
	"github.com/gorilla/mux"
)

type Post struct {
	Id   int
	Text string
	Link string
}

type BoardPageData struct {
	PageTitle string
	Posts     []Post
}

type NewPost struct {
	text string
	link string
}

func HandleFunc() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/board", boardPage)
	r.HandleFunc("/post", postPage)

	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	r.PathPrefix("/style/").Handler(http.StripPrefix("/style/", http.FileServer(http.Dir("views/style/"))))

	return r
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func boardPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/board.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	results, err := db.Query("SELECT id, text, link FROM posts ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}

	data := BoardPageData{
		PageTitle: "Board",
		Posts:     []Post{},
	}
	var sqlPost Post
	for results.Next() {
		err = results.Scan(&sqlPost.Id, &sqlPost.Text, &sqlPost.Link)
		if err != nil {
			log.Fatal(err)
		}
		data.Posts = append(data.Posts, Post{Id: sqlPost.Id, Text: sqlPost.Text, Link: sqlPost.Link})
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmpl, err := template.ParseFiles("views/post.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	post := NewPost{
		text: r.FormValue("post"),
		link: r.FormValue("link"),
	}

	if post.link != "" && post.text != "" && utils.IsURL(post.link) {
		_, err = db.Exec(`INSERT INTO posts (text, link) VALUES (?,?)`, post.text, post.link)
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, struct{ Success bool }{true})
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Redirect(w, r, "http://localhost:8080/post", 301)
	}
}
