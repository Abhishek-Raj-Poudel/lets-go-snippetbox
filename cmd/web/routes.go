package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	//Unprotected routes will use the normal dynamic middleware
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	//snippets
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	//authentication
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.useSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.useSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))


	//For protected ones we will first add the requireAuthntication middleware in dynamic
	protected := dynamic.Append(app.requireAuthentication)

	//snippets
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	//authentication
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogOutPost))

	// a convinent way of writing middlewares then just returning function inside a function.
	// it shows chains more properly
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// it will go through the secureHeaders middleware and add those OWASP secure headers
	return standard.Then(router)
}
