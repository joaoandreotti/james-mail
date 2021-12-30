package pages

import (
  "james-mail/database"
  "james-mail/pgp"
)

func GetEmailHtmlList(encryptedEmailList []database.StoredEmail, privateKey, password string) (emailHtmlList string) {
  for _, encryptedEmail := range encryptedEmailList {
    decryptedEmail := pgp.DecryptMessage(privateKey, encryptedEmail.EmailBody, password)
    if len(decryptedEmail) > 0 {
      emailHtmlList += "<td>"+encryptedEmail.EmailAddress+": "+decryptedEmail+"</td><hr>"
    } else {
      emailHtmlList += "<td>"+encryptedEmail.EmailAddress+": "+encryptedEmail.EmailBody+"</td><hr>"
    }
  }
  return emailHtmlList
}
