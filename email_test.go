package email

import (
	"testing"
)

var testDataVerify = []struct {
	email string
	valid bool
}{}

func TestVerify(t *testing.T) {
	for _, data := range testDataVerify {
		valid, err := Verify(data.email)
		if valid == data.valid {
			continue
		}
		t.Fatal(data, valid, err)
	}
}
