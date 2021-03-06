package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
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
	PageTitle    string
	Posts        []Post
	IsPrevPage   bool
	IsNextPage   bool
	CurrentPage  int
	PreviousPage int
	NextPage     int
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
	r.HandleFunc("/random", randomPage)

	// NotFoundHandler handles routes to /about, /rules and not found
	r.NotFoundHandler = http.HandlerFunc(defaultPage)

	r.HandleFunc("/boards/{category}/{page}", boardsPage)

	r.HandleFunc("/boards/{category}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			category := vars["category"]
			http.Redirect(w, r, "/boards/"+category+"/0", 301)
		})

	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	r.PathPrefix("/style/").Handler(http.StripPrefix("/style/", http.FileServer(http.Dir("./style/"))))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	return r
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	for _, page := range []string{"about", "rules"} {
		if strings.Contains(r.RequestURI, page) {
			tmpl, err := template.ParseFiles("views/"+page+".html", "views/templates/head.html", "views/templates/header.html")
			if err != nil {
				log.Fatalf("Can not parse "+page+" page : %v", err)
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				log.Fatalf("Can not execute templates for "+page+" page : %v", err)
			}
			return
		}
	}

	// Redirect to 404 page if not found
	tmpl, err := template.ParseFiles("views/404.html", "views/templates/head.html", "views/templates/header.html")
	if err != nil {
		log.Fatalf("Can not parse home page : %v", err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalf("Can not execute templates for home page : %v", err)
	}
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

	// Get category and page from the URL
	vars := mux.Vars(r)
	category := vars["category"]
	page := vars["page"]
	limit := 15

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		log.Warn(err)
	}

	offset := pageNumber * limit

	//Redirect client to /all if a wrong url is entered
	if utils.CheckCategory(category) != true {
		http.Redirect(w, r, "/boards/all/0", 301)
	}

	data := BoardPageData{
		PageTitle:    category,
		Posts:        []Post{},
		CurrentPage:  pageNumber,
		PreviousPage: pageNumber - 1,
		NextPage:     pageNumber + 1,
	}

	var results *sql.Rows

	if category == "all" {
		results, err = db.Query("SELECT id, text, link, category FROM posts ORDER BY id DESC LIMIT ? OFFSET ?", limit+1, offset)
	} else {
		results, err = db.Query("SELECT id, text, link, category FROM posts WHERE category= ? ORDER BY id DESC LIMIT ? OFFSET ?", category, limit+1, offset)
	}

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

	if len(data.Posts) == limit+1 {
		data.IsNextPage = true
		// Index begins at 0 but a half-open range
		data.Posts = data.Posts[0:limit]
	} else {
		data.IsNextPage = false
	}

	if data.PreviousPage < 0 {
		data.IsPrevPage = false
	} else {
		data.IsPrevPage = true
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

	// utils.CanClientPost returns false if a RemoteAddr can not post
	if !utils.CanClientPost(r.RemoteAddr) {
		// Execute Anti spam alert message
		err = tmpl.Execute(w, struct {
			Success    bool
			Timeout    bool
			BadContent bool
		}{false, true, false})

		if err != nil {
			log.Fatalf("Can not execute templates for post page : %v", err)
		}

		return
	}

	// Add RemoteAddr to AntiSpam list
	utils.AddIPToAntiSpam(r.RemoteAddr)

	// Integrity check
	// Check if link is empty or if text is empty or is url is not an url or if text length is < 500 characters
	integrityCheck := true

	switch {
	case post.link == "":
		integrityCheck = false
		break
	case post.text == "":
		integrityCheck = false
		break
	case !utils.IsURL(post.link):
		integrityCheck = false
		break
	case len(post.text) > 500:
		integrityCheck = false
		break
	}

	if !integrityCheck {
		// Something gone wrong, but a normal user should never arrive here
		// so there is no alert message, but BadRequest status code (400)
		http.Redirect(w, r, "/post", 400)
	}

	// Content check
	// Check if the request contain censured content, or a non-existent category
	contentCheck := true
	switch {
	case !utils.AuthorizedURL(post.link):
		contentCheck = false
		break

	case !utils.AuthorizedText(post.text):
		contentCheck = false
		break

	case !utils.CheckCategory(post.category):
		contentCheck = false
		break
	}

	if !contentCheck {
		// Execute Bad Content alert message
		err = tmpl.Execute(w, struct {
			Success    bool
			Timeout    bool
			BadContent bool
		}{false, false, true})

		if err != nil {
			log.Fatalf("Can not execute templates for post page : %v", err)
		}

		return
	}

	// All verification checks passed
	_, err = db.Exec(`INSERT INTO posts (text, link, category) VALUES (?,?,?)`, post.text, post.link, post.category)
	if err != nil {
		log.Fatal(err)
	}

	// Execute Success alert message
	err = tmpl.Execute(w, struct {
		Success    bool
		Timeout    bool
		BadContent bool
	}{true, false, false})

	if err != nil {
		log.Fatalf("Can not execute templates for post page : %v", err)
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
