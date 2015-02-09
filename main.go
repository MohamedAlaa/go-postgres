package main

import (
  "database/sql"
  "github.com/codegangsta/martini"
  "github.com/martini-contrib/render"
  _ "github.com/lib/pq"
  "net/http"
)

type Book struct {
  Title string
  Author string
  Description string
}

func SetupDB() *sql.DB {
  db, err := sql.Open("postgres", "dbname=lesson4 sslmode=disable")
  PanicIf(err)
  return db
}

func PanicIf(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {
  m := martini.Classic()
  m.Map(SetupDB())
  m.Use(render.Renderer(render.Options{
    Layout: "layout",
  }))

  m.Get("/", ShowBooks)
  m.Get("/create", NewBook)
  m.Post("/books", Create)

  m.Run()
}

func ShowBooks(ren render.Render, r *http.Request, db *sql.DB) {

  searchTerm := "%" + r.URL.Query().Get("search") + "%"

  rows, err := db.Query(`SELECT title, author, description FROM books
    WHERE author ILIKE $1
    OR author ILIKE $1
    OR description ILIKE $1`, searchTerm)

  PanicIf(err)
  defer rows.Close()

  books := []Book{}
  for rows.Next() {
    b := Book{}
    PanicIf(rows.Err())
    err := rows.Scan(&b.Title, &b.Author, &b.Description)
    PanicIf(err)
    books = append(books, b)
  }

  ren.HTML(200, "books", books)

}

func NewBook(ren render.Render) {
  ren.HTML(200, "create", nil)
}

func Create(ren render.Render, r *http.Request, db *sql.DB) {
  rows, err := db.Query("INSERT INTO books (title, author, description) VALUES ($1, $2, $3)",
    r.FormValue("title"),
    r.FormValue("author"),
    r.FormValue("description"),
  )

  PanicIf(err)
  defer rows.Close()

  ren.Redirect("/")
}