package handlers

import (
	"html/template"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/eze-kiel/freeboard/database"
	"github.com/eze-kiel/freeboard/utils"
	"github.com/gorilla/mux"
)

// Post represents a post which will be displayed
type Post struct {
	ID       int
	Text     string
	Link     string
	Category string
}

// BoardPageData contains the data sent to the board page
type BoardPageData struct {
	PageTitle string
	Posts     []Post
}

// NewPost contains data sent via the from in Post section
type NewPost struct {
	text     string
	link     string
	category string
}

// HandleFunc handles functions
func HandleFunc() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/post", postPage)
	r.HandleFunc("/rules", rulesPage)
	r.HandleFunc("/random", randomPage)
	r.HandleFunc("/about", aboutPage)
	r.NotFoundHandler = http.HandlerFunc(notFoundPage)
	r.HandleFunc("/boards/{category}", boardsPage)

	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	r.PathPrefix("/style/").Handler(http.StripPrefix("/style/", http.FileServer(http.Dir("./style/"))))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	return r
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home.html", "views/templates/head.html", "views/templates/header.html", "views/templates/announcements.html")
	if err != nil {
		log.Fatalf("Can not parse home page : %v", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for home page : %v", err)
	}
}

func boardsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/boards.html", "views/templates/head.html", "views/templates/header.html")
	if err != nil {
		log.Fatalf("Can not parse board page : %v", err)
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get category from the URL
	vars := mux.Vars(r)
	category := vars["category"]

	// Redirect client to /all if a wrong url is entered
	if utils.CheckCategory(category) != true {
		http.Redirect(w, r, "/boards/all", 301)
	}

	data := BoardPageData{
		PageTitle: strings.Title(category),
		Posts:     []Post{},
	}

	// Send all the content of the database
	if category == "all" {
		results, err := db.Query("SELECT id, text, link, category FROM posts ORDER BY id DESC")
		if err != nil {
			log.Fatal(err)
		}
		var sqlPost Post
		for results.Next() {
			err = results.Scan(&sqlPost.ID, &sqlPost.Text, &sqlPost.Link, &sqlPost.Category)
			if err != nil {
				log.Fatal(err)
			}
			data.Posts = append(data.Posts, Post{ID: sqlPost.ID, Text: sqlPost.Text, Link: sqlPost.Link, Category: sqlPost.Category})
		}

		// Send only the content of the requested category
	} else {
		results, err := db.Query("SELECT id, text, link, category FROM posts WHERE category= ? ORDER BY id DESC", category)
		if err != nil {
			log.Fatal(err)
		}
		var sqlPost Post
		for results.Next() {
			err = results.Scan(&sqlPost.ID, &sqlPost.Text, &sqlPost.Link, &sqlPost.Category)
			if err != nil {
				log.Fatal(err)
			}
			data.Posts = append(data.Posts, Post{ID: sqlPost.ID, Text: sqlPost.Text, Link: sqlPost.Link, Category: sqlPost.Category})
		}
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatalf("Can not execute templates for board page : %v", err)
	}
}

func postPage(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmpl, err := template.ParseFiles("views/post.html", "views/templates/head.html", "views/templates/header.html")
	if err != nil {
		log.Fatalf("Can not parse post page : %v", err)
	}

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	post := NewPost{
		text:     r.FormValue("post"),
		link:     r.FormValue("link"),
		category: r.FormValue("category"),
	}

	// Integrity check
	// Check if link is empty or if text is empty or is url is not an url or if text length is < 500 characters
	if post.link != "" && post.text != "" && utils.IsURL(post.link) && len(post.text) <= 500 {

		// Content check
		// Check if the request contain censured content, or a non-existent category
		if utils.AuthorizedURL(post.link) && utils.AuthorizedText(post.text) && utils.CheckCategory(post.category) {
			_, err = db.Exec(`INSERT INTO posts (text, link, category) VALUES (?,?,?)`, post.text, post.link, post.category)
			if err != nil {
				log.Fatal(err)
			}

			err = tmpl.Execute(w, struct{ Success bool }{true})
			if err != nil {
				log.Fatalf("Can not execute templates for post page : %v", err)
			}
		} else {
			http.Redirect(w, r, "/post", 301)
		}

	} else {
		http.Redirect(w, r, "/post", 301)
	}
}

func rulesPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/rules.html", "views/templates/head.html", "views/templates/header.html")
	if err != nil {
		log.Fatalf("Can not parse rules page : %v", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for rules page : %v", err)
	}
}

func randomPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/random.html", "views/templates/head.html", "views/templates/header.html")
	if err != nil {
		log.Fatalf("Can not parse random page : %v", err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Sort a random element frmo the database
	results, err := db.Query("SELECT id, text, link, category FROM posts ORDER BY RAND() LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}

	data := BoardPageData{
		PageTitle: "Random",
		Posts:     []Post{},
	}
	var sqlPost Post
	for results.Next() {
		err = results.Scan(&sqlPost.ID, &sqlPost.Text, &sqlPost.Link, &sqlPost.Category)
		if err != nil {
			log.Fatal(err)
		}
		data.Posts = append(data.Posts, Post{ID: sqlPost.ID, Text: sqlPost.Text, Link: sqlPost.Link, Category: sqlPost.Category})
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Fatalf("Can not execute templates for random page : %v", err)
	}
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/404.html", "views/templates/head.html", "views/templates/header.html")
	if err != nil {
		log.Fatalf("Can not parse 404 page : %v", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for 404 page : %v", err)
	}
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/about.html", "views/templates/head.html", "views/templates/header.html")
	if err != nil {
		log.Fatalf("Can not parse about page : %v", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for about page : %v", err)
	}
}
