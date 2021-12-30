package middleware

import (
  "log"
  "net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging() Middleware {
  return func(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      log.Println(r.URL.Path)
      f(w, r)
    }
  }
}

func BasicChain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
  for _, m := range middlewares {
    f = m(f)
  }
  return f
}

func Chain(response http.HandlerFunc) http.HandlerFunc {
  return BasicChain(response, Logging())
}
