package main

import (
	"encoding/base64"
	"encoding/gob"
	"log"
	"math/rand"
	"os"
	"strings"
)

const (
	pwdChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length   = 16
)

func genPassword() string {
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = pwdChars[rand.Intn(len(pwdChars))]
	}
	return string(buf)
}

func writeWrapper(n int, err error) {
	if err != nil {
		log.Println(n, err)
	}
}

func eapKey(email string) string {
	prefix := "eap-"
	suffix := strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(email)), "=")
	return prefix + suffix

}

func SaveEAPtoFile(secretsMap EAPSecretsMap) {
	filePath := "swanctl.eap.secrets"
	err := os.WriteFile(filePath, []byte(secretsMap.String()), 0644)
	if err != nil {
		log.Println("Couldn't write to  file " + filePath + ". err: " + err.Error())
	}

}

func DumpEAPtoFile(secretsMap EAPSecretsMap) {
	filePath := "swanctl.eap.dump"
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("Couldn't open file " + filePath + ". err: " + err.Error())
	}
	enc := gob.NewEncoder(f)
	if err := enc.Encode(secretsMap); err != nil {
		log.Println("Couldn't write to  file " + filePath + ". err: " + err.Error())
	}

}

func RestoreDumpFromFile() EAPSecretsMap {
	filePath := "swanctl.eap.dump"
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		log.Println("Couldn't open file " + filePath + ". err: " + err.Error())
		return make(EAPSecretsMap)
	}
	dec := gob.NewDecoder(f)
	var eap EAPSecretsMap
	if err := dec.Decode(&eap); err != nil {
		log.Println("Couldn't decode file " + filePath + ". err: " + err.Error())
		return make(EAPSecretsMap)
	}

	return eap
}
