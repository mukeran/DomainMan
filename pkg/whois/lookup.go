package whois

import (
	"DomainMan/models"
	"DomainMan/pkg/database"
	"DomainMan/pkg/errors"
	"fmt"
	"github.com/jinzhu/gorm"
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
	if v := db.Where("domain_name = ?", domain).Order("created_at desc").First(&cachedWhois); gorm.IsRecordNotFoundError(v.Error) {
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
	db := database.DB.Begin()
	defer db.RollbackUnlessCommitted()
	levels := strings.Split(domain, ".")
	var server string
	for i := 0; i < len(levels); i++ {
		name := strings.Join(levels[i:], ".")
		var suffix models.Suffix
		if v := db.Where("name = ?", name).First(&suffix); gorm.IsRecordNotFoundError(v.Error) {
			continue
		} else if v.Error != nil {
			panic(v.Error)
		}
		server = suffix.WhoisServer
	}
	if server == "" {
		return nil, errors.ErrUnsupportedSuffix
	}
	if v := db.Commit(); v.Error != nil {
		panic(v.Error)
	}
	return LookupWithWhoisServer(domain, server)
}

func LookupWithWhoisServer(domain string, server string) (*models.Whois, error) {
	conn, err := net.Dial("tcp", server)
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
