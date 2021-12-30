package authentication

import (
  "net/http"
  "github.com/gorilla/sessions"
  "fmt"
)

var authCookie string = "Authentication"
var authSecretKey string = "lasidfj198034y1r j34r190u34r891u34m1r89p 4rh1780 r"

var (
  key = []byte(authSecretKey)
  store = sessions.NewCookieStore(key)
)

func ValidateAuthentication(r *http.Request) bool {
  session, _ := store.Get(r, authCookie)
  auth, ok := session.Values["authenticated"].(bool);
  return ok && auth
}

func Authenticate(w http.ResponseWriter, r *http.Request, username, email string) {
  session, _ := store.Get(r, authCookie)
  session.Values["authenticated"] = true
  session.Values["username"] = username
  session.Values["email"] = email
  session.Save(r, w)
}

func Unauthenticate(w http.ResponseWriter, r *http.Request) {
  session, _ := store.Get(r, authCookie)
  session.Values["authenticated"] = false
  session.Save(r, w)
}

func GetSessionValue(r *http.Request, value string) string {
  session, _ := store.Get(r, authCookie)
  return fmt.Sprint(session.Values[value])
}
