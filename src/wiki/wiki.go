package main

import (
	"fmt"
	"io/ioutil"
  "net/http"
  "html/template"
  "strings"
)

type Page struct {
  Title string
  Body []byte
}

// Page Methods
func (page * Page) save() error {
  filename := page.Title + ".txt"
  return ioutil.WriteFile("./pages/" + filename, page.Body, 0600)
}

// Helper
func loadPage(title string) (*Page ,error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile("./pages/" + filename)
  if err != nil {
    return nil, err
  }

  return &Page{Title: title, Body: body}, nil
}

func allowBreakline(content string) interface{} {
  safeHTML := template.HTMLEscapeString(content)
  safeHTML = strings.Replace(safeHTML, "\n", "<br />", -1)
  return template.HTML(safeHTML)
}

// Handler
func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/view/"):]
  page, err := loadPage(title)
  if err != nil {
    page = &Page{Title: "404 Not Found", Body: []byte("Sorry this page can't found!")}
  }

  t := template.New("")
  t = t.Funcs(template.FuncMap{"br": allowBreakline})
  t, _ = t.ParseFiles("./template/view.html", "./template/header.html", "./template/footer.html")
  err = t.ExecuteTemplate(w, "view.html", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/edit/"):]
  page, err := loadPage(title)
  if err != nil {
    page = &Page{Title: title}
  }

  t, _ := template.ParseFiles("./template/edit.html", "./template/header.html", "./template/footer.html")
  t.Execute(w, page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[len("/save/"):]
  body := r.FormValue("body")
  page := &Page{Title: title, Body: []byte(body)}
  page.save()
  http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  fmt.Printf("Server start on :8888\n")
  http.ListenAndServe(":8888", nil)
}
