package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	flagEnableSES = flag.Bool("enable-ses", false, "enable ses")
	emailBcc      = flag.String("email-bcc", "", "Bcc: email address")
	emailFrom     = flag.String("email-from", "", "From: email address")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", indexHtmlHandler)
	http.Handle("/add/", newSecretsManagerHandler())
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
