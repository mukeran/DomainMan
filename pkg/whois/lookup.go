package whois

import (
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net"
	"strings"
	"time"
)

const (
	maxBufferSize = 100
	lookupTimeGap = 240 * time.Hour
)

func LookupWithCache(domain string) (*models.Whois, error) {
	db := database.DB
	var cachedWhois models.Whois
	if v := db.Where("domain_name = ?", domain).Order("created_at desc").First(&cachedWhois); errors.Is(v.Error, gorm.ErrRecordNotFound) {
		return Lookup(domain)
	} else if v.Error != nil {
		panic(v.Error)
	}
	if cachedWhois.CreatedAt.Add(lookupTimeGap).Before(time.Now()) {
		return Lookup(domain)
	}
	return &cachedWhois, nil
}

func Lookup(domain string) (*models.Whois, error) {
	var (
		mode   uint
		server string
	)
	db := database.DB
	levels := strings.Split(domain, ".")
	if len(levels) >= 3 {
		for i := 0; i < len(levels); i++ {
			name := strings.Join(levels[i:], ".")
			var suffix models.Suffix
			if v := db.Where("name = ?", name).First(&suffix); errors.Is(v.Error, gorm.ErrRecordNotFound) {
				continue
			} else if v.Error != nil {
				panic(v.Error)
			}
			mode = suffix.Mode
			server = suffix.WhoisServer
			break
		}
	} else {
		var suffix models.Suffix
		if v := db.Where("name = ?", levels[len(levels)-1]).First(&suffix); errors.Is(v.Error, gorm.ErrRecordNotFound) {
		} else if v.Error != nil {
			panic(v.Error)
		} else {
			mode = suffix.Mode
			server = suffix.WhoisServer
		}
	}
	if server == "" {
		return nil, errors.ErrUnsupportedSuffix
	}
	switch mode {
	case models.ModeWhois:
		return LookupWithWhoisServer(domain, server)
	case models.ModeWeb:
		fallthrough
	default:
		return nil, errors.ErrUnsupportedWhoisMode
	}
}

func LookupWithWhoisServer(domain string, server string) (*models.Whois, error) {
	conn, err := net.DialTimeout("tcp", server, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(domain + "\r\n"))
	if err != nil {
		return nil, err
	}
	var raw []byte
	buf := make([]byte, maxBufferSize)
	size := 0
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		raw = append(raw, buf[:n]...)
		size += n
	}
	fmt.Printf("Received %d bytes from server\n", size)
	return Parse(raw)
}
