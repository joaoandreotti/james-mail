package pages

import (
  "log"
  "sync"
  "net/http"
  "strings"

  "github.com/gorilla/mux"
)

func GetEmailHtmlList(emailList []string) (emailHtmlList string) {
  for _, email := range emailList {
    emailInfo := strings.Split(email, ";")
    emailAddress := emailInfo[0]
    emailBody := emailInfo[1]
    emailHtmlList += "<td>"+emailAddress+": "+emailBody+"</td><hr>"
  }
  return emailHtmlList
}

func StartHttpServer (r *mux.Router, wg *sync.WaitGroup) {
  defer wg.Done()
  log.Println("Starting server at :8080")
  http.ListenAndServe(":8080", r)
}
