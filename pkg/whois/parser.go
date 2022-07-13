package whois

import (
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"bytes"
	"github.com/araddon/dateparse"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func convertTime(str string) (t time.Time) {
	switch {
	case regexp.MustCompile("^\\d{2}-(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)-\\d{4}$").MatchString(str):
		t, _ = time.Parse("02-Jan-2006", str)
	case regexp.MustCompile("^\\d{2}-\\d{2}-\\d{4}$").MatchString(str):
		t, _ = time.Parse("02-01-2006", str)
	case strings.Contains(str, "(UTC+8)"):
		str = strings.ReplaceAll(str, "(UTC+8)", "+08:00")
		t, _ = time.Parse("2006-01-02 15:04:05 Z07:00", str)
	case regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}$").MatchString(str):
		t, _ = time.Parse("2006-01-02 15:04:05", str)
	case regexp.MustCompile("^\\d{2}/\\d{2}/\\d{4} \\d{2}:\\d{2}:\\d{2}$").MatchString(str):
		t, _ = time.Parse("01/02/2006 15:04:05", str)
	default:
		t, _ = dateparse.ParseAny(str)
	}
	return
}

var redactedKeywords = []string{"rdds", "privacy", "please"}

func isRedacted(value string) bool {
	value = strings.ToLower(value)
	for _, keyword := range redactedKeywords {
		if strings.Contains(value, keyword) {
			return true
		}
	}
	return false
}

func Parse(raw []byte) (*models.Whois, error) {
	whois := &models.Whois{
		FetchedAt: time.Now(),
	}
	whois.Raw = string(raw)
	lines := bytes.Split(raw, []byte("\n"))
	for _, line := range lines {
		line = bytes.Trim(line, "\r\n\t ")
		switch {
		case strings.Contains(string(line), "Record expires on "):
			whois.ExpirationDate = convertTime(string(line[len("Record expires on "):]))
			continue
		case strings.Contains(string(line), "Record created on "):
			whois.RegistrationDate = convertTime(string(line[len("Record created on "):]))
			continue
		}
		parts := bytes.Split(line, []byte(":"))
		key := string(bytes.ToLower(bytes.Trim(parts[0], "\t ")))
		value := string(bytes.Trim(bytes.Join(parts[1:], []byte(":")), "\t "))
		switch key {
		case "domain name", "domain":
			whois.DomainName = strings.ToLower(value)
		case "domain status":
			index := strings.IndexRune(value, ' ')
			if index == -1 {
				index = len(value)
			}
			status := strings.ToLower(value[:index])
			b, ok := StatusFromString[status]
			if ok {
				whois.Status |= b
			}
		case "registrant",
			"registrant name",
			"registrant organization",
			"company english name (it should be the same as the registered/corporation" + // for .hk
				" name on your business register certificate or relevant documents)":
			if isRedacted(value) {
				value = "REDACTED"
			}
			whois.Registrant = value
		case "registrant contact email", "registrant email":
			if isRedacted(value) {
				value = "REDACTED"
			}
			whois.RegistrantEmail = value
		case "sponsoring registrar", // for .cn
			"registrar",
			"registrar name": // for .hk
			whois.Registrar = value
		case "registrar iana id":
			id, _ := strconv.ParseUint(value, 10, 64)
			whois.RegistrarIANAID = uint(id)
		case "name server":
			whois.NameServer = append(whois.NameServer, value)
		case "updated date",
			"domain record last updated": // for .edu
			whois.UpdatedDate = convertTime(value)
		case "registration time",
			"creation date",
			"domain record activated",       // for .edu
			"domain name commencement date": // for .hk
			whois.RegistrationDate = convertTime(value)
		case "expiration time",
			"registrar registration expiration date",
			"registry expiry date",
			"domain expires", // for .edu
			"expiry date":    // for .hk and .im
			if whois.ExpirationDate.IsZero() {
				whois.ExpirationDate = convertTime(value)
			}
		case "dnssec":
			whois.DNSSEC = value
		case "dnssec ds data":
			whois.DSData = value
		}
	}
	if whois.DomainName == "" {
		return nil, errors.ErrBadWhoisFormatOrNotRegistered
	}
	return whois, nil
}
