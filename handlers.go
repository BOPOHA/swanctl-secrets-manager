package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"
)

type EAPSecret struct {
	Secret, Id string
	Expdate    int64
}
type EAPSecretsMap map[string]EAPSecret

func (m EAPSecretsMap) String() string {
	out := ""
	for x, y := range m {
		out += x + " {\n    " +
			"id = " + y.Id + "\n    " +
			"secret = " + y.Secret + "\n" +
			"}\n"
	}
	return out
}

type smConfig struct {
	emDomains []string
}

type secretsManager struct {
	mu     sync.Mutex
	client *sesv2.Client
	conf   smConfig
	eap    EAPSecretsMap
}

func (m *secretsManager) setEAPSecret(email, password string) {
	m.mu.Lock()
	m.eap[eapKey(email)] = EAPSecret{
		Secret: password,
		Id:     email,
		//Expdate: time.Now().AddDate(0, 0, 7).UTC().Unix(),
		Expdate: time.Now().Add(time.Minute * 60).UTC().Unix(),
	}
	log.Println("added user " + email + "/" + password)
	SaveEAPtoFile(m.eap)
	DumpEAPtoFile(m.eap)
	m.mu.Unlock()
}
func (m *secretsManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reqEmail := r.FormValue("email")
		ok, email := m.emailValidator(reqEmail)
		if ok != true {
			w.WriteHeader(http.StatusBadRequest)
			writeWrapper(w.Write([]byte("bad email")))
			return
		}
		password := genPassword()
		m.setEAPSecret(email, password)
		if m.client == nil { // for test and debug
			log.Println("mail client is nil")
			return
		}
		// send email
		emailConf := getDefaultEmailConfig()
		emailConf.setEmailTo(email)
		emailConf.body = "your temporary password is: \n" + password + "\n"
		input := getSESEmailInput(emailConf)
		_, err := m.client.SendEmail(context.TODO(), input)
		if err != nil {
			log.Printf("Unable to send email %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		writeWrapper(w.Write([]byte("please check your email: " + email)))
		return

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (m *secretsManager) CleanExpiredEAP() {
	m.mu.Lock()
	for p, s := range m.eap {
		anyBodyDeleted := false
		if s.Expdate < time.Now().Unix() {
			delete(m.eap, p)
			log.Println("deleted user " + s.Id + "/" + s.Secret)
			anyBodyDeleted = true
		}
		if anyBodyDeleted {
			SaveEAPtoFile(m.eap)
			DumpEAPtoFile(m.eap)
		}
	}
	//log.Printf("secrets:\n%s\n\n", m.eap)
	m.mu.Unlock()
}

func (m *secretsManager) emailValidator(email string) (ok bool, addres string) {
	adr, err := mail.ParseAddress(email)
	if err != nil {
		// return false
		return
	}
	addres = adr.Address

	// if no trusted domains - pass all valid emails
	if len(m.conf.emDomains) == 0 {
		return true, addres
	}
	// checking that email domain in trusted list
	for _, d := range m.conf.emDomains {
		if strings.HasSuffix(adr.Address, "@"+d) {
			return true, addres
		}
	}
	// return false
	return
}

func newSecretsManagerHandler() http.Handler {
	var sesclient *sesv2.Client
	var emDomains []string

	eap := RestoreDumpFromFile()
	if *flagEnableSES {
		sesclient = getSesClient()
	}
	if *trustedDomains != "" {
		emDomains = strings.Split(*trustedDomains, ",")
	}
	config := smConfig{
		emDomains: emDomains,
	}

	log.Printf("restored from dump:\n%s\n\n", eap)
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	sm := &secretsManager{
		client: sesclient,
		conf:   config,
		eap:    eap,
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				sm.CleanExpiredEAP()
			}
		}
	}()
	return sm
}

func indexHtmlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		indexPage := `<!DOCTYPE html>
<html>
<body>

<h2>Add/update user</h2>

<form action="/add/" method="post">
  <label for="fname">Email:</label><br>
  <input type="email" id="email" name="email">
  <input type="submit" value="Submit">
</form>

<p>Note ...</p>

</body>
</html>
`
		writeWrapper(w.Write([]byte(indexPage)))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
