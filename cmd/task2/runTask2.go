package task2

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"xkcdcomics/cmd"
	"xkcdcomics/model"
)

var templates *template.Template
var deleteMutex sync.Mutex

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

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	deleteMutex.Lock()
	defer deleteMutex.Unlock()

	// Define the number of comics to delete
	numToDelete := 10

	// Retrieve 10 random comics from the database
	randomComics, err := cmd.GetRandomComicsFromDatabase(numToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the selected comics from the database
	err = cmd.DeleteComicsFromDatabase(randomComics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the deleted comics in JSON format
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(randomComics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RunTask2() {
	// Set up database
	var err error
	err = cmd.SetupDatabase()
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	defer cmd.CloseDBConn()

	// Prepare templates
	templates, err = template.ParseGlob("static/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/delete", deleteHandler)

	// Start server
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
