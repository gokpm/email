package email

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
)

func Verify(email string) (valid bool, err error) {
	_, err = mail.ParseAddress(email)
	if err != nil {
		return
	}
	i := strings.LastIndex(email, "@")
	if i < 0 || i >= len(email)-1 {
		err = ErrInvalidEmail
		return
	}
	domain := email[i+1:]
	_, err = net.LookupNS(domain)
	if err != nil {
		return
	}
	resp, err := http.Get(disposableEmailsURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	disposableDomains := []string{}
	err = json.Unmarshal(bytes, &disposableDomains)
	if err != nil {
		return
	}
	for _, disposableDomain := range disposableDomains {
		if !strings.EqualFold(domain, disposableDomain) {
			continue
		}
		err = ErrDisposableEmail
		break
	}
	if err != nil {
		return
	}
	records, err := net.LookupMX(domain)
	if err != nil {
		return
	}
	if len(records) < 1 {
		err = ErrNoMXRecords
		return
	}
	host := records[0].Host
	pref := records[0].Pref
	for _, record := range records {
		if record.Pref >= pref {
			continue
		}
		pref = record.Pref
		host = record.Host
	}
	addr := fmt.Sprintf("%[1]s:%[2]d", host, smtpPort)
	client, err := smtp.Dial(addr)
	if err != nil {
		return
	}
	defer client.Close()
	err = client.Mail(fromEmail)
	if err != nil {
		return
	}
	err = client.Rcpt(email)
	if err != nil {
		return
	}
	valid = true
	return
}