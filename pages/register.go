package pages

import (
  "net/http"
  "html/template"
  "james-mail/middleware"
  "james-mail/database"
  "james-mail/validator"
  "james-mail/authentication"
  "github.com/gorilla/mux"
)

const registerPagePath string = "/register"

const registerFormUsername string = "username"
const registerFormEmail string = "email"
const registerFormPassword string = "password"

type RegisterPageData struct {
  RegisterFormUsername string
  RegisterFormEmail string
  RegisterFormPassword string
  RegisterPagePath string
}

func GetRegisterPageData() (data RegisterPageData) {
  data = RegisterPageData {
    RegisterFormUsername: registerFormUsername,
    RegisterFormEmail: registerFormEmail,
    RegisterFormPassword: registerFormPassword,
    RegisterPagePath: registerPagePath,
  }
  return data
}

func RegisterGet(w http.ResponseWriter, r *http.Request) {
  if authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, emailsPagePath, 302)
    return
  }

  page := template.Must(template.ParseFiles("pages/html/register.html"))
  data := GetRegisterPageData()
  page.Execute(w, data)
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
  if authentication.ValidateAuthentication(r) {
    http.Redirect(w, r, emailsPagePath, 302)
    return
  }

  data := database.UserRegistrationData {
    Username: r.FormValue(registerFormUsername),
    Email: r.FormValue(registerFormEmail),
    Password: r.FormValue(registerFormPassword),
  }
  if !validator.ValidateStruct(data) {
    http.Redirect(w, r, errorPagePath, 302)
    return
  }

  database.InsertNewUser(data.Username, data.Email, data.Password)
  http.Redirect(w, r, loginPagePath, 302)
}

func ConfigRegister(r *mux.Router) {
  r.HandleFunc(registerPagePath, middleware.Chain(RegisterGet)).Methods("GET")
  r.HandleFunc(registerPagePath, middleware.Chain(RegisterPost)).Methods("POST")
}
