package task2

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"xkcdcomics/cmd"
	"xkcdcomics/model"
)

var templates *template.Template

func indexHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * 10
	comics, err := cmd.GetComicsFromDatabase(offset, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prevPage := page - 1
	if prevPage <= 0 {
		prevPage = 1
	}
	nextPage := page + 1

	err = templates.ExecuteTemplate(w, "index.html", struct {
		Comics   []model.Comic
		PrevPage int
		NextPage int
	}{Comics: comics, PrevPage: prevPage, NextPage: nextPage})
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func RunTask2() {
	// Set up database
	var err error
	cmd.Db, err = cmd.SetupDatabase()
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	defer cmd.Db.Close()

	// Prepare templates
	templates, err = template.ParseGlob("static/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up routes
	http.HandleFunc("/", indexHandler)

	// Start server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
