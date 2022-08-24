# email

## Summary

Verify if an email address is syntactically valid and reachable.

## Description

1. Verify if the email address is syntactically valid.
2. Verify if the domain has valid nameserver (NS) records.
3. Verify if the domain is not part of a disposable domain list.
4. Verify if the domain has valid mail exchanger (MX) records.
5. Verify if the Mail Transfer Agent (MTA) is reachable.

## Installation

```
go get -u github.com/gokpm/email
```

## Example

```
valid, err = email.Verify(`user@example.com`)
```
