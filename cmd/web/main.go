package main

import (
  "database/sql"
	"flag"
	"log"
	"net/http"
	"os"
  "html/template"

  "snippitbox.chronoabi.com/internal/models"

  _"github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
  snippets *models.SnippetModel
  templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4001", "HTTP network address")

  dsn := flag.String("dsn","web:root_123@/snippetbox?parseTime=true","MySQL data Source name")

	flag.Parse()


	// too log error in more detail
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  db,err := openDB(*dsn)
  if err != nil{
    errorLog.Fatal(err)
  }

  defer db.Close()

  templateCache,err := newTemplateCache() 

  if err != nil {
    errorLog.Fatal(err)
  }

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
    snippets: &models.SnippetModel{DB: db},
    templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server started in port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error){
  db, err := sql.Open("mysql",dsn)
  if err != nil {
    return nil, err
  }

  if err := db.Ping(); err != nil {
    return nil, err
  }

  return db, nil
}
