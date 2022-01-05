package smtp

import (
  "james-mail/database"
  "james-mail/pgp"
)

func GetReceivedEmailsList(userEmail, password string) (emailList []string) {
  _, privateKey := database.GetKeyPairFromEmail(userEmail)
  encryptedEmailList := database.GetReceivedEmailsFromEmail(userEmail)
  for _, encryptedEmail := range encryptedEmailList {
    decryptedEmail := pgp.DecryptMessage(privateKey, encryptedEmail.EmailBody, password)
    emailAddress := encryptedEmail.EmailAddress + ";"
    if len(decryptedEmail) > 0 {
      emailList = append(emailList, emailAddress + decryptedEmail)
    } else {
      emailList = append(emailList, emailAddress + encryptedEmail.EmailBody)
    }
  }
  return emailList
}

func GetSentEmailsList(userEmail, password string) (emailList []string) {
  _, privateKey := database.GetKeyPairFromEmail(userEmail)
  encryptedEmailList := database.GetSentEmailsFromEmail(userEmail)
  for _, encryptedEmail := range encryptedEmailList {
    decryptedEmail := pgp.DecryptMessage(privateKey, encryptedEmail.EmailBody, password)
    emailAddress := encryptedEmail.EmailAddress + ";"
    if len(decryptedEmail) > 0 {
      emailList = append(emailList, emailAddress + decryptedEmail)
    } else {
      emailList = append(emailList, emailAddress + encryptedEmail.EmailBody)
    }
  }
  return emailList
}
