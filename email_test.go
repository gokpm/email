package email

import "testing"

func TestVerify(t *testing.T) {
	for _, data := range testDataVerify {
		valid, err := Verify(data.email)
		if valid == data.valid {
			continue
		}
		t.Fatal(data, valid, err)
	}
}

var testDataVerify = []struct {
	email string
	valid bool
}{
	{
		email: ``,
		valid: false,
	},
	{
		email: `abigail@example.com`,
		valid: false,
	},
	{
		email: `email@example.com`,
		valid: false,
	},
	{
		email: `firstname.lastname@example.com`,
		valid: false,
	},
	{
		email: `email@subdomain.example.com`,
		valid: false,
	},
	{
		email: `firstname+lastname@example.com`,
		valid: false,
	},
	{
		email: `email@123.123.123.123`,
		valid: false,
	},
	{
		email: `email@[123.123.123.123]`,
		valid: false,
	},
	{
		email: `firstname-lastname@example.com`,
		valid: false,
	},
	{
		email: `email@example.co.jp`,
		valid: false,
	},
	{
		email: `email@example.museum`,
		valid: false,
	},
	{
		email: `email@example.name`,
		valid: false,
	},
	{
		email: `_______@example.com`,
		valid: false,
	},
	{
		email: `email@example-one.com`,
		valid: false,
	},
	{
		email: `1234567890@example.com`,
		valid: false,
	},
	{
		email: `plainaddress`,
		valid: false,
	},
	{
		email: `just"not"right@example.com`,
		valid: false,
	},
	{
		email: `this\ is"really"not\allowed@example.com`,
		valid: false,
	},
	{
		email: `“(),:;<>[\]@example.com`,
		valid: false,
	},
	{
		email: `Abc..123@example.com`,
		valid: false,
	},
	{
		email: `email@example..com`,
		valid: false,
	},
	{
		email: `email@111.222.333.44444`,
		valid: false,
	},
	{
		email: `email@example.web`,
		valid: false,
	},
	{
		email: `email@-example.com`,
		valid: false,
	},
	{
		email: `email@example`,
		valid: false,
	},
	{
		email: `email@example.com (Joe Smith)`,
		valid: false,
	},
	{
		email: `あいうえお@example.com`,
		valid: false,
	},
	{
		email: `email..email@example.com`,
		valid: false,
	},
	{
		email: `email.@example.com`,
		valid: false,
	},
	{
		email: `.email@example.com`,
		valid: false,
	},
	{
		email: `email@example@example.com`,
		valid: false,
	},
	{
		email: `email.example.com`,
		valid: false,
	},
	{
		email: `Joe Smith <email@example.com>`,
		valid: false,
	},
	{
		email: `@example.com`,
		valid: false,
	},
	{
		email: `#@%^%#$@#$@#.com`,
		valid: false,
	},
	{
		email: `user@09gmail.com`,
		valid: false,
	},
	{
		email: `user@netplix.site`,
		valid: false,
	},
	{
		email: "mail.gokulpm@gmail.com",
		valid: true,
	},
}
