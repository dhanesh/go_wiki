package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

//Page being created or requested from the server
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(res http.ResponseWriter, tmpl string, page *Page) {
	tmp, err := template.ParseFiles("./templates/" + tmpl + ".html")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(res, page)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// ViewHandler Used for viewing pages
type ViewHandler struct {
}

func (viewHandler ViewHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/view/"):]
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(res, req, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(res, "view", page)
}

// EditHandler Used to edit pages
type EditHandler struct {
}

func (editHandler EditHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(res, "edit", page)
}

// SaveHandler Used to save pages that have been edited
type SaveHandler struct {
}

func (saveHandler SaveHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/edit/"):]
	body := req.FormValue("body")
	p := Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

func main() {
	http.Handle("/view/", http.Handler(ViewHandler{}))
	http.Handle("/edit/", http.Handler(EditHandler{}))
	http.Handle("/save/", http.Handler(SaveHandler{}))
	http.ListenAndServe(":8080", nil)
}
