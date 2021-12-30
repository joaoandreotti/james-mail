package database

import (
  "golang.org/x/crypto/bcrypt"
  "encoding/base64"
  "log"
)

const hashCost = 16

func GenerateHashedPassword(password string) string {
  bytePassword := []byte(password)
  hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, hashCost)
  if err != nil {
    log.Fatal(err)
  }
  return base64.URLEncoding.EncodeToString(hashedPassword)
}

func ValidateHashedPassword(inputedPassword, storedPassword string) bool {
  hashedPassword, err := base64.URLEncoding.DecodeString(storedPassword)
  if err != nil {
    return false
  }
  err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(inputedPassword))
  return err == nil
}
