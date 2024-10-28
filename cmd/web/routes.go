package main

import (
	"net/http"
"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
  router := httprouter.New()

  router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    app.notFound(w)
  })

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet,"/", app.home)
  router.HandlerFunc(http.MethodGet,"/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet,"/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost,"/snippet/create", app.snippetCreatePost)

	// a convinent way of writing middlewares then just returning function inside a function.
	// it shows chains more properly
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// it will go through the secureHeaders middleware and add those OWASP secure headers
	return standard.Then(router)
}
