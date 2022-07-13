package containers

import (
	"DomainMan/pkg/gui/instance"
	"DomainMan/pkg/models"
	"DomainMan/pkg/whois"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/idna"
	"strings"
)

func whoisDetailParsed(inst *instance.Instance, w *models.Whois) fyne.CanvasObject {
	c := container.NewGridWithColumns(4)
	c.Add(canvas.NewText("Domain Name:", theme.ForegroundColor()))
	if strings.Contains(w.DomainName, "xn--") {
		parsed, err := idna.ToUnicode(w.DomainName)
		if err != nil {
			parsed = "Error parsing domain name"
		}
		c.Add(canvas.NewText(fmt.Sprintf("%s (%s)", w.DomainName, parsed), theme.ForegroundColor()))
	} else {
		c.Add(canvas.NewText(w.DomainName, theme.ForegroundColor()))
	}
	c.Add(canvas.NewText("Fetched At:", theme.ForegroundColor()))
	c.Add(canvas.NewText(w.FetchedAt.Format("2006-01-02 15:04:05"), theme.ForegroundColor()))
	c.Add(canvas.NewText("Expiration Date:", theme.ForegroundColor()))
	if w.ExpirationDate.IsZero() {
		c.Add(canvas.NewText("-", theme.ForegroundColor()))
	} else {
		c.Add(canvas.NewText(w.ExpirationDate.Format("2006-01-02 15:04:05"), theme.ForegroundColor()))
	}
	c.Add(canvas.NewText("Registration Date:", theme.ForegroundColor()))
	if w.RegistrationDate.IsZero() {
		c.Add(canvas.NewText("-", theme.ForegroundColor()))
	} else {
		c.Add(canvas.NewText(w.RegistrationDate.Format("2006-01-02 15:04:05"), theme.ForegroundColor()))
	}
	if w.Registrar != "" {
		c.Add(canvas.NewText("Registrar:", theme.ForegroundColor()))
		c.Add(canvas.NewText(w.Registrar, theme.ForegroundColor()))
		if w.RegistrarIANAID != 0 {
			c.Add(canvas.NewText("Registrar IANA ID:", theme.ForegroundColor()))
			c.Add(canvas.NewText(fmt.Sprintf("%d", w.RegistrarIANAID), theme.ForegroundColor()))
		} else {
			c.Add(canvas.NewText("", theme.ForegroundColor()))
			c.Add(canvas.NewText("", theme.ForegroundColor()))
		}
	}
	if w.Registrant != "" {
		c.Add(canvas.NewText("Registrant:", theme.ForegroundColor()))
		c.Add(canvas.NewText(w.Registrant, theme.ForegroundColor()))
		c.Add(canvas.NewText("Registrant Email: ", theme.ForegroundColor()))
		c.Add(canvas.NewText(w.RegistrantEmail, theme.ForegroundColor()))
	}
	if len(w.NameServer) != 0 || w.Status != 0 {
		var statusStrings []string
		for _, status := range whois.StatusList {
			if w.Status&status == status {
				statusStrings = append(statusStrings, whois.StatusToString[status])
			}
		}
		larger := len(w.NameServer)
		if len(statusStrings) > larger {
			larger = len(statusStrings)
		}
		for index := 0; index < larger; index++ {
			if index < len(w.NameServer) {
				c.Add(canvas.NewText(fmt.Sprintf("Name Server %d:", index+1), theme.ForegroundColor()))
				c.Add(canvas.NewText(strings.ToLower(w.NameServer[index]), theme.ForegroundColor()))
			} else {
				c.Add(canvas.NewText("", theme.ForegroundColor()))
				c.Add(canvas.NewText("", theme.ForegroundColor()))
			}
			if index < len(statusStrings) {
				c.Add(canvas.NewText(fmt.Sprintf("Domain Status %d:", index+1), theme.ForegroundColor()))
				c.Add(canvas.NewText(statusStrings[index], theme.ForegroundColor()))
			} else {
				c.Add(canvas.NewText("", theme.ForegroundColor()))
				c.Add(canvas.NewText("", theme.ForegroundColor()))
			}
		}
	}
	c.Add(canvas.NewText("Raw Whois Data:", theme.ForegroundColor()))
	return container.NewPadded(c)
}

func whoisDetailRaw(inst *instance.Instance, whois *models.Whois) fyne.CanvasObject {
	text := widget.NewMultiLineEntry()
	text.SetText(whois.Raw)
	text.Disable()
	return text
}

func WhoisDetail(inst *instance.Instance, whois *models.Whois) fyne.CanvasObject {
	return container.NewBorder(whoisDetailParsed(inst, whois), nil, nil, nil, whoisDetailRaw(inst, whois))
}
