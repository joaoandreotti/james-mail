package pgp

import (
  "github.com/ProtonMail/gopenpgp/v2/crypto"
  "github.com/ProtonMail/gopenpgp/v2/helper"
)

const openPgpKeyType = "rsa"
const openPgpRsaSize = 4096

func GeneratePgpKeyPair(username, email, password string) (string, string) {
  privateKey, err := helper.GenerateKey(username, email, []byte(password),
    openPgpKeyType, openPgpRsaSize)
  if err != nil {
    return "", ""
  }
  key, _ := crypto.NewKeyFromArmored(privateKey)
  publicKey, _ := key.GetArmoredPublicKey()
  return publicKey, privateKey
}

func EncryptMessage(publicKey, message string) string {
  armor, err := helper.EncryptMessageArmored(publicKey, message)
  if err != nil {
    return ""
  }
  return armor
}

func DecryptMessage(privateKey, message, password string) string {
  plain, err := helper.DecryptMessageArmored(privateKey, []byte(password), message)
  if err != nil {
    return ""
  }
  return plain
}
