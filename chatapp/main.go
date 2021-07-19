package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html", "templates/front.html", "templates/data.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title
	if strings.Contains(title, ".txt") {
		body, err := ioutil.ReadFile("./textFiles/" + filename)

		if err != nil {
			return nil, err
		}
		return &Page{Title: title, Body: body}, nil
	} else {
		body, err := ioutil.ReadFile("./textFiles/" + filename + ".txt")
		if err != nil {
			return nil, err
		}
		return &Page{Title: title, Body: body}, nil
	}
}

// なぜ第３引数に型指定してるのにnilを許容する・・・？
// →初期値はnilとPageコンストラクタの型を定義した?
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateForSlice(w http.ResponseWriter, tmpl string, p []*Page) {
	for _, file := range p {
		err := templates.ExecuteTemplate(w, tmpl+".html", file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// http.Redirect(w, r, "/view/"+title, http.StatusFound)

func redirectFrontPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/frontPage", http.StatusFound)
}

func frontHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "front", nil)
}

// 作成されている.txtファイルへのアクセスをまとめる
func dataHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./textFiles/")
	if err != nil {
		log.Fatal(err)
	}
	var textFiles []*Page
	for _, file := range files {
		p, err := loadPage(file.Name())
		if err != nil {
			fmt.Printf("data HandlerdeでErrorが発生している")
			return
		}
		textFiles = append(textFiles, p)
	}
	renderTemplateForSlice(w, "data", textFiles)

}

func main() {

	http.HandleFunc("/", redirectFrontPage)
	http.HandleFunc("/frontPage/", frontHandler)
	http.HandleFunc("/data/", dataHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
