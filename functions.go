package email

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
)

func Verify(ctx context.Context, email string) (valid bool, err error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	select {
	case err = <-signalVerification(ctx, email):
	case <-ctx.Done():
		err = fmt.Errorf("%[1]s: %[2]v", email, ErrTimeout)
	}
	if err != nil {
		return
	}
	valid = true
	return
}

func signalVerification(ctx context.Context, email string) (ch chan error) {
	if ctx == nil {
		ctx = context.TODO()
	}
	ch = make(chan error, 3)
	defer close(ch)
	/*
		Verify if the email address is syntactically valid.
	*/
	e, err := mail.ParseAddress(email)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	i := strings.LastIndex(e.Address, "@")
	if i < 0 || i >= len(email)-1 {
		err = fmt.Errorf("%[1]s: %[2]v", email, ErrInvalidEmail)
		ch <- err
		return
	}
	domain := email[i+1:]
	/*
		Verify if the domain has valid nameserver (NS) records.
	*/
	_, err = net.LookupNS(domain)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	/*
		Verify if the domain is not part of a disposable domain list.
	*/
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, disposableEmailsURL, nil)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			err = fmt.Errorf("%[1]s: %[2]v", email, err)
			ch <- err
			return
		}
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	disposableDomains := []string{}
	err = json.Unmarshal(bytes, &disposableDomains)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	for _, disposableDomain := range disposableDomains {
		if !strings.EqualFold(domain, disposableDomain) {
			continue
		}
		err = fmt.Errorf("%[1]s: %[2]v", email, ErrDisposableEmail)
		break
	}
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	/*
		Verify if the domain has valid mail exchanger (MX) records.
	*/
	records, err := net.LookupMX(domain)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	if len(records) < 1 {
		err = fmt.Errorf("%[1]s: %[2]v", email, ErrNoMXRecords)
		ch <- err
		return
	}
	/*
		Verify if the Mail Transfer Agent (MTA) is reachable.
	*/
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
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	defer func() {
		err = client.Close()
		if err != nil {
			err = fmt.Errorf("%[1]s: %[2]v", email, err)
			ch <- err
			return
		}
	}()
	err = client.Mail(fromEmail)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	err = client.Rcpt(email)
	if err != nil {
		err = fmt.Errorf("%[1]s: %[2]v", email, err)
		ch <- err
		return
	}
	return
}
