package main

import (
	"fmt"
	"net/http"
)

// middleware flow as of now : recoverPanic -> logRequest -> secureHeaders servmux -> application handler

// for this application flow : secureHeaders -> servemux -> application handler
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //CSP error are show in your frontend's developer tools so keep your eye on that as well
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)

	})
}

// Log's all requests, application flow : logRequest -> secureHeaders servmux -> application handler

func (app *application) logRequest(next http.Handler) http.Handler{
  return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){

  app.infoLog.Printf("%s - %s %s %s",r.RemoteAddr,r.Proto,r.Method,r.URL.RequestURI())

    next.ServeHTTP(w,r)
  })
}


//Helps recover panic  

func (app *application) recoverPanic(next http.Handler) http.Handler{
  return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){

    defer func(){
      if err := recover(); err != nil{
        //close connection
        w.Header().Set("Connection","close")

        //report an error
        app.serverError( w,fmt.Errorf("%s",err) )

      }
    }()

    next.ServeHTTP(w,r)
  })
}
