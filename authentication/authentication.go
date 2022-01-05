package authentication

import (
  "fmt"
  "net/http"
  "github.com/gorilla/sessions"
  "github.com/gorilla/securecookie"
)

const authCookie string = "Authentication"
var authSecretKey []byte = securecookie.GenerateRandomKey(32)

var (
  key = authSecretKey
  store = sessions.NewCookieStore(key)
)

func GenereateAuthSecretKey() {
}

func ValidateAuthentication(r *http.Request) bool {
  session, _ := store.Get(r, authCookie)
  auth, ok := session.Values["authenticated"].(bool);
  return ok && auth
}

func Authenticate(w http.ResponseWriter, r *http.Request, username, email, password string) {
  session, _ := store.Get(r, authCookie)
  session.Values["authenticated"] = true
  session.Values["username"] = username
  session.Values["email"] = email
  session.Values["password"] = password
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
