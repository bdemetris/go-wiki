package wiki

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("static/edit.html", "static/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	p, err := GetPages()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("static/index.html")
	t.Execute(w, p)
}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := GetPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := GetPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	title = r.FormValue("title")
	p := &Page{Title: string(title), Body: []byte(body)}
	err := Save(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
