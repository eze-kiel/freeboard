package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/eze-kiel/freeboard/database"
	"github.com/eze-kiel/freeboard/utils"
	"github.com/gorilla/mux"
)

// Post represents a post which will be displayed
type Post struct {
	ID   int
	Text string
	Link string
}

// BoardPageData contains the data sent to the board page
type BoardPageData struct {
	PageTitle string
	Posts     []Post
}

// NewPost contains data sent via the from in Post section
type NewPost struct {
	text string
	link string
}

// HandleFunc handles functions
func HandleFunc() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/board", boardPage)
	r.HandleFunc("/post", postPage)
	r.HandleFunc("/rules", rulesPage)

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
		err = results.Scan(&sqlPost.ID, &sqlPost.Text, &sqlPost.Link)
		if err != nil {
			log.Fatal(err)
		}
		data.Posts = append(data.Posts, Post{ID: sqlPost.ID, Text: sqlPost.Text, Link: sqlPost.Link})
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

	if post.link != "" && post.text != "" && utils.IsURL(post.link) && len(post.text) <= 500 {
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

func rulesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/rules.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
