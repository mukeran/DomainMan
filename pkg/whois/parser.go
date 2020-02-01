package whois

import (
	"DomainMan/models"
	"DomainMan/pkg/errors"
	"bytes"
	"strconv"
	"strings"
	"time"
)

func Parse(raw []byte) (*models.Whois, error) {
	whois := &models.Whois{}
	whois.Raw = string(raw)
	lines := bytes.Split(raw, []byte("\n"))
	for _, line := range lines {
		line = bytes.Trim(line, "\r\n\t ")
		parts := bytes.Split(line, []byte(":"))
		key := string(bytes.ToLower(bytes.Trim(parts[0], "\t ")))
		value := string(bytes.Trim(bytes.Join(parts[1:], []byte(":")), "\t "))
		switch key {
		case "domain name":
			whois.DomainName = strings.ToLower(value)
		case "domain status":
			index := strings.IndexRune(value, ' ')
			status := strings.ToLower(value[:index])
			b, ok := StatusFromString[status]
			if ok {
				whois.Status |= b
			}
		case "registrant", "registrant name":
			whois.Registrant = value
		case "registrant contact email", "registrant Email":
			whois.RegistrantEmail = value
		case "sponsoring registrar", "registrar":
			whois.Registrar = value
		case "registrar iana id":
			id, _ := strconv.ParseUint(value, 10, 64)
			whois.RegistrarIANAID = uint(id)
		case "name server":
			whois.NameServer = append(whois.NameServer, value)
		case "updated date":
			whois.UpdatedDate, _ = time.Parse(time.RFC3339, value)
		case "registration time", "creation date":
			whois.RegistrationDate, _ = time.Parse(time.RFC3339, value)
		case "expiration time", "registrar registration expiration date", "registry expiry date":
			whois.ExpirationDate, _ = time.Parse(time.RFC3339, value)
		case "dnssec":
			whois.DNSSEC = value
		case "dnssec ds data":
			whois.DSData = value
		}
	}
	if whois.DomainName == "" {
		return nil, errors.ErrBadFormat
	}
	return whois, nil
}
